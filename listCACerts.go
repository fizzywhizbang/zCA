package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func listCACerts(app *widgets.QApplication, window *widgets.QMainWindow) *widgets.QVBoxLayout {
	config := ConfigParser()
	verticalLayout := widgets.NewQVBoxLayout()

	treeWidget := widgets.NewQTreeWidget(nil)

	verticalLayout.AddWidget(treeWidget, 0, 0)
	treeWidget.SetColumnCount(3)
	treeWidget.SetObjectName("treewidget")
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetSortingEnabled(true)
	treeWidget.SortByColumn(0, core.Qt__SortOrder(0))
	treeWidget.SetAlternatingRowColors(true)
	treeWidget.HorizontalScrollBar().SetHidden(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	treeWidget.SetStyleSheet(tableColors)
	treeWidget.Header()
	treeWidget.SetHeaderLabels([]string{"Certificate", "Not Before", "Not After"})

	certs := getCerts(config.RootDIR)

	for _, val := range certs {
		file := config.RootDIR + "/" + val + ".pem"
		fmt.Println(file)
		notBefore, notAfter, _ := readDates(file)
		// notBefore := ""
		// notAfter := ""

		treewidgetItem := widgets.NewQTreeWidgetItem2([]string{val, notBefore, notAfter}, 0)
		treewidgetItem.SetData(0, int(core.Qt__UserRole), core.NewQVariant12(val))
		treeWidget.AddTopLevelItem(treewidgetItem)
	}
	treeWidget.ConnectDoubleClicked(func(index *core.QModelIndex) {
		certName := treeWidget.CurrentItem().Text(0)
		showCert(certName, "root", config, app)

	})
	treeWidget.ResizeColumnToContents(1)
	treeWidget.ResizeColumnToContents(2)
	treeWidget.SetColumnWidth(0, 400)

	treeWidget.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		certName := treeWidget.CurrentItem().Text(0)
		contextMenu2(certName, config, window, app, event)
	})
	return verticalLayout

}

func readDates2(file string) (string, string) {
	// Read and parse the PEM certificate file
	pemData, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	block, rest := pem.Decode([]byte(pemData))
	if block == nil || len(rest) > 0 {
		log.Fatal("Certificate decoding error")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	notbefore := cert.NotBefore.Format("Jan 2 15:04:05 2006 MST")
	notafter := cert.NotAfter.Format("Jan 2 15:04:05 2006 MST")
	return notbefore, notafter

}

func contextMenu2(certName string, config ZcaConfig, w *widgets.QMainWindow, app *widgets.QApplication, event *gui.QContextMenuEvent) {
	menu := widgets.NewQMenu(w)

	menu.AddAction("View Certificate Info").ConnectTriggered(func(checked bool) {
		showCert(certName, "root", config, app)
	})
	menu.AddAction("View Certificate and Key").ConnectTriggered(func(checked bool) {
		showCertKey(certName, "root", config, app)
	})

	menu.Exec2(event.GlobalPos().QPoint_PTR(), nil)

}
