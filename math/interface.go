package math

type RandomnessGenerator interface {
	Read(p []byte) (n int, err error)
	IsInterfaceNil() bool
}
