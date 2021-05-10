package elliptic_curves

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

type ellipticCurve struct {
	curve elliptic.CurveParams
}

func NewEllipticCurve(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int) *ellipticCurve {
	myEc := &ellipticCurve{}
	myEc.curve = elliptic.CurveParams{P: P, N: N, B: B, Gx: Gx, Gy: Gy, BitSize: BitSize, Name: "EC"}
	return myEc
}

func Add(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x1 *big.Int, y1 *big.Int, x2 *big.Int, y2 *big.Int) (x *big.Int, y *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.Add(x1, y1, x2, y2)
}

func Double(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x *big.Int, y *big.Int) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.Double(x, y)
}

func IsOnCurve(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x *big.Int, y *big.Int) bool {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.IsOnCurve(x, y)
}

func ScalarBaseMult(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, k []byte) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.ScalarBaseMult(k)
}

func ScalarMult(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, Bx *big.Int, By *big.Int, k []byte) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return myEc.curve.ScalarMult(Bx, By, k)
}

func Marshal(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x *big.Int, y *big.Int) []byte {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.Marshal(&myEc.curve, x, y)
}

func Unmarshal(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, data []byte) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.Unmarshal(&myEc.curve, data)
}

func GenerateKey(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int) (priv []byte, x, y *big.Int, err error) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.GenerateKey(&myEc.curve, rand.Reader)
}

func MarshalCompressed(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, x *big.Int, y *big.Int) []byte {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.MarshalCompressed(&myEc.curve, x, y)
}

func UnmarshalCompressed(P *big.Int, N *big.Int, B *big.Int, Gx *big.Int, Gy *big.Int, BitSize int, data []byte) (*big.Int, *big.Int) {
	myEc := NewEllipticCurve(P, N, B, Gx, Gy, BitSize)
	return elliptic.UnmarshalCompressed(&myEc.curve, data)
}
