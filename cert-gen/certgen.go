package certgen

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func GeneratePEMFiles() {
	ca := createCA()
	caCertPEM, caPrivateKeyPEM := createCACert(ca)
	createCertificate(ca, caCertPEM, caPrivateKeyPEM)

}

func createCA() *x509.Certificate {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1000),
		Subject: pkix.Name{
			Country:      []string{"DE"},
			Organization: []string{"TestCompany"},
			Locality:     []string{"Leipzig"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), //today  + 10 years validity period
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	return ca
}

func createCACert(ca *x509.Certificate) (*bytes.Buffer, *bytes.Buffer) {
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4069)

	if err != nil {
		fmt.Println("CA - Private Key generation failed")
		fmt.Println(err)
	}

	caCertificateBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)

	if err != nil {
		fmt.Println("CA - Cert creation failed")
		fmt.Println(err)
	}

	caPEM := new(bytes.Buffer)

	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertificateBytes,
	})

	caPrivateKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey),
	})

	return caPEM, caPrivateKeyPEM
}

func createCertificate(ca *x509.Certificate, caCertPEM *bytes.Buffer, caPrivateKeyPEM *bytes.Buffer) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1001),
		Subject: pkix.Name{
			Country:      []string{"DE"},
			Organization: []string{"TestCompany"},
			Locality:     []string{"Leipzig"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println("New Cert - RSA Key Gen Failed")
		fmt.Println(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, certPrivKey)
	if err != nil {
		fmt.Println("New Cert generation failed")
		fmt.Println(err)
	}

	certPEMFile, err := os.Create("./certs/cert.pem")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	certPEMFile.Close()

	pem.Encode(certPEMFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivateKey, err := os.Create("./certs/private_key.pem")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pem.Encode(certPrivateKey, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	certPrivateKey.Close()
}
