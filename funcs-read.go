package main

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func split(s string) (results []string) {
	if len(s) > 0 {
		return strings.Split(s, ",")
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getCas() []string {
	config := ConfigParser()

	files, err := os.ReadDir(config.RootDIR)
	if err != nil {
		log.Fatal(err)
	}
	list := []string{}
	for _, file := range files {
		if !strings.Contains(file.Name(), "-key") && !file.IsDir() && !strings.Contains(file.Name(), "DS_Store") {

			list = append(list, strings.Replace(file.Name(), ".pem", "", -1))
		}

	}
	return list
}

func getCertStatus(crl, serial string) string {
	readFile, err := os.Open(crl)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	/*
		column0 (status): Valid Revoked or Expired (V,R,E)
		column1 (currentTime + y): Expiration time
		column2: revokation time if R is set
		column3: Serial number (use serial number)
		column4: filename of the certificate (use filename)
		column5: certificate subject name (use CN)
	*/
	for fileScanner.Scan() {
		//split
		slice := strings.Split(fileScanner.Text(), "\t")
		if slice[3] == serial {
			return slice[0]
		}

	}

	readFile.Close()
	return "F"
}

func getCerts(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	list := []string{}
	for _, file := range files {
		if !strings.Contains(file.Name(), "DS_Store") && !strings.Contains(file.Name(), "key") && !strings.Contains(file.Name(), "txt") {

			list = append(list, strings.Replace(file.Name(), ".pem", "", -1))
		}

	}
	return list
}

func getCACerts() []string {
	config := ConfigParser()
	files, err := os.ReadDir(config.RootDIR)
	if err != nil {
		log.Fatal(err)
	}

	list := []string{}
	for _, file := range files {
		if !strings.Contains(file.Name(), "DS_Store") {

			list = append(list, strings.Replace(file.Name(), ".pem", "", -1))
		}

	}
	return list
}

func getIssuer(keyFile, certFile string) (*issuer, error) {
	keyContents, keyErr := ioutil.ReadFile(keyFile)
	certContents, certErr := ioutil.ReadFile(certFile)
	//if neither key or cert exist then 'makeIssuer' make issuer and run this function again
	if keyErr != nil {
		return nil, fmt.Errorf("%s (but %s exists)", keyErr, certFile)
		//if cert does not exist throw error because we don't have a matching pair
	} else if certErr != nil {
		return nil, fmt.Errorf("%s (but %s exists)", certErr, keyFile)
	}
	////////////////////////
	//read private key to test for errors
	key, err := readPrivateKey(keyContents)
	if err != nil {
		return nil, fmt.Errorf("reading private key from %s: %s", keyFile, err)
	}

	cert, err := readCert(certContents)
	if err != nil {
		return nil, fmt.Errorf("reading CA certificate from %s: %s", certFile, err)
	}

	equal, err := publicKeysEqual(key.Public(), cert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("comparing public keys: %s", err)
	} else if !equal {
		return nil, fmt.Errorf("public key in CA certificate %s doesn't match private key in %s",
			certFile, keyFile)
	}
	return &issuer{key, cert}, nil
}

func readPrivateKey(keyContents []byte) (crypto.Signer, error) {
	block, _ := pem.Decode(keyContents)
	if block == nil {
		return nil, fmt.Errorf("no PEM found")
	} else if block.Type != "RSA PRIVATE KEY" && block.Type != "ECDSA PRIVATE KEY" {
		return nil, fmt.Errorf("incorrect PEM type %s", block.Type)
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func readCert(certContents []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certContents)
	if block == nil {
		return nil, fmt.Errorf("no PEM found")
	} else if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("incorrect PEM type %s", block.Type)
	}
	return x509.ParseCertificate(block.Bytes)
}
func parseIPs(ipAddresses []string) ([]net.IP, error) {
	var parsed []net.IP
	for _, s := range ipAddresses {
		p := net.ParseIP(s)
		if p == nil {
			return nil, fmt.Errorf("invalid IP address %s", s)
		}
		parsed = append(parsed, p)
	}
	return parsed, nil
}

func publicKeysEqual(a, b interface{}) (bool, error) {
	aBytes, err := x509.MarshalPKIXPublicKey(a)
	if err != nil {
		return false, err
	}
	bBytes, err := x509.MarshalPKIXPublicKey(b)
	if err != nil {
		return false, err
	}
	return bytes.Compare(aBytes, bBytes) == 0, nil
}
func calculateSKID(pubKey crypto.PublicKey) ([]byte, error) {
	spkiASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}
	_, err = asn1.Unmarshal(spkiASN1, &spki)
	if err != nil {
		return nil, err
	}
	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)
	return skid[:], nil
}
