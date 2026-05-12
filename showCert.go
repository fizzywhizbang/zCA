package main

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func showCert(certName, serial, certtype string, config ZcaConfig, app *qt.QApplication) {
	file := config.CertDir + "/" + serial + "/" + certName + ".pem"
	if certtype == "root" {
		file = config.RootDIR + "/" + certName + ".pem"
	}
	// Read and parse the PEM certificate file
	pemData, err := os.ReadFile(file)
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
	window.OnKeyPressEvent(func(super func(*qt.QKeyEvent), event *qt.QKeyEvent) {
		if int32(event.Key()) == int32(qt.Key_Escape) {
			//close window
			window.Close()
		}
	})
	centralWidget := qt.NewQWidget(nil)
	window.SetCentralWidget(centralWidget)
	verticalLayout := qt.NewQVBoxLayout(nil)

	textBox := qt.NewQTextEdit(nil)
	text, _ := CertificateText(cert)
	textBox.SetText(text)
	verticalLayout.AddWidget(textBox.QWidget)
	centralWidget.SetLayout(verticalLayout.QLayout)

	// make the window visible
	window.Show()

}
