package main

import "github.com/therecipe/qt/widgets"

func home(app *widgets.QApplication, window *widgets.QMainWindow) *widgets.QVBoxLayout {
	config := ConfigParser()

	verticalLayout := widgets.NewQVBoxLayout()
	formLayout := widgets.NewQFormLayout(nil)

	if getRootCt(config.RootDIR) == 0 {
		warningText := widgets.NewQLabel(nil, 0)
		style := "color:#FF0000; font-size: 12px;"
		warningText.SetStyleSheet(style)
		warningText.SetText("You need to create a root certificate authority before begining")
		formLayout.AddRow5(warningText)

	}

	rootText := widgets.NewQLineEdit(nil)
	rootText.SetText(config.RootDIR)
	rootText.SetFixedWidth(580)
	formLayout.AddRow3("Root Certificate Directory: ", rootText)

	certText := widgets.NewQLineEdit(nil)
	certText.SetText(config.CertDir)
	certText.SetFixedWidth(580)
	formLayout.AddRow3("Certificate Directory: ", certText)

	ocspText := widgets.NewQLineEdit(nil)
	ocspText.SetText(config.OCSP)
	ocspText.SetFixedWidth(580)
	formLayout.AddRow3("OCSP URL: ", ocspText)

	crlText := widgets.NewQLineEdit(nil)
	crlText.SetText(config.CRL)
	crlText.SetFixedWidth(580)
	formLayout.AddRow3("CRL Location: ", crlText)

	verticalLayout.AddLayout(formLayout, 0)
	return verticalLayout
}
