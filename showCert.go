package main

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/therecipe/qt/widgets"
)

func showCert(certName string, config ZcaConfig, app *widgets.QApplication) {
	file := config.CertDir + "/" + certName + "/" + certName + "cert.pem"
	// Read and parse the PEM certificate file
	pemData, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	block, rest := pem.Decode([]byte(pemData))
	if block == nil || len(rest) > 0 {
		log.Fatal("Certificate decoding error")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	//create new window
	window := mkWindow(app)
	centralWidget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(centralWidget)
	verticalLayout := widgets.NewQVBoxLayout()

	textBox := widgets.NewQTextEdit(nil)
	text, _ := CertificateText(cert)
	textBox.SetText(text)
	verticalLayout.AddWidget(textBox, 0, 0)
	centralWidget.SetLayout(verticalLayout)

	// make the window visible
	window.Show()

}
