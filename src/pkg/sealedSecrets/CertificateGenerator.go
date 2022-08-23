package sealedSecrets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"bytes"
	utils "gepaplexx/day-x-generator/pkg/util"
	"os"
)

func GenerateCertificate(env string) (string, string, error) {
	utils.PrintAction("Creating private key and cert...")
	_, err := os.Stat(fmt.Sprintf("%s/%s.crt", utils.TARGET_DIR, env))
	if err == nil {
		utils.PrintSuccess()
		fmt.Println("found existing certificate for env", env)
		crt, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.crt", utils.TARGET_DIR, env))
		if err != nil {
			return "", "", err
		}

		key, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.key", utils.TARGET_DIR, env))
		if err != nil {
			return "", "", err
		}

		return certToString(crt, key)
	}

	publicKey, privateKey, err := createKeyAndCert(env)
	if err != nil {
		utils.PrintFailure()
		return "", "", err
	}
	utils.PrintSuccess()

	vaultPassword, err := utils.ReadFromStdin("Enter vault password")
	if err != nil {
		return "", "", err
	}

	utils.PrintAction("\nEncrypting certificate for vault...")
	encryptedPublicKey, err := encryptForVault(publicKey, vaultPassword)
	if err != nil {
		utils.PrintFailure()
		return "", "", err
	}
	utils.PrintSuccess()

	utils.PrintAction("Encrypting private key for vault...")
	encryptedPrivateKey, err := encryptForVault(privateKey, vaultPassword)
	if err != nil {
		utils.PrintFailure()
		return "", "", err
	}
	utils.PrintSuccess()

	return encryptedPublicKey, encryptedPrivateKey, nil
}

func createKeyAndCert(env string) (string, string, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2022),
		Subject: pkix.Name{
			Organization: []string{"Gepaplexx"},
			Country:      []string{"AT"},
			Locality:     []string{"Linz"},
			PostalCode:   []string{"4030"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		panic(err)
	}

	certificatePath := fmt.Sprintf("%s/%s.crt", utils.TARGET_DIR, env)
	caPEM, err := os.Create(certificatePath)
	if err != nil {
		panic(err)
	}

	pem.Encode(caPEM, buildPemBlock("CERTIFICATE", caBytes))

	caPrivKeyPEM, err := os.Create(fmt.Sprintf("%s/%s.key", utils.TARGET_DIR, env))
	if err != nil {
		panic(err)
	}

	pem.Encode(caPrivKeyPEM, buildPemBlock("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(caPrivKey)))

	return certToString(caBytes, x509.MarshalPKCS1PrivateKey(caPrivKey))
}

func certToString(crt []byte, key []byte) (string, string, error) {
	var PublicKey bytes.Buffer
	var PrivateKey bytes.Buffer
	err := pem.Encode(&PublicKey, buildPemBlock("CERTIFICATE", crt))
	if err != nil {
		return "", "", err
	}

	err = pem.Encode(&PrivateKey, buildPemBlock("RSA PRIVATE KEY", key))
	if err != nil {
		return "", "", err
	}

	return PublicKey.String(), PrivateKey.String(), nil
}

func buildPemBlock(typ string, bytes []byte) *pem.Block {
	return &pem.Block{
		Type:  typ,
		Bytes: bytes,
	}
}
