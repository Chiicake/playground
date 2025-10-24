package elliptic

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// 常量定义
var (
	COEFA         = big.NewInt(4)
	COEFB         = big.NewInt(27)
	PRIMESECURITY = 500
	TWO           = big.NewInt(2)
	THREE         = big.NewInt(3)
)

// EllipticCurve 表示椭圆曲线 y^2 = x^3 + ax + b (mod p)
type EllipticCurve struct {
	A, B, P   *big.Int
	Generator *ECPoint
	N         *big.Int // 曲线阶数
}

// NewEllipticCurve 创建一个新的椭圆曲线实例
func NewEllipticCurve(a, b, p, n *big.Int) (*EllipticCurve, error) {
	// 检查p是否为质数
	if !isProbablePrime(p, PRIMESECURITY) {
		return nil, errors.New("不安全参数: p不是质数")
	}

	curve := &EllipticCurve{A: a, B: b, P: p, N: n}

	// 检查曲线是否为奇异曲线
	if curve.IsSingular() {
		return nil, errors.New("曲线是奇异的")
	}

	return curve, nil
}

// NewEllipticCurveFromSecP256 从SecP256参数创建椭圆曲线
func NewEllipticCurveFromSecP256(secP256 *SecP256) *EllipticCurve {
	curve, err := NewEllipticCurve(secP256.A, secP256.B, secP256.P, secP256.N)
	if err != nil {
		panic(err)
	}
	curve.Generator = NewECPoint(curve, secP256.Gx, secP256.Gy)
	return curve
}

// GetSecretKey 生成私钥
func (e *EllipticCurve) GetSecretKey() (*big.Int, error) {
	// 生成一个长度为p的比特长度+17的随机数
	bitLength := e.P.BitLen() + 17
	key, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(bitLength)))
	if err != nil {
		return nil, err
	}
	return key, nil
}

// GetPublicKey 生成公钥
func (e *EllipticCurve) GetPublicKey() (*ECPoint, error) {
	if e.Generator == nil {
		return nil, errors.New("请设置生成点")
	}

	secretKey, err := e.GetSecretKey()
	if err != nil {
		return nil, err
	}

	return e.Generator.Multiply(secretKey), nil
}

// GetPublicKeyWithPrk 使用给定的私钥生成公钥
func (e *EllipticCurve) GetPublicKeyWithPrk(prk *big.Int) (*ECPoint, error) {
	if e.Generator == nil {
		return nil, errors.New("请设置生成点")
	}
	return e.Generator.Multiply(prk), nil
}

// IsSingular 检查曲线是否为奇异曲线 (4a³ + 27b² != 0)
func (e *EllipticCurve) IsSingular() bool {
	aa := new(big.Int).Exp(e.A, big.NewInt(3), nil) // a³
	bb := new(big.Int).Exp(e.B, big.NewInt(2), nil) // b²

	term1 := new(big.Int).Mul(aa, COEFA) // 4a³
	term2 := new(big.Int).Mul(bb, COEFB) // 27b²

	result := new(big.Int).Add(term1, term2)
	result.Mod(result, e.P) // (4a³ + 27b²) mod p

	return result.Cmp(big.NewInt(0)) == 0
}

// OnCurve 检查点是否在曲线上
func (e *EllipticCurve) OnCurve(q *ECPoint) bool {
	if q.IsZero {
		return true
	}

	ySquare := new(big.Int).Exp(q.Y, big.NewInt(2), e.P) // y² mod p
	xCube := new(big.Int).Exp(q.X, big.NewInt(3), e.P)   // x³ mod p

	// x³ + ax + b mod p
	dum := new(big.Int).Add(xCube, new(big.Int).Mul(e.A, q.X))
	dum.Add(dum, e.B)
	dum.Mod(dum, e.P)

	return ySquare.Cmp(dum) == 0
}

func (e *EllipticCurve) String() string {
	return "y^2 = x^3 + " + e.A.String() + "x + " + e.B.String() + " (mod " + e.P.String() + ")"
}

// 辅助函数：模逆运算
func modInverse(a, mod *big.Int) (*big.Int, error) {
	gcd := new(big.Int).GCD(nil, nil, a, mod)
	if gcd.Cmp(big.NewInt(1)) != 0 {
		return nil, errors.New("逆元不存在")
	}

	result := new(big.Int).ModInverse(a, mod)
	return result, nil
}

// 辅助函数：检查一个数是否可能是质数
func isProbablePrime(n *big.Int, iterations int) bool {
	return n.ProbablyPrime(iterations)
}
