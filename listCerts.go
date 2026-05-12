package main

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func listCerts(app *qt.QApplication, window *qt.QMainWindow) *qt.QVBoxLayout {
	config := ConfigParser()
	verticalLayout := qt.NewQVBoxLayout(nil)

	treeWidget := qt.NewQTreeWidget(nil)

	verticalLayout.AddWidget(treeWidget.QWidget)
	treeWidget.SetColumnCount(4)
	treeWidget.SetObjectName(*qt.NewQAnyStringView3("treewidget"))
	treeWidget.Header().SetSectionsClickable(true)
	treeWidget.SetSortingEnabled(true)
	treeWidget.SortByColumn(0, qt.SortOrder(0))
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
			treewidgetItem := qt.NewQTreeWidgetItem2([]string{certs[i][5], certs[i][0], certs[i][3], certs[i][1]})
			treewidgetItem.SetData(0, int(qt.UserRole), qt.NewQVariant11(certs[i][3]))
			treeWidget.AddTopLevelItem(treewidgetItem)
		}

	}

	treeWidget.OnDoubleClicked(func(index *qt.QModelIndex) {
		certName := treeWidget.CurrentItem().Text(0)
		serial := treeWidget.CurrentItem().Text(2)
		showCert(certName, serial, "cert", config, app)

	})
	treeWidget.ResizeColumnToContents(2)
	treeWidget.ResizeColumnToContents(3)
	treeWidget.SetColumnWidth(1, 50)
	treeWidget.SetColumnWidth(0, 300)

	treeWidget.OnContextMenuEvent(func(super func(event *qt.QContextMenuEvent), event *qt.QContextMenuEvent) {
		certName := treeWidget.CurrentItem().Text(0)
		serial := treeWidget.CurrentItem().Text(2)
		contextMenu(certName, serial, config, window, app, event)
	})
	return verticalLayout

}

func readDates(file string) (string, string, string) {
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
	serial := cert.SerialNumber.String()
	return notbefore, notafter, serial

}

func contextMenu(certName, serial string, config ZcaConfig, window *qt.QMainWindow, app *qt.QApplication, event *qt.QContextMenuEvent) {
	menu := qt.NewQMenu(window.QWidget)

	menu.AddActionWithText("View Certificate Info").OnTriggered(func() {
		showCert(certName, serial, "cert", config, app)
	})
	menu.AddActionWithText("View Certificate and Key").OnTriggered(func() {
		showCertKey(certName, serial, "cert", config, app)
	})
	menu.AddActionWithText("Revoke Certificate").OnTriggered(func() {
		revoke(serial, config)

	})
	menu.Exec3(event.GlobalPos().ToPointF().ToPoint(), nil)

}
