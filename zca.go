package main

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

type issuer struct {
	key  crypto.Signer
	cert *x509.Certificate
}

var GlobalForm = ""
var GlobalCert = ""

func main() {
	//if first run create config file and directories
	CkConfig()
	app := qt.NewQApplication(os.Args)
	window := mkWindow(app)
	mkgui(app, window)
	qt.QApplication_Exec()
}
func mkWindow(app *qt.QApplication) *qt.QMainWindow {
	window := qt.NewQMainWindow(nil)
	window.SetMinimumSize2(850, 600)
	window.SetWindowTitle("ZCA")
	return window
}

func mkgui(app *qt.QApplication, window *qt.QMainWindow) *qt.QVBoxLayout {

	centralWidget := qt.NewQWidget(nil)
	window.SetCentralWidget(centralWidget)

	verticalLayout := qt.NewQVBoxLayout(nil)

	toolbar := toolbarInit(app, window, qt.NewQToolBar3(), centralWidget)

	verticalLayout.SetMenuBar(toolbar.QWidget)

	centralWidget.SetLayout(verticalLayout.QLayout)
	if GlobalForm == "" {
		//show home view
		verticalLayout.AddLayout(home(app, window).QLayout)
	}
	// make the window visible
	window.Show()

	return verticalLayout
}

func toolbarInit(app *qt.QApplication, window *qt.QMainWindow, toolbar *qt.QToolBar, centralWidget *qt.QWidget) *qt.QToolBar {
	toolbar.SetToolButtonStyle(qt.ToolButtonTextOnly)
	toolbar.SetMovable(true)

	label := qt.NewQLabel3("Select Action")
	toolbar.AddWidget(label.QWidget)

	newSelector := qt.NewQComboBox(nil)
	items := []string{"Select Action", "List CA Certificates", "List Certificates", "New Root", "New Intermediate", "New Server/Client", "New User Certificate", "Config"}
	newSelector.AddItems(items)
	fmt.Println(GlobalForm)
	newSelector.SetCurrentText(GlobalForm)

	newSelector.OnCurrentTextChanged(func(text string) {

		if text == "Select Action" || GlobalForm == "" {
			//show home view
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(home(app, window).QLayout)
		}

		if text == "New User Certificate" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showUserForm(app, window).QLayout)

		}
		if text == "List CA Certificates" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(listCACerts(app, window).QLayout)
		}
		if text == "List Certificates" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(listCerts(app, window).QLayout)

		}
		if text == "New Root" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showCAForm(app, window).QLayout)

		}
		if text == "New Intermediate" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showCAIntForm(app, window).QLayout)

		}

		if text == "New Server/Client" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showServerForm(app, window).QLayout)

		}

		if text == "Config" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(configEdit(app, window).QLayout)

		}

		window.Show()

	})
	toolbar.AddWidget(newSelector.QWidget)

	return toolbar
}
