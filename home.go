package main

import qt "github.com/mappu/miqt/qt6"

func home(app *qt.QApplication, window *qt.QMainWindow) *qt.QVBoxLayout {
	config := ConfigParser()

	verticalLayout := qt.NewQVBoxLayout(nil)
	formLayout := qt.NewQFormLayout(nil)

	if getRootCt(config.RootDIR) == 0 {
		warningText := qt.NewQLabel(nil)
		style := "color:#FF0000; font-size: 12px;"
		warningText.SetStyleSheet(style)
		warningText.SetText("You need to create a root certificate authority before begining")
		formLayout.AddRowWithWidget(warningText.QWidget)

	}

	rootText := qt.NewQLineEdit(nil)
	rootText.SetText(config.RootDIR)
	rootText.SetFixedWidth(580)
	formLayout.AddRow3("Root Certificate Directory: ", rootText.QWidget)

	certText := qt.NewQLineEdit(nil)
	certText.SetText(config.CertDir)
	certText.SetFixedWidth(580)
	formLayout.AddRow3("Certificate Directory: ", certText.QWidget)

	ocspText := qt.NewQLineEdit(nil)
	ocspText.SetText(config.OCSP)
	ocspText.SetFixedWidth(580)
	formLayout.AddRow3("OCSP URL: ", ocspText.QWidget)

	crlText := qt.NewQLineEdit(nil)
	crlText.SetText(config.CRL)
	crlText.SetFixedWidth(580)
	formLayout.AddRow3("CRL Location: ", crlText.QWidget)

	verticalLayout.AddLayout(formLayout.QLayout)
	return verticalLayout
}
