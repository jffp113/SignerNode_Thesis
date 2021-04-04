package api

import "SignerNode/signermanager"

type GenericRespChan func() <-chan signermanager.ManagerResponse
type SignFunc func(data []byte) <-chan signermanager.ManagerResponse
type VerifyFunc func(data []byte) <-chan signermanager.ManagerResponse
type MembershipFunc func(data []byte) <-chan signermanager.ManagerResponse
type InstallShareFunc func(data []byte) <-chan signermanager.ManagerResponse

func ConvertToGeneric(f GenericRespChan) func(data []byte) <-chan signermanager.ManagerResponse {
	return func([]byte) <-chan signermanager.ManagerResponse {
		return f()
	}
}
