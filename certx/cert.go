package certx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

// 解析 PEM 格式证书，提取关键信息
func ParseCertificate(pemPath string) (*x509.Certificate, error) {
	// 读取 PEM 文件内容
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		return nil, fmt.Errorf("读取证书失败: %v", err)
	}

	// 解码 PEM 格式（提取 DER 数据）
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("无效的 PEM 证书格式")
	}

	// 解析 DER 格式证书为 x509.Certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %v", err)
	}

	return cert, nil
}

// 生成 RootCA
func GenerateRootCA(certPath, keyPath string) error {
	// 1. 生成 RSA 私钥（4096 位）
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("生成私钥失败: %v", err)
	}

	// 2. 准备根 CA 证书模板
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("生成序列号失败: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "My Root CA",
			Organization: []string{"My Org"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour), // 有效期 10 年
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
	}

	// 3. 生成根 CA 证书（自签名）
	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	if err != nil {
		return fmt.Errorf("生成根 CA 证书失败: %v", err)
	}

	// 4. 保存根 CA 证书（PEM 格式）
	_ = os.MkdirAll(filepath.Dir(certPath), 0755)
	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("创建根 CA 证书文件失败: %v", err)
	}
	defer certFile.Close()
	if err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("编码根 CA 证书失败: %v", err)
	}

	// 5. 保存根 CA 私钥（PEM 格式）
	_ = os.MkdirAll(filepath.Dir(keyPath), 0755)
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("创建根 CA 私钥文件失败: %v", err)
	}
	defer keyFile.Close()
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return fmt.Errorf("编码根 CA 私钥失败: %v", err)
	}

	return nil
}

// 生成自签名证书和私钥
func GenerateSelfSignedCert(certPath, keyPath string) error {
	// 1. 生成 RSA 私钥（2048 位）
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("生成私钥失败: %v", err)
	}

	// 2. 准备证书模板
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return fmt.Errorf("生成序列号失败: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "example.com", // 证书绑定的域名
			Organization: []string{"My Org"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 有效期 1 年
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, // 用于服务器认证
		BasicConstraintsValid: true,
	}

	// 3. 生成证书（自签名：用私钥签名自身）
	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,             // 证书模板
		&template,             // 颁发者（自签名时与模板相同）
		&privateKey.PublicKey, // 公钥
		privateKey,            // 签名用的私钥
	)
	if err != nil {
		return fmt.Errorf("生成证书失败: %v", err)
	}

	// 4. 保存证书（PEM 格式）
	_ = os.MkdirAll(filepath.Dir(certPath), 0755)
	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("创建证书文件失败: %v", err)
	}
	defer certFile.Close()
	if err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("编码证书失败: %v", err)
	}

	// 5. 保存私钥（PEM 格式）
	_ = os.MkdirAll(filepath.Dir(keyPath), 0755)
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("创建私钥文件失败: %v", err)
	}
	defer keyFile.Close()
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey) // 转换为 PKCS#1 格式
	if err := pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes}); err != nil {
		return fmt.Errorf("编码私钥失败: %v", err)
	}

	return nil
}

// 验证证书有效性（需传入根证书池）
func VerifyCertificate(cert *x509.Certificate, rootCAs *x509.CertPool) error {
	opts := x509.VerifyOptions{
		Roots:       rootCAs,       // 信任的根证书池
		CurrentTime: time.Now(),    // 验证当前时间是否在有效期内
		DNSName:     "example.com", // 验证证书绑定的域名（若有）
	}

	// 验证证书链
	chains, err := cert.Verify(opts)
	if err != nil {
		return fmt.Errorf("证书验证失败: %v", err)
	}

	fmt.Println("证书验证成功，有效链数:", len(chains))
	return nil
}
