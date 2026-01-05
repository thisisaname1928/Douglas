package testsvr

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// generate self sign cert using ecdsa eliptic, bruh this isn't work bro
func genSelfSignCert() ([]byte, []byte) {
	private, e := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if e != nil {
		fmt.Println("GEN CERT FAILED!")
	}

	serialNumber, e := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if e != nil {
		fmt.Println("GEN CERT FAILED! SẺIAL")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"BLAH BLAH DOUGLAS"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, e := x509.CreateCertificate(rand.Reader, &template, &template, &private.PublicKey, private)
	if e != nil {
		fmt.Println("GEN CERT FAILED! DẺ")
	}

	cert := new(bytes.Buffer)
	key := new(bytes.Buffer)

	pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	b, e := x509.MarshalECPrivateKey(private)
	if e != nil {
		fmt.Println("GEN CERT FAILED! Marshal ECP")
	}

	pem.Encode(key, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})

	return cert.Bytes(), key.Bytes()
}

func genECDHKey() *ecdsa.PrivateKey {
	private, e := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if e != nil {
		fmt.Println("GEN KEY FAILED!")
	}

	return private
}
