package certx

import (
	"crypto/x509"
	"testing"
)

func TestParse(t *testing.T) {
	// 替换为你的证书路径（如 ./server.crt）
	cert, err := ParseCertificate("testdata/server.crt")
	if err != nil {
		t.Fatalf("解析证书失败: %v", err)
	}

	// 输出证书关键信息（识别结果）
	t.Log("证书信息:")
	t.Logf("  版本: %d", cert.Version)
	t.Logf("  序列号: %x", cert.SerialNumber)
	t.Logf("  颁发者: %s", cert.Issuer.String())
	t.Logf("  主体: %s", cert.Subject.String())
	t.Logf("  有效期开始: %s", cert.NotBefore.Format("2006-01-02 15:04:05"))
	t.Logf("  有效期结束: %s", cert.NotAfter.Format("2006-01-02 15:04:05"))
	t.Logf("  公钥算法: %s", cert.PublicKeyAlgorithm.String())
	t.Logf("  签名算法: %s", cert.SignatureAlgorithm.String())
}

func TestGenerateSelfSignedCert(t *testing.T) {
	certPath := "testdata/server.crt"
	keyPath := "testdata/server.key"
	err := GenerateSelfSignedCert(certPath, keyPath)
	if err != nil {
		t.Fatalf("生成自签名证书失败: %v", err)
	}
	t.Log("生成自签名证书成功")
}

func TestVerifyCertificate(t *testing.T) {
	// 加载根证书池（示例：加载系统根证书）
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		t.Fatalf("加载根证书池失败: %v", err)
	}

	// 加载待验证的证书
	cert, err := ParseCertificate("testdata/server.crt")
	if err != nil {
		t.Fatalf("解析证书失败: %v", err)
	}

	// 验证证书
	err = VerifyCertificate(cert, rootCAs)
	if err != nil {
		t.Fatalf("证书验证失败: %v", err)
	}

	t.Log("证书验证成功")
}

func TestGenerateRootCA(t *testing.T) {
	certPath := "testdata/ca.crt"
	keyPath := "testdata/ca.key"
	err := GenerateRootCA(certPath, keyPath)
	if err != nil {
		t.Fatalf("生成 RootCA 失败: %v", err)
	}
	t.Log("生成 RootCA 成功")
}
