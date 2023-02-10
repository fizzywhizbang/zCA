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

func listCerts(app *widgets.QApplication, window *widgets.QMainWindow) *widgets.QVBoxLayout {
	config := ConfigParser()
	verticalLayout := widgets.NewQVBoxLayout()

	treeWidget := widgets.NewQTreeWidget(nil)

	verticalLayout.AddWidget(treeWidget, 0, 0)
	treeWidget.SetColumnCount(4)
	treeWidget.SetObjectName("treewidget")
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetSortingEnabled(true)
	treeWidget.SortByColumn(0, core.Qt__SortOrder(0))
	treeWidget.SetAlternatingRowColors(true)
	treeWidget.HorizontalScrollBar().SetHidden(true)
	tableColors := "alternate-background-color: #88DD88; background-color:#FFFFFF; color:#000000; font-size: 12px;"
	treeWidget.SetStyleSheet(tableColors)
	treeWidget.Header()
	treeWidget.SetHeaderLabels([]string{"Certificate", "Status", "Serial", "Not After"})

	/*
		column0 (status): Valid Revoked or Expired (V,R,E)
		column1 (currentTime + y): Expiration time
		column2: revokation time if R is set
		column3: Serial number (use serial number)
		column4: filename of the certificate (use filename)
		column5: certificate subject name (use CN)
	*/
	certs := readCRL(config.CRL)
	for i := 0; i < len(certs); i++ {
		if len(certs[i]) > 0 {
			treewidgetItem := widgets.NewQTreeWidgetItem2([]string{certs[i][5], certs[i][0], certs[i][3], certs[i][1]}, 0)
			treewidgetItem.SetData(0, int(core.Qt__UserRole), core.NewQVariant12(certs[i][3]))
			treeWidget.AddTopLevelItem(treewidgetItem)
		}

	}

	treeWidget.ConnectDoubleClicked(func(index *core.QModelIndex) {
		certName := treeWidget.CurrentItem().Text(0)
		serial := treeWidget.CurrentItem().Text(2)
		showCert(certName, serial, "cert", config, app)

	})
	treeWidget.ResizeColumnToContents(2)
	treeWidget.ResizeColumnToContents(3)
	treeWidget.SetColumnWidth(1, 50)
	treeWidget.SetColumnWidth(0, 300)

	treeWidget.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		certName := treeWidget.CurrentItem().Text(0)
		serial := treeWidget.CurrentItem().Text(2)
		contextMenu(certName, serial, config, window, app, event)
	})
	return verticalLayout

}

func readDates(file string) (string, string, string) {
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
	serial := cert.SerialNumber.String()
	return notbefore, notafter, serial

}

func contextMenu(certName, serial string, config ZcaConfig, w *widgets.QMainWindow, app *widgets.QApplication, event *gui.QContextMenuEvent) {
	menu := widgets.NewQMenu(w)

	menu.AddAction("View Certificate Info").ConnectTriggered(func(checked bool) {
		showCert(certName, serial, "cert", config, app)
	})
	menu.AddAction("View Certificate and Key").ConnectTriggered(func(checked bool) {
		showCertKey(certName, serial, "cert", config, app)
	})
	menu.AddAction("Revoke Certificate").ConnectTriggered(func(checked bool) {
		fmt.Println(serial)
	})
	menu.Exec2(event.GlobalPos().QPoint_PTR(), nil)

}
