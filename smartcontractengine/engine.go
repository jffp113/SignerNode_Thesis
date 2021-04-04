package smartcontractengine

import (
	"github.com/jffp113/SignerNode_Thesis/messaging"
	"github.com/jffp113/SignerNode_Thesis/smartcontractengine/pb"
	"github.com/golang/protobuf/proto"
	"github.com/ipfs/go-log"
	zmq "github.com/pebbe/zmq4"
	"io"
)

var logger = log.Logger("smartcontract_engine")

const RegisterChanSize = 5

type scClient struct {
	requests            map[string]string
	handlers            map[string]string
	registerHandlerChan chan *pb.Message
	conn                *messaging.ZmqConnection
	clients             *messaging.ZmqConnection
	context             *zmq.Context
}

type context struct {
	scAddress string
	client    *scClient
	worker    *messaging.ZmqConnection
}

func (c *context) Close() error {
	c.worker.Close()
	return nil
}

func (c *context) InvokeSmartContract(payload []byte) ScResponse {

	req := pb.SmartContractValidationRequest{
		Payload:              payload,
		SmartContractAddress: c.scAddress,
	}

	handlerId := c.client.handlers[c.scAddress]

	msg, _, _ := pb.CreateHandlerMessage(pb.Message_SMART_CONTRACT_VALIDATE_REQUEST, &req, handlerId)

	reply, err := c.sendHandlerMessageAndReceiveResponse(msg)

	if err != nil {
		panic("error requesting smartcontract validation")
	}

	if reply.MessageType != pb.Message_SMART_CONTRACT_VALIDATE_RESPONSE {
		panic("Wrong message received")
	}

	replyTHS := pb.SmartContractValidationResponse{}

	err = proto.Unmarshal(reply.Content, &replyTHS)

	if err != nil {
		panic(err)
	}

	return ScResponse{
		T:      int(replyTHS.T),
		N:      int(replyTHS.N),
		Scheme: replyTHS.SignatureScheme,
		Valid:  replyTHS.Status == pb.SmartContractValidationResponse_OK,
		Error:  replyTHS.Status == pb.SmartContractValidationResponse_INTERNAL_ERROR,
	}
}

func (c *scClient) GetContext(scAddress string) (SCContext, io.Closer) {
	worker, err := messaging.NewConnection(c.context, zmq.DEALER, "inproc://workers", false)

	if err != nil {
		panic(err)
	}

	r := context{scAddress, c, worker}
	return &r, &r
}

func NewSmartContractClientFactory(uri string) (SCContextFactory, error) {
	context, _ := zmq.NewContext()

	conn, err := messaging.NewConnection(context, zmq.ROUTER, uri, true)

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	clients, err := messaging.NewConnection(context, zmq.ROUTER, "inproc://workers", true)
	if err != nil {
		return nil, err
	}

	c := scClient{
		requests:            make(map[string]string),
		handlers:            make(map[string]string),
		registerHandlerChan: make(chan *pb.Message, RegisterChanSize),
		conn:                conn,
		clients:             clients,
		context:             context,
	}

	go c.receive()
	go c.processNewHandlers()

	return &c, nil
}

func (c *scClient) Close() error {
	c.conn.Close()
	return nil
}

func (c *scClient) receive() {
	poller := zmq.NewPoller()

	poller.Add(c.conn.Socket(), zmq.POLLIN)
	poller.Add(c.clients.Socket(), zmq.POLLIN)

	for {
		polled, err := poller.Poll(-1)
		logger.Debug("Received messages")
		if err != nil {
			logger.Error("Error Polling messages from socket")
			return
		}
		for _, ready := range polled {
			switch socket := ready.Socket; socket {
			case c.conn.Socket():
				c.handleConnSocket()
			case c.clients.Socket():
				c.handleClientSocket()
			}
		}
	}

}

func (c *scClient) processNewHandlers() {
	worker, err := messaging.NewConnection(c.context, zmq.DEALER, "inproc://workers", false)

	if err != nil {
		panic(err)
	}

	for newMsg := range c.registerHandlerChan {
		req := pb.SmartContractRegisterRequest{}

		err := proto.Unmarshal(newMsg.Content, &req)

		if err != nil {
			logger.Warnf("Error Ignoring register handler MSG: %v", err)
			continue
		}
		logger.Debugf("Registering %v from %v", req.SmartContractAddress, newMsg.HandlerId)

		c.handlers[req.SmartContractAddress] = newMsg.HandlerId

		rep := pb.SmartContractRegisterResponse{Status: pb.SmartContractRegisterResponse_OK}

		handlerMsg, _, err := pb.CreateMessageWithCorrelationId(pb.Message_SMART_CONTRACT_REGISTER_RESPONSE,
			&rep, newMsg.CorrelationId, newMsg.HandlerId)

		if err != nil {
			logger.Warnf("Error Ignoring register handler MSG: %v", err)
			continue
		}

		data, err := proto.Marshal(handlerMsg)

		worker.SendData("", data)

	}
}

func (c *scClient) handleConnSocket() {
	handlerId, data, err := c.conn.RecvData()

	if err != nil {
		logger.Warnf("Error Ignoring MSG: %v", err)
		return
	}

	msg := pb.Message{}

	err = proto.Unmarshal(data, &msg)

	if err != nil {
		logger.Warnf("Error Ignoring MSG: %v", err)
		return
	}

	if msg.MessageType == pb.Message_SMART_CONTRACT_REGISTER_REQUEST {
		logger.Debug("Register Handler MSG received")
		msg.HandlerId = handlerId
		c.registerHandlerChan <- &msg
	} else {
		logger.Debug("Received response from a handler")
		v, present := c.requests[msg.CorrelationId]
		if !present {
			logger.Warnf("MSG not expected, ignoring MSG")
			return
		}

		err := c.clients.SendData(v, data)
		if err != nil {
			logger.Error("Error Sending the message")
			return
		}

		delete(c.requests, msg.CorrelationId)
	}
}

func (c *scClient) handleClientSocket() {
	logger.Debug("Received data")

	clientId, data, err := c.clients.RecvData()

	handlerMsg := pb.Message{}
	err = proto.Unmarshal(data, &handlerMsg)

	if err != nil {
		logger.Error("Error sending out: %v", err)
		return
	}

	c.requests[handlerMsg.CorrelationId] = clientId

	err = c.conn.SendData(handlerMsg.HandlerId, data)
	logger.Debug("Sent msg out to signer")

	if err != nil {
		logger.Warnf("Error retransmitting message", err)
	}

}

func (c *context) sendHandlerMessageAndReceiveResponse(msg *pb.Message) (*pb.Message, error) {
	data, err := proto.Marshal(msg)

	if err != nil {
		return nil, err
	}

	err = c.worker.SendData("", data)

	if err != nil {
		return nil, err
	}

	_, recvData, err := c.worker.RecvData()

	if err != nil {
		return nil, err
	}

	reply, err := pb.UnmarshallSignMessage(recvData)

	if err != nil {
		return nil, err
	}

	return reply, nil
}
