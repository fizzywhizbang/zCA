package main

import (
	"fmt"
	"strconv"

	qt "github.com/mappu/miqt/qt6"
)

func showUserForm(app *qt.QApplication, window *qt.QMainWindow) *qt.QFormLayout {
	//userCertificate(iss *issuer, cn string, y int, O, OU []string, config ZcaConfig)
	config := ConfigParser()
	formLayout := qt.NewQFormLayout(nil)
	formLayout.SetFieldGrowthPolicy(qt.QFormLayout__ExpandingFieldsGrow)
	label := qt.NewQLabel3("Create New Server or Client Certificate")
	formLayout.AddWidget(label.QWidget)

	selector := qt.NewQComboBox(nil)
	selector.AddItems(getCas())
	formLayout.AddRow3("Select Certificate Authority: ", selector.QWidget)

	O := qt.NewQLineEdit(nil)
	O.SetPlaceholderText("Widgets INTL")
	formLayout.AddRow3("Organization: ", O.QWidget)

	OU := qt.NewQLineEdit(nil)
	OU.SetPlaceholderText("Widgets Web Services")
	formLayout.AddRow3("Organizational Unit: ", OU.QWidget)

	// CN := qt.NewQLineEdit(nil)
	// CN.SetPlaceholderText("Widgets International")
	// formLayout.AddRow3("Common Name: ", CN)

	CName := qt.NewQLineEdit(nil)
	CName.SetPlaceholderText("john.doe")
	formLayout.AddRow3("Certificate Name: ", CName.QWidget)

	Email := qt.NewQLineEdit(nil)
	Email.SetPlaceholderText("john.doe@gmail.com")
	formLayout.AddRow3("Email Address: ", Email.QWidget)

	age := qt.NewQComboBox(nil)
	age.AddItems([]string{"3", "2", "1"})
	formLayout.AddRow3("Age (years): ", age.QWidget)

	// domains := qt.NewQTextEdit(nil)
	// domains.SetMaximumHeight(100)
	// formLayout.AddRow3("Domains: ", domains)

	// ips := qt.NewQTextEdit(nil)
	// ips.SetMaximumHeight(100)
	// formLayout.AddRow3("IPs: ", ips)

	optionGroup := qt.NewQHBoxLayout(nil)

	//add button
	addButton := qt.NewQPushButton(nil)
	addButton.SetText("Add")
	optionGroup.AddWidget(addButton.QWidget)
	//cancel button
	cancelButton := qt.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	optionGroup.AddWidget(cancelButton.QWidget)

	formLayout.InsertRow6(11, optionGroup.QLayout)

	cancelButton.OnClicked(func() {
		window.Close()
		GlobalForm = ""
		mkgui(app, window)
	})
	addButton.OnClicked(func() {
		//check if file exists
		caKey := config.RootDIR + "/" + selector.CurrentText() + "-key.pem"
		caCert := config.RootDIR + "/" + selector.CurrentText() + ".pem"

		issuer, err := getIssuer(caKey, caCert)
		if err != nil {
			fmt.Println(err)
		}

		o := []string{}
		o = append(o, O.Text())
		ou := []string{}
		ou = append(ou, OU.Text())
		i, _ := strconv.Atoi(age.CurrentText())
		email := []string{}
		email = append(email, Email.Text())

		// domainSlice := []string{}
		// scanner := bufio.NewScanner(strings.NewReader(domains.ToPlainText()))
		// for scanner.Scan() {
		// 	domainSlice = append(domainSlice, scanner.Text())
		// }
		// ipSlice := []string{}
		// scanner2 := bufio.NewScanner(strings.NewReader(ips.ToPlainText()))
		// for scanner.Scan() {
		// 	ipSlice = append(ipSlice, scanner2.Text())
		// }

		//sign(iss *issuer, cn string, y int, domains, ipAddresses []string)
		_, err = userCertificate(issuer, CName.Text(), i, o, ou, email, config)

		if err != nil {
			fmt.Println("Error creating Cert", err)
		} else {
			GlobalForm = ""
			window.Close()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(listCerts(app, window).QLayout)
			GlobalForm = "List Certificates"
		}
		// }
	})

	return formLayout
}
