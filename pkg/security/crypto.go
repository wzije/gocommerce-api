package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/ecommerce-api/pkg/config"
	errors2 "github.com/ecommerce-api/pkg/exception"
	log "github.com/sirupsen/logrus"
	"os"
)

type CryptoUtil interface {
	GenerateKey()
	GetAuthorizedKey(key string) (*rsa.PublicKey, error)
	GetPrivateKey(key string) (*rsa.PrivateKey, error)
	Encrypt(payload string, key *rsa.PublicKey) (string, error)
	Decrypt(cipherText string, key *rsa.PrivateKey) (string, error)
}

type cryptoUtil struct {
}

func (e cryptoUtil) GenerateKey() {

	targetLocation := config.BasePath + ".keys/"

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println("Error generating RSA private key:", err)
		os.Exit(1)
	}

	// Encode the private key to the PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyFile, err := os.Create(targetLocation + "private.pem")

	if err != nil {
		fmt.Println("Error creating private key file:", err)
		os.Exit(1)
	}

	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return
	}

	if err := privateKeyFile.Close(); err != nil {
		return
	}

	// Extract the public key from the private key
	publicKey := &privateKey.PublicKey

	// Encode the public key to the PEM format
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}

	publicKeyFile, err := os.Create(targetLocation + "public.pem")

	if err != nil {
		fmt.Println("Error creating public key file:", err)
		os.Exit(1)
	}

	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return
	}

	if err := publicKeyFile.Close(); err != nil {
		return
	}

	fmt.Println("RSA key pair generated successfully!")
}

func (e cryptoUtil) GetPrivateKey(key string) (*rsa.PrivateKey, error) {
	fileKey := fmt.Sprintf("%s/.keys/%s", config.BasePath, key)
	keyData, err := os.ReadFile(fileKey)
	if err != nil {
		log.Printf("ERROR: fail get idrsa, %s", err.Error())
		os.Exit(1)
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		log.Printf("ERROR: fail get idrsa, invalid key")
		os.Exit(1)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		log.Printf("ERROR: fail get idrsa, %s", err.Error())
		os.Exit(1)
	}

	return privateKey, nil
}

func (e cryptoUtil) GetAuthorizedKey(key string) (*rsa.PublicKey, error) {
	fileKey := fmt.Sprintf("%s/.keys/authorized/%s", config.BasePath, key)

	keyData, err := os.ReadFile(fileKey)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(keyData)

	publicKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)

	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		return publicKey, nil
	default:
		return nil, errors2.ErrParsing
	}

	//return publicKey, err
}

func (e cryptoUtil) Encrypt(payload string, key *rsa.PublicKey) (string, error) {
	// params
	msg := []byte(payload)
	rnd := rand.Reader
	hash := sha256.New()

	// encrypt with OAEP
	cipherText, err := rsa.EncryptOAEP(hash, rnd, key, msg, nil)

	if err != nil {
		log.Printf("ERROR: fail to encrypt, %s", err.Error())
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (e cryptoUtil) Decrypt(cipherText string, key *rsa.PrivateKey) (string, error) {
	// decode base64 encoded signature
	msg, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		log.Printf("ERROR: fail to base64 decode, %s", err.Error())
		return "", err
	}

	// params
	rnd := rand.Reader
	hash := sha256.New()

	// decrypt with OAEP
	plainText, err := rsa.DecryptOAEP(hash, rnd, key, msg, nil)
	if err != nil {
		log.Printf("ERROR: fail to decrypt, %s", err.Error())
		return "", err
	}

	return string(plainText), nil

}

func NewCryptoUtil() CryptoUtil {
	return cryptoUtil{}
}
