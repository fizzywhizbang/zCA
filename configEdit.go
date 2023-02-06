package main

import "github.com/therecipe/qt/widgets"

func configEdit(app *widgets.QApplication, window *widgets.QMainWindow) *widgets.QVBoxLayout {
	config := ConfigParser()

	verticalLayout := widgets.NewQVBoxLayout()
	formLayout := widgets.NewQFormLayout(nil)
	rootText := widgets.NewQLineEdit(nil)
	rootText.SetText(config.RootDIR)
	rootText.SetFixedWidth(580)
	formLayout.AddRow3("Root Certificate Directory: ", rootText)

	certText := widgets.NewQLineEdit(nil)
	certText.SetText(config.CertDir)
	certText.SetFixedWidth(580)
	formLayout.AddRow3("Certificate Directory: ", certText)

	optionGroup := widgets.NewQHBoxLayout()

	//add button
	savButton := widgets.NewQPushButton(nil)
	savButton.SetText("Save")
	savButton.ConnectClicked(func(checked bool) {
		if updateConfig(rootText.Text(), certText.Text()) {
			widgets.QMessageBox_Information(nil, "Saved", "Updated Configuration File\nBe sure to restart the program to apply the new settings", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		} else {
			widgets.QMessageBox_Warning(nil, "Warning", "Something went wrong, I need my meds", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
	})
	optionGroup.AddWidget(savButton, 0, 0)

	//cancel button
	cancelButton := widgets.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	cancelButton.ConnectClicked(func(checked bool) {
		window.Close()
		GlobalForm = ""
		mkgui(app, window)
	})
	optionGroup.AddWidget(cancelButton, 0, 0)
	formLayout.InsertRow6(3, optionGroup)

	verticalLayout.AddLayout(formLayout, 0)
	return verticalLayout
}
