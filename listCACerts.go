package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func listCACerts(app *qt.QApplication, window *qt.QMainWindow) *qt.QVBoxLayout {
	config := ConfigParser()
	verticalLayout := qt.NewQVBoxLayout(nil)

	treeWidget := qt.NewQTreeWidget(nil)

	verticalLayout.AddWidget(treeWidget.QWidget)
	treeWidget.SetColumnCount(3)
	treeWidget.SetObjectName(*qt.NewQAnyStringView3("treewidget"))
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetSortingEnabled(true)
	treeWidget.SortByColumn(0, qt.SortOrder(0))
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

		treewidgetItem := qt.NewQTreeWidgetItem2([]string{val, notBefore, notAfter})
		treewidgetItem.SetData(0, int(qt.UserRole), qt.NewQVariant11(val))
		treeWidget.AddTopLevelItem(treewidgetItem)
	}
	treeWidget.OnDoubleClicked(func(index *qt.QModelIndex) {
		certName := treeWidget.CurrentItem().Text(0)
		showCert(certName, "", "root", config, app)

	})
	treeWidget.ResizeColumnToContents(1)
	treeWidget.ResizeColumnToContents(2)
	treeWidget.SetColumnWidth(0, 400)

	treeWidget.OnContextMenuEvent(func(super func(event *qt.QContextMenuEvent), event *qt.QContextMenuEvent) {
		certName := treeWidget.CurrentItem().Text(0)
		contextMenu2(certName, config, window, app, event)
	})
	return verticalLayout

}

func readDates2(file string) (string, string) {
	// Read and parse the PEM certificate file
	pemData, err := os.ReadFile(file)
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

func contextMenu2(certName string, config ZcaConfig, w *qt.QMainWindow, app *qt.QApplication, event *qt.QContextMenuEvent) {
	menu := qt.NewQMenu(w.QWidget)

	menu.AddActionWithText("View Certificate Info").OnTriggered(func() {
		showCert(certName, "", "root", config, app)
	})
	menu.AddActionWithText("View Certificate and Key").OnTriggered(func() {
		showCertKey(certName, "", "root", config, app)
	})

	menu.ExecWithPos(event.GlobalPos().ToPointF().ToPoint())

}
