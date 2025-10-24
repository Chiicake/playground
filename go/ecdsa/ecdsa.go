package elliptic

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"
)

// ECDSASign 椭圆曲线签名算法(ECDSASign)
func ECDSASign(curve *EllipticCurve, privateKey *big.Int, message []byte) (*big.Int, *big.Int, error) {
	// 1. 检查私钥有效性
	if privateKey.Cmp(big.NewInt(0)) <= 0 || privateKey.Cmp(curve.N) >= 0 {
		return nil, nil, errors.New("invalid private key")
	}

	// 2. 生成随机数k (1 < k < n)
	k, err := rand.Int(rand.Reader, new(big.Int).Sub(curve.N, big.NewInt(1)))
	if err != nil {
		return nil, nil, err
	}
	// 确保k >= 1
	if k.Cmp(big.NewInt(0)) == 0 {
		k = big.NewInt(1)
	}

	// 3. 计算点R = k * G
	R := curve.Generator.Multiply(k)

	// 4. 计算r = R.x mod n
	r := new(big.Int).Mod(R.X, curve.N)
	if r.Cmp(big.NewInt(0)) == 0 {
		return nil, nil, errors.New("r is zero, need to generate new k")
	}

	// 5. 计算消息哈希 e = H(message)
	hash := sha256.Sum256(message)
	e := new(big.Int).SetBytes(hash[:])
	e.Mod(e, curve.N) // 将哈希值规约到模n下

	// 6. 计算k的逆元 k^-1 mod n
	kInv := new(big.Int).ModInverse(k, curve.N)

	// 7. 计算s = k^-1 * (e + d * r) mod n
	dotProduct := new(big.Int).Mul(privateKey, r) // d*r
	s := new(big.Int).Add(e, dotProduct)          // e + d*r
	s.Mul(s, kInv)                                // (e + d*r) * k^-1
	s.Mod(s, curve.N)                             // 模n运算

	if s.Cmp(big.NewInt(0)) == 0 {
		return nil, nil, errors.New("s is zero, need to generate new k")
	}

	return r, s, nil
}

// ECDSAVerify 椭圆曲线验证算法(ECDSASign)
func ECDSAVerify(curve *EllipticCurve, publicKey *ECPoint, message []byte, r, s *big.Int) (bool, error) {
	// 1. 检查r和s的有效性
	if r.Cmp(big.NewInt(0)) <= 0 || r.Cmp(curve.N) >= 0 {
		return false, errors.New("invalid r value")
	}
	if s.Cmp(big.NewInt(0)) <= 0 || s.Cmp(curve.N) >= 0 {
		return false, errors.New("invalid s value")
	}

	// 2. 检查公钥有效性 (公钥必须在曲线上且不是无穷远点)
	if publicKey == nil || publicKey.X == nil || publicKey.Y == nil {
		return false, errors.New("invalid public key")
	}

	// 3. 计算消息哈希 e = H(message)
	hash := sha256.Sum256(message)
	e := new(big.Int).SetBytes(hash[:])
	e.Mod(e, curve.N) // 将哈希值规约到模n下

	// 4. 计算w = s^-1 mod n
	w := new(big.Int).ModInverse(s, curve.N)

	// 5. 计算u1 = e * w mod n 和 u2 = r * w mod n
	u1 := new(big.Int).Mul(e, w)
	u1.Mod(u1, curve.N)

	u2 := new(big.Int).Mul(r, w)
	u2.Mod(u2, curve.N)

	// 6. 计算点P = u1*G + u2*Q (Q是公钥)
	P1 := curve.Generator.Multiply(u1)
	P2 := publicKey.Multiply(u2)
	P := P1.Add(P2)

	// 7. 如果P是无穷远点，验证失败
	if P.X == nil || P.Y == nil {
		return false, nil
	}

	// 8. 计算v = P.x mod n
	v := new(big.Int).Mod(P.X, curve.N)

	// 9. 验证v是否等于r
	return v.Cmp(r) == 0, nil
}
