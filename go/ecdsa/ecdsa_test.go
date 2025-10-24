package elliptic

import (
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	// 生成随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	// scepk1
	curve := NewEllipticCurveFromSecP256(K1)

	// 从私钥中提取公钥
	publicKey, err := curve.GetPublicKeyWithPrk(privateKey.D)
	if err != nil {
		t.Fatal(err)
	}

	// 要签名的消息
	message := []byte("Hello, World!")

	// 对消息进行签名
	r, s, err := ECDSASign(curve, privateKey.D, message)
	if err != nil {
		t.Fatalf("签名过程失败: %v", err)
	}
	if r == nil || s == nil {
		t.Fatal("生成的签名为空")
	}

	// 验证签名
	valid, err := ECDSAVerify(curve, publicKey, message, r, s)
	if err != nil {
		t.Fatalf("验证过程失败: %v", err)
	}
	if !valid {
		t.Error("验证失败: 合法签名未通过验证")
	}

	// 额外测试：验证篡改后的消息
	tamperedMessage := []byte("Hello, World?") // 篡改原始消息
	valid, err = ECDSAVerify(curve, publicKey, tamperedMessage, r, s)
	if err != nil {
		t.Fatalf("篡改消息验证失败: %v", err)
	}
	if valid {
		t.Error("验证失败: 篡改后的消息通过了验证")
	}
}
