package main

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
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
	app := widgets.NewQApplication(len(os.Args), os.Args)
	window := mkWindow(app)
	mkgui(app, window)
	app.Exec()
}

func mkWindow(app *widgets.QApplication) *widgets.QMainWindow {
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(850, 600)
	window.SetWindowTitle("ZCA")
	return window
}

func mkgui(app *widgets.QApplication, window *widgets.QMainWindow) *widgets.QVBoxLayout {

	centralWidget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(centralWidget)

	verticalLayout := widgets.NewQVBoxLayout()

	toolbar := toolbarInit(app, window, widgets.NewQToolBar2(nil), centralWidget)

	verticalLayout.SetMenuBar(toolbar)

	centralWidget.SetLayout(verticalLayout)
	if GlobalForm == "" {
		//show home view
		verticalLayout.AddLayout(home(app, window), 0)
	}
	// make the window visible
	window.Show()

	return verticalLayout
}

func toolbarInit(app *widgets.QApplication, window *widgets.QMainWindow, toolbar *widgets.QToolBar, centralWidget *widgets.QWidget) *widgets.QToolBar {
	toolbar.SetToolButtonStyle(core.Qt__ToolButtonTextOnly)
	toolbar.SetMovable(true)

	label := widgets.NewQLabel2("Select Action", nil, 0)
	toolbar.AddWidget(label)

	newSelector := widgets.NewQComboBox(nil)
	items := []string{"Select Action", "List CA Certificates", "List Certificates", "New Root", "New Intermediate", "New Server/Client", "New User Certificate", "Config"}
	newSelector.AddItems(items)
	fmt.Println(GlobalForm)
	newSelector.SetCurrentText(GlobalForm)

	newSelector.ConnectCurrentTextChanged(func(text string) {

		if text == "Select Action" || GlobalForm == "" {
			//show home view
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(home(app, window), 0)
		}

		if text == "New User Certificate" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showUserForm(app, window), 0)

		}
		if text == "List CA Certificates" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(listCACerts(app, window), 0)
		}
		if text == "List Certificates" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(listCerts(app, window), 0)

		}
		if text == "New Root" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showCAForm(app, window), 0)

		}
		if text == "New Intermediate" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showCAIntForm(app, window), 0)

		}

		if text == "New Server/Client" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showServerForm(app, window), 0)

		}

		if text == "Config" {
			GlobalForm = text
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(configEdit(app, window), 0)

		}

		window.Show()

	})
	toolbar.AddWidget(newSelector)

	return toolbar
}
