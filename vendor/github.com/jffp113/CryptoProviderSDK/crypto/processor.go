package crypto

import (
	"github.com/ipfs/go-log"
	"github.com/jffp113/go-util/messaging/routerdealerhandlers/processor"
)

var logger = log.Logger("signer_processor")

type SignerProcessor struct {
	proc *processor.HandlerProcessor
}

func NewSignerProcessor(uri string) *SignerProcessor {
	return &SignerProcessor{processor.NewHandlerProcessor(uri)}
}

func (self *SignerProcessor) AddHandler(handler THSignerHandler) {
	self.proc.AddHandler(&handlerDecorator{handler})
}

func (self *SignerProcessor) Start() error {
	return self.proc.Start()
}

