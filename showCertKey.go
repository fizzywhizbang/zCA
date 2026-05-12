package main

import (
	"fmt"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func showCertKey(file, serial, certtype string, config ZcaConfig, app *qt.QApplication) {

	fileCert := config.CertDir + "/" + serial + "/" + file + ".pem"
	fileKey := config.CertDir + "/" + serial + "/" + file + "-key.pem"

	if certtype == "root" {
		fileCert = config.RootDIR + "/" + file + ".pem"
		fileKey = config.RootDIR + "/" + file + "-key.pem"
	}

	cert, err := os.ReadFile(fileCert)
	if err != nil {
		fmt.Println("unable to read %s", fileCert)
	}
	key, err := os.ReadFile(fileKey)
	if err != nil {
		fmt.Println("unable to read %s", fileKey)
	}

	// Read and parse the PEM certificate file

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

	textBoxCert := qt.NewQTextEdit(nil)
	textBoxCert.SetText(string(cert))
	verticalLayout.AddWidget(textBoxCert.QWidget)

	textBoxKey := qt.NewQTextEdit(nil)
	textBoxKey.SetText(string(key))
	verticalLayout.AddWidget(textBoxKey.QWidget)

	centralWidget.SetLayout(verticalLayout.QLayout)

	// make the window visible
	window.Show()

}
