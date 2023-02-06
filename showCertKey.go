package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/widgets"
)

func showCertKey(file string, app *widgets.QApplication) {

	fileCert := "crt/" + file + "/" + file + "cert.pem"
	fileKey := "crt/" + file + "/" + file + "-key.pem"
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
