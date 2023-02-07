package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func showCertKey(file, certtype string, config ZcaConfig, app *widgets.QApplication) {

	fileCert := config.CertDir + "/" + file + "/" + file + ".pem"
	fileKey := config.CertDir + "/" + file + "/" + file + "-key.pem"

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
	window.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		if int32(e.Key()) == int32(core.Qt__Key_Escape) {
			//close window
			window.Close()
		}
	})
	centralWidget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(centralWidget)
	verticalLayout := widgets.NewQVBoxLayout()

	textBoxCert := widgets.NewQTextEdit(nil)
	textBoxCert.SetText(string(cert))
	verticalLayout.AddWidget(textBoxCert, 0, 0)

	textBoxKey := widgets.NewQTextEdit(nil)
	textBoxKey.SetText(string(key))
	verticalLayout.AddWidget(textBoxKey, 0, 0)

	centralWidget.SetLayout(verticalLayout)

	// make the window visible
	window.Show()

}
