package main

import (
	"crypto"
	"crypto/x509"
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
	window.SetMinimumSize2(800, 600)
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

	// make the window visible
	window.Show()

	return verticalLayout
}

func toolbarInit(app *widgets.QApplication, window *widgets.QMainWindow, toolbar *widgets.QToolBar, centralWidget *widgets.QWidget) *widgets.QToolBar {
	toolbar.SetToolButtonStyle(core.Qt__ToolButtonTextOnly)
	toolbar.SetMovable(true)

	label := widgets.NewQLabel2("Select CA", nil, 0)
	toolbar.AddWidget(label)

	newSelector := widgets.NewQComboBox(nil)
	items := []string{"Select Action", "New Root", "New Intermediate", "New Server/Client", "List Certificates", "List CA Certificates"}
	newSelector.AddItems(items)

	newSelector.ConnectCurrentTextChanged(func(text string) {

		if text == "New Root" {
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showCAForm(app, window), 0)
			GlobalForm = "New Root"
		}
		if text == "New Intermediate" {
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showCAIntForm(app, window), 0)
			GlobalForm = text
		}

		if text == "New Server/Client" {
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showServerForm(app, window), 0)
			GlobalForm = "New Server/Client"
		}

		if text == "List Certificates" {
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(listCerts(app, window), 0)
			GlobalForm = text
		}
		if text == "List CA Certificates" {
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(listCACerts(app, window), 0)
			GlobalForm = text
		}
		window.Show()
	})
	toolbar.AddWidget(newSelector)

	// showCerts := widgets.NewQPushButton2("List Certs", nil)
	// toolbar.AddWidget(showCerts)
	// showCerts.ConnectClicked(func(checked bool) {
	// 	if GlobalForm != "LC" {
	// 		centralWidget.DeleteLater()
	// 		vlayout2 := mkgui(app, window)
	// 		vlayout2.AddLayout(listCerts(app, window), 0)
	// 		GlobalForm = "LC"
	// 	}
	// })

	showConfig := widgets.NewQPushButton2("Config", nil)
	toolbar.AddWidget(showConfig)
	showConfig.ConnectClicked(func(checked bool) {
		if GlobalForm != "C" {
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(configEdit(app, window), 0)
			GlobalForm = "C"
		}
	})

	return toolbar
}
