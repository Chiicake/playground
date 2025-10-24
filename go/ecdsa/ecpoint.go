package elliptic

import "math/big"

// ECPoint 表示椭圆曲线上的点
type ECPoint struct {
	Mother    *EllipticCurve
	X, Y      *big.Int
	IsZero    bool
	fastCache []*ECPoint
	cache     []*ECPoint
}

// Add 点加法
func (p *ECPoint) Add(q *ECPoint) *ECPoint {
	if !p.HasCommonMother(q) {
		panic("不是同一个椭圆曲线")
	}

	if p.IsZero {
		return q
	} else if q.IsZero {
		return p
	}

	y1, y2 := p.Y, q.Y
	x1, x2 := p.X, q.X
	pMod := p.Mother.P

	// 计算斜率alpha
	var alpha *big.Int

	if x2.Cmp(x1) == 0 {
		if y2.Cmp(y1) != 0 {
			// 对称点相加，结果为零元素
			return NewZeroECPoint(p.Mother)
		} else {
			// 点加倍，alpha = (3x² + a) / (2y)
			x1Squared := new(big.Int).Exp(x1, TWO, pMod)
			numerator := new(big.Int).Mul(x1Squared, THREE)
			numerator.Add(numerator, p.Mother.A)
			numerator.Mod(numerator, pMod)

			denominator := new(big.Int).Mul(y1, TWO)
			denominator.Mod(denominator, pMod)
			denominatorInv, err := modInverse(denominator, pMod)
			if err != nil {
				panic(err)
			}

			alpha = new(big.Int).Mul(numerator, denominatorInv)
			alpha.Mod(alpha, pMod)
		}
	} else {
		// 不同点相加，alpha = (y2 - y1) / (x2 - x1)
		numerator := new(big.Int).Sub(y2, y1)
		numerator.Mod(numerator, pMod)

		denominator := new(big.Int).Sub(x2, x1)
		denominator.Mod(denominator, pMod)
		denominatorInv, err := modInverse(denominator, pMod)
		if err != nil {
			panic(err)
		}

		alpha = new(big.Int).Mul(numerator, denominatorInv)
		alpha.Mod(alpha, pMod)
	}

	// 计算结果点坐标
	x3 := new(big.Int).Exp(alpha, TWO, pMod)
	x3.Sub(x3, x2)
	x3.Sub(x3, x1)
	x3.Mod(x3, pMod)

	y3 := new(big.Int).Sub(x1, x3)
	y3.Mul(y3, alpha)
	y3.Sub(y3, y1)
	y3.Mod(y3, pMod)

	// 确保结果为正数
	if y3.Sign() < 0 {
		y3.Add(y3, pMod)
	}
	if x3.Sign() < 0 {
		x3.Add(x3, pMod)
	}

	return NewECPoint(p.Mother, x3, y3)
}

// Multiply 点乘法（ scalar multiplication）
func (p *ECPoint) Multiply(coef *big.Int) *ECPoint {
	result := NewZeroECPoint(p.Mother)
	coefBytes := coef.Bytes()

	if p.fastCache != nil {
		// 使用FastCache加速计算
		for _, b := range coefBytes {
			result = result.Times256().Add(p.fastCache[int(b)&0xff])
		}
		return result
	}

	// 如果Cache为空，则先计算Cache内容
	if p.cache == nil {
		p.cache = make([]*ECPoint, 16)
		p.cache[0] = NewZeroECPoint(p.Mother)
		for i := 1; i < len(p.cache); i++ {
			p.cache[i] = p.cache[i-1].Add(p)
		}
	}

	for _, b := range coefBytes {
		result = result.Times16().Add(p.cache[int(b>>4)&0x0f])
		result = result.Times16().Add(p.cache[int(b)&0x0f])
	}

	return result
}

// Times16 计算点的16倍（通过连续4次加倍）
func (p *ECPoint) Times16() *ECPoint {
	result := p
	for i := 0; i < 4; i++ {
		result = result.Add(result)
	}
	return result
}

// Times256 计算点的256倍（通过连续8次加倍）
func (p *ECPoint) Times256() *ECPoint {
	result := p
	for i := 0; i < 8; i++ {
		result = result.Add(result)
	}
	return result
}

func (p *ECPoint) String() string {
	return "(" + p.X.String() + ", " + p.Y.String() + ")"
}

// HasCommonMother 检查两个点是否属于同一椭圆曲线
func (p *ECPoint) HasCommonMother(q *ECPoint) bool {
	// 简单比较曲线参数是否相同
	return p.Mother.A.Cmp(q.Mother.A) == 0 &&
		p.Mother.B.Cmp(q.Mother.B) == 0 &&
		p.Mother.P.Cmp(q.Mother.P) == 0
}

// NewECPoint 创建一个新的椭圆曲线上的点
func NewECPoint(mother *EllipticCurve, x, y *big.Int) *ECPoint {
	point := &ECPoint{
		Mother: mother,
		X:      new(big.Int).Set(x),
		Y:      new(big.Int).Set(y),
		IsZero: false,
	}

	if !mother.OnCurve(point) {
		panic("点不在椭圆曲线上")
	}

	return point
}

// NewZeroECPoint 创建一个零元素点
func NewZeroECPoint(mother *EllipticCurve) *ECPoint {
	return &ECPoint{
		Mother: mother,
		X:      big.NewInt(0),
		Y:      big.NewInt(0),
		IsZero: true,
	}
}

// FastCache 预计算缓存以加速乘法运算
func (p *ECPoint) FastCache() {
	if p.fastCache == nil {
		p.fastCache = make([]*ECPoint, 256)
		p.fastCache[0] = NewZeroECPoint(p.Mother)
		for i := 1; i < len(p.fastCache); i++ {
			p.fastCache[i] = p.fastCache[i-1].Add(p)
		}
	}
}
