package signermanager

type Status int

const (
	Ok Status = iota
	Error
	InvalidTransaction
)

type ManagerResponse struct {
	ResponseStatus Status
	err error
}


func sendErrorMessage(c chan<- ManagerResponse, err error){
	c<- ManagerResponse{
		ResponseStatus: Error,
		err:            err,
	}
}

func sendInvalidTransactionMessage(c chan<- ManagerResponse){
	c<- ManagerResponse{
		ResponseStatus: InvalidTransaction,
		err:            nil,
	}
}

func sendOkMessage(c chan<- ManagerResponse){
	c<- ManagerResponse{
		ResponseStatus: Ok,
		err:            nil,
	}
}
