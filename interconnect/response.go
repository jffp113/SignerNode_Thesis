package interconnect

type Status int

//Protocol

const (
	Ok Status = iota
	Error
	InvalidTransaction
)

type HandlerResponse struct {
	ResponseStatus Status
	ResponseData   []byte
	Err            error
}

func CreateErrorMessage(err error) HandlerResponse{
	logger.Error(err)
	return HandlerResponse{
		ResponseStatus: Error,
		Err:            err,
	}
}

func CreateInvalidTransactionMessage() HandlerResponse{
	return HandlerResponse{
		ResponseStatus: InvalidTransaction,
		Err:            nil,
	}
}

func CreateOkMessage(data []byte) HandlerResponse{
	return HandlerResponse{
		ResponseStatus: Ok,
		ResponseData:   data,
		Err:            nil,
	}
}

func SendErrorMessage(c chan<- HandlerResponse, err error) {
	logger.Error(err)
	c <- CreateErrorMessage(err)
}

func SendInvalidTransactionMessage(c chan<- HandlerResponse) {
	c <- CreateInvalidTransactionMessage()
}

func SendOkMessage(c chan<- HandlerResponse, data []byte) {
	c <- CreateOkMessage(data)
}