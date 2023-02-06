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

	selector := widgets.NewQComboBox(nil)
	selector.AddItems(getCas())
	toolbar.AddWidget(selector)

	toolbar.AddSeparator()

	goButton := widgets.NewQPushButton2("New CA", nil)
	toolbar.AddWidget(goButton)
	goButton.ConnectClicked(func(checked bool) {
		if GlobalForm != "CA" {
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showCAForm(app, window), 0)
			GlobalForm = "CA"
		}

		window.Show()
	})

	serverCert := widgets.NewQPushButton2("New Server Cert", nil)
	toolbar.AddWidget(serverCert)
	serverCert.ConnectClicked(func(checked bool) {
		GlobalCert = selector.CurrentText()
		if GlobalForm != "SC" {
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(showServerForm(app, window), 0)
			GlobalForm = "SC"
		}
	})

	showCerts := widgets.NewQPushButton2("List Certs", nil)
	toolbar.AddWidget(showCerts)
	showCerts.ConnectClicked(func(checked bool) {
		if GlobalForm != "LC" {
			centralWidget.DeleteLater()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(listCerts(app, window), 0)
			GlobalForm = "LC"
		}
	})

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
