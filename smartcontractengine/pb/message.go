package pb

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// Generate a new UUID
func GenerateId() string {
	return fmt.Sprint(uuid.NewV4())
}

func CreateMessage(msgType Message_MessageType, data []byte) ([]byte, string, error) {
	corrId := GenerateId()
	return CreateSignMessageWithCorrelationId(msgType, data, corrId)
}

func CreateSignMessageWithCorrelationId(msgType Message_MessageType, data []byte, corrId string) ([]byte, string, error) {
	b, err := proto.Marshal(&Message{
		MessageType:          msgType,
		CorrelationId: corrId,
		Content:       data,
	})
	return b, corrId, err
}

func CreateMessageWithCorrelationId(msgType Message_MessageType, data proto.Message, corrId string, handlerId string) (*Message, string, error) {
	bytes, err := proto.Marshal(data)

	if err != nil {
		return nil, "", err
	}

	hd := Message{
		MessageType:   msgType,
		CorrelationId: corrId,
		Content:       bytes,
		HandlerId:     handlerId,
	}
	return &hd, corrId, nil
}

func CreateHandlerMessage(msgType Message_MessageType, data proto.Message, handlerId string) (*Message, string, error) {
	return CreateMessageWithCorrelationId(msgType, data, GenerateId(), handlerId)
}

func UnmarshallSignMessage(data []byte) (*Message, error) {
	msg := Message{}

	err := proto.Unmarshal(data, &msg)

	return &msg, err
}
