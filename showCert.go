package main

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func showCert(certName, certtype string, config ZcaConfig, app *widgets.QApplication) {
	file := config.CertDir + "/" + certName + "/" + certName + "cert.pem"
	if certtype == "root" {
		file = config.RootDIR + "/" + certName + ".pem"
	}
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
	window.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		if int32(e.Key()) == int32(core.Qt__Key_Escape) {
			//close window
			window.Close()
		}
	})
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
