package interconnect

type HandlerType int

const (
	SignClientRequest HandlerType = iota
	VerifyClientRequest
	InstallClientRequest
	MembershipClientRequest

	NetworkMessage
)
