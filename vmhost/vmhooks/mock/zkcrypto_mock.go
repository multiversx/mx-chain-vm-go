package mock

import "github.com/stretchr/testify/mock"

// EllipticCurveMock is a mock for the lowLevelFeatures.EllipticCurve interface
type EllipticCurveMock struct {
	mock.Mock
}

func (ec *EllipticCurveMock) Add(point1Bytes, point2Bytes []byte) ([]byte, error) {
	args := ec.Called(point1Bytes, point2Bytes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (ec *EllipticCurveMock) Mul(pointBytes, scalarBytes []byte) ([]byte, error) {
	args := ec.Called(pointBytes, scalarBytes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (ec *EllipticCurveMock) MultiExp(pointsVec, scalarsVec [][]byte) ([]byte, error) {
	args := ec.Called(pointsVec, scalarsVec)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (ec *EllipticCurveMock) MapToCurve(element []byte) ([]byte, error) {
	args := ec.Called(element)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// PairingMock is a mock for the lowLevelFeatures.Pairing interface
type PairingMock struct {
	mock.Mock
}

func (p *PairingMock) PairingCheck(pointsG1, pointsG2 [][]byte) (bool, error) {
	args := p.Called(pointsG1, pointsG2)
	return args.Bool(0), args.Error(1)
}
