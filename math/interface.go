package math

// RandomnessGenerator will provide the interface to the main functionalities of the VM where randomness is needed
type RandomnessGenerator interface {
	Read(p []byte) (n int, err error)
	IsInterfaceNil() bool
}
