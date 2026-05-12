package main

import qt "github.com/mappu/miqt/qt6"

func configEdit(app *qt.QApplication, window *qt.QMainWindow) *qt.QVBoxLayout {
	config := ConfigParser()

	verticalLayout := qt.NewQVBoxLayout(nil)
	formLayout := qt.NewQFormLayout(nil)
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

	optionGroup := qt.NewQHBoxLayout(nil)

	//add button
	savButton := qt.NewQPushButton(nil)
	savButton.SetText("Save")
	savButton.OnClicked(func() {
		if updateConfig(rootText.Text(), certText.Text(), ocspText.Text(), crlText.Text()) {
			qt.QMessageBox_Information(nil, "Saved", "Updated Configuration File\nBe sure to restart the program to apply the new settings")
		} else {
			qt.QMessageBox_Warning(nil, "Warning", "Something went wrong, I need my meds")
		}
	})
	optionGroup.AddWidget(savButton.QWidget)

	//cancel button
	cancelButton := qt.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	cancelButton.OnClicked(func() {
		window.Close()
		GlobalForm = ""
		mkgui(app, window)
	})
	optionGroup.AddWidget(cancelButton.QWidget)
	formLayout.InsertRow6(5, optionGroup.QLayout)

	verticalLayout.AddLayout(formLayout.QLayout)
	return verticalLayout
}
