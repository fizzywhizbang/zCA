package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"time"
)

func makeDirectory(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeIssuer(keyFile, certFile, caname string, C, S, L, O, OU []string, y int) error {
	//make key file and pass the key to makeRootCert
	// keyFile = "root/" + keyFile
	key, err := makeKey(keyFile)
	if err != nil {
		return err
	}
	// certFile = "root/" + certFile
	_, err = makeRootCert(key, certFile, caname, C, S, L, O, OU, y)
	if err != nil {
		return err
	}
	return nil
}

func makeKey(filename string) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	der := x509.MarshalPKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = pem.Encode(file, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	})
	if err != nil {
		return nil, err
	}
	return key, nil
}

func makeRootCert(key crypto.Signer, filename, caname string, C, S, L, O, OU []string, y int) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	skid, err := calculateSKID(key.Public())
	if err != nil {
		return nil, err
	}
	template := &x509.Certificate{
		Subject: pkix.Name{
			Country:            C,
			Province:           S,
			Locality:           L,
			Organization:       O,
			OrganizationalUnit: OU,
			CommonName:         caname,
		},
		SerialNumber: serial,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(y, 0, 0),

		SubjectKeyId:          skid,
		AuthorityKeyId:        skid,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = pem.Encode(file, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(der)
}

func intermediate(iss *issuer, cn string, y int, C, S, L, O, OU []string, config ZcaConfig) (*x509.Certificate, error) {
	key, err := makeKey(fmt.Sprintf("%s/%s-key.pem", config.RootDIR, cn))
	if err != nil {
		return nil, err
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	template := &x509.Certificate{
		Subject: pkix.Name{
			Country:            C,
			Province:           S,
			Locality:           L,
			Organization:       O,
			OrganizationalUnit: OU,
			CommonName:         cn,
		},
		SerialNumber: serial,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(y, 0, 0),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, iss.cert, key.Public(), iss.key)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%s.pem", config.RootDIR, cn), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = pem.Encode(file, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(der)
}

func sign(iss *issuer, cn string, y int, domains, ipAddresses, C, S, L, O, OU []string, config ZcaConfig) (*x509.Certificate, error) {
	var cnFolder = config.CertDir + "/" + strings.Replace(cn, "*", "_", -1)
	err := os.Mkdir(cnFolder, 0700)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	key, err := makeKey(fmt.Sprintf("%s/%s-key.pem", cnFolder, cn))
	if err != nil {
		return nil, err
	}
	parsedIPs, err := parseIPs(ipAddresses)
	if err != nil {
		return nil, err
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	template := &x509.Certificate{
		DNSNames:    domains,
		IPAddresses: parsedIPs,
		Subject: pkix.Name{
			Country:            C,
			Province:           S,
			Locality:           L,
			Organization:       O,
			OrganizationalUnit: OU,
			CommonName:         cn,
		},
		SerialNumber: serial,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(y, 0, 0),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, iss.cert, key.Public(), iss.key)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%scert.pem", cnFolder, cn), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = pem.Encode(file, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(der)
}

func userCertificate(iss *issuer, cn string, y int, domains, ipAddresses, C, S, L, O, OU []string, config ZcaConfig) (*x509.Certificate, error) {
	var cnFolder = config.CertDir + "/" + strings.Replace(cn, "*", "_", -1)
	err := os.Mkdir(cnFolder, 0700)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	key, err := makeKey(fmt.Sprintf("%s/%s-key.pem", cnFolder, cn))
	if err != nil {
		return nil, err
	}
	parsedIPs, err := parseIPs(ipAddresses)
	if err != nil {
		return nil, err
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	template := &x509.Certificate{
		DNSNames:    domains,
		IPAddresses: parsedIPs,
		Subject: pkix.Name{
			Country:            C,
			Province:           S,
			Locality:           L,
			Organization:       O,
			OrganizationalUnit: OU,
			CommonName:         cn,
		},
		SerialNumber: serial,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(y, 0, 0),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection, x509.ExtKeyUsageCodeSigning},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, iss.cert, key.Public(), iss.key)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%scert.pem", cnFolder, cn), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = pem.Encode(file, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(der)
}
