package interconnect

import (
	"context"
	"errors"
	"github.com/ipfs/go-log/v2"
)

var logger = log.Logger("interconnect")

const (
	//Size of waiting events
	EventPoolSize = 100

	//Default Number of Workers
	WORKERS = 20
)

//P2pContext defines a P2P network P2pContext used by different
//defined protocols.
type P2pContext struct {
	//Broadcast broadcasts a message to all P2P network members.
	Broadcast func(msg []byte) error

	//BroadcastToGroup only broadcast to a specific group.
	BroadcastToGroup func(groupId string ,msg []byte) error

	//JoinGroup joins group can create a group over a P2P network.
	JoinGroup func(groupId string) error

	//LeaveGroup leave a previously joined group.
	LeaveGroup func(groupId string) error
}

//Handler defines the type of a interconnect event
type Handler func(content []byte, ctx P2pContext) HandlerResponse


type Interconnect interface {
	RegisterHandler(t HandlerType, handler Handler)
	EmitEvent(t HandlerType, content []byte) HandlerResponse
}

//interconnect defines a event base system.
//interconnect supports Handlers to be register
//to react to events
type interconnect struct {

	//handlers used to react to certain event types (HandlerType)
	//used by handler workers
	handlers map[HandlerType][]Handler

	//eventChan is a chan used to send events to handler
	//workers
	eventChan chan event

	//nWorkers define the number of handler workers
	nWorkers int

	//P2P network P2pContext. User by the handlers to communicate in a P2P network
	p2pCtx P2pContext

	cancel context.CancelFunc
}

//event defines a emitted event by EmitEvent
type event struct {
	//handlerType is used to choose the process handlers
	handlerType HandlerType

	//content to be used by a handler
	content []byte

	//respChan is used by a handler worker to send
	//the final response produced by the handlers
	respChan chan<- HandlerResponse
}

//NewInterconnect creates a new interconnect using configurations.
func NewInterconnect(configs ...Config) (*interconnect,error){
	ic := interconnect{
		handlers:  make(map[HandlerType][]Handler),
		eventChan: make(chan event,EventPoolSize),
		nWorkers:  WORKERS,
		cancel:    func() {},
	}

	//process every given configuration
	for _,conf := range configs {
		err := conf(&ic)
		if err != nil {
			return &ic,err
		}
	}


	//create a context to terminate workers
	ctx,cancel := context.WithCancel(context.Background())
	ic.cancel = cancel

	//initialize different handler workers
	for i := 0 ; i < ic.nWorkers; i++ {
		go handlerWorker(&ic,ctx)
	}


	return &ic,nil
}

func (ic *interconnect) Done(){
	ic.cancel()
}

//handlerWorker represents a handler which waits for new events and process them
func handlerWorker(ic *interconnect, ctx context.Context) {
	for {
		select {
			case <-ctx.Done():
				logger.Info("Terminating handler worker")
				return
			case event := <- ic.eventChan:
				handlers := ic.handlers[event.handlerType]
				event.respChan <- processEvent(handlers,event.content,ic.p2pCtx)
		}
	}
}

//processEvent processes a event using all the available handlers to a HandlerType
func processEvent(handlers []Handler, content []byte, ctx P2pContext) HandlerResponse{
	var resp HandlerResponse

	if len(handlers) == 0 {
		return CreateErrorMessage(errors.New("no handlers available to process request"))
	}

	for _,handler := range handlers{
		resp = handler(content,ctx)
		content = resp.ResponseData
	}
	return resp
}

//RegisterHandler registers a handler to react to certain
//events. If more than one handler is register for the
//same HandlerType the interconnect will process by the register order
//returning the last value to the event publisher.
//Response from a HandlerResponse will be feed to the next handler.
//RegisterHandler is not reentrant. In presence of concurrency it is
//necessary to synchronize.
func (i *interconnect) RegisterHandler(t HandlerType, handler Handler) {
	i.handlers[t] = append(i.handlers[t],handler)
}


//EmitEvent emits events to be processed by a Handler
func (i *interconnect) EmitEvent(t HandlerType, content []byte) HandlerResponse{
	respChan := make(chan HandlerResponse)
	i.eventChan<-event{t,content,respChan}
	return <-respChan
}




