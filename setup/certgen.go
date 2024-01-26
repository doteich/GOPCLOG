package setup

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/url"
	"os"
	"time"
)

func CreateKeyPair() error {

	_, err0 := os.Stat("./certs")

	if err0 != nil {
		os.Mkdir("certs", 0777)
	}

	_, err1 := os.Stat("./certs/cert.pem")
	_, err2 := os.Stat("./certs/key.pem")

	if err1 == nil && err2 == nil {
		return nil
	}

	pk, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return err
	}

	hostName, err := os.Hostname()

	if err != nil {
		return err
	}

	u, err := url.Parse(hostName)

	if err != nil {
		return err
	}

	ca := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Country: []string{"DE"}, Organization: []string{"Guanaco"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * 3650 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign | x509.KeyUsageContentCommitment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		URIs:                  []*url.URL{u},
	}

	bArr, err := x509.CreateCertificate(rand.Reader, &ca, &ca, &pk.PublicKey, pk)

	if err != nil {
		return err
	}
	cert, err := os.Create("./certs/cert.pem")

	if err != nil {
		return err
	}

	defer cert.Close()

	if err := pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: bArr}); err != nil {
		return err
	}

	key, err := os.Create("./certs/key.pem")

	if err != nil {
		return err
	}

	defer key.Close()

	if err := pem.Encode(key, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)}); err != nil {
		return err
	}

	return nil
}
