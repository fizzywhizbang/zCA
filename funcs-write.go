package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
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

func mkLog(status string, y int, serial, filename, cname string, config ZcaConfig) {

	file, err := os.OpenFile(config.CRL, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Could not open CRL")
		return
	}

	defer file.Close()

	/*
		tab delimited always 6 columns openssl format
			column0 (status): Valid Revoked or Expired (V,R,E)
			column1 (currentTime + y): Expiration time
			column2: revokation time if R is set
			column3: Serial number (use serial number)
			column4: filename of the certificate (use filename)
			column5: certificate subject name (use CN)
			Not all columns are used but a tab is there just in case
			0   1(YYMMDD 00:00:01) 2  3                                   4       5
	*/
	status = strings.ToUpper(status)
	currentTime := time.Now()
	currentTime = currentTime.AddDate(y, 0, 0)
	notBefore := currentTime.Format("060102000001")

	revocationTime := ""
	if status == "R" {
		revocationTime = time.Now().Format("060102000001")
	}

	formatedString := status + "\t" + notBefore + "\t" + revocationTime + "\t" + serial + "\t" + filename + "\t" + cname + "\n"

	_, err2 := file.WriteString(formatedString)

	if err2 != nil {
		fmt.Println("Could not write text to example.txt")
	}

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
	//log the action
	mkLog("V", y, serial.String(), filename, caname, ConfigParser())

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
	//log the action
	mkLog("V", y, serial.String(), cn, cn, ConfigParser())

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
	fmt.Println("Parse IP")
	if err != nil {
		return nil, err
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	fmt.Println("Gen Serial")
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
		OCSPServer:            []string{config.OCSP},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	fmt.Println("Create Certificate")
	der, err := x509.CreateCertificate(rand.Reader, template, iss.cert, key.Public(), iss.key)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%s.pem", cnFolder, cn), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
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

	//log the action
	mkLog("V", y, serial.String(), cn, cn, ConfigParser())

	return x509.ParseCertificate(der)
}

func userCertificate(iss *issuer, cn string, y int, O, OU, email []string, config ZcaConfig) (*x509.Certificate, error) {
	var cnFolder = config.CertDir + "/" + strings.Replace(cn, "*", "_", -1)
	err := os.Mkdir(cnFolder, 0700)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	key, err := makeKey(fmt.Sprintf("%s/%s-key.pem", cnFolder, cn))
	if err != nil {
		return nil, err
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	template := &x509.Certificate{
		Subject: pkix.Name{
			Organization:       O,
			OrganizationalUnit: OU,
			CommonName:         cn,
		},
		SerialNumber:          serial,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(y, 0, 0),
		EmailAddresses:        email,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection, x509.ExtKeyUsageClientAuth},
		OCSPServer:            []string{config.OCSP},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	extSubjectAltName := pkix.Extension{}
	extSubjectAltName.Id = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 6, 6}
	extSubjectAltName.Critical = false
	extSubjectAltName.Value = []byte(`Principal Name:levine.marc.s`)
	template.ExtraExtensions = []pkix.Extension{extSubjectAltName}

	der, err := x509.CreateCertificate(rand.Reader, template, iss.cert, key.Public(), iss.key)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%s.pem", cnFolder, cn), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
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
	// log the action
	mkLog("V", y, serial.String(), cn, cn, ConfigParser())

	return x509.ParseCertificate(der)
}
