package elliptic_curves

import (
	"crypto/elliptic"
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
