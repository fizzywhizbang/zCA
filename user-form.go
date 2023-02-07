package main

import (
	"fmt"
	"strconv"

	"github.com/therecipe/qt/widgets"
)

func showUserForm(app *widgets.QApplication, window *widgets.QMainWindow) *widgets.QFormLayout {
	//userCertificate(iss *issuer, cn string, y int, O, OU []string, config ZcaConfig)
	config := ConfigParser()
	formLayout := widgets.NewQFormLayout(nil)
	formLayout.SetFieldGrowthPolicy(widgets.QFormLayout__ExpandingFieldsGrow)
	label := widgets.NewQLabel2("Create New Server or Client Certificate", nil, 0)
	formLayout.AddWidget(label)

	selector := widgets.NewQComboBox(nil)
	selector.AddItems(getCas())
	formLayout.AddRow3("Select Certificate Authority: ", selector)

	O := widgets.NewQLineEdit(nil)
	O.SetPlaceholderText("Widgets INTL")
	formLayout.AddRow3("Organization: ", O)

	OU := widgets.NewQLineEdit(nil)
	OU.SetPlaceholderText("Widgets Web Services")
	formLayout.AddRow3("Organizational Unit: ", OU)

	// CN := widgets.NewQLineEdit(nil)
	// CN.SetPlaceholderText("Widgets International")
	// formLayout.AddRow3("Common Name: ", CN)

	CName := widgets.NewQLineEdit(nil)
	CName.SetPlaceholderText("marc.levine")
	formLayout.AddRow3("Certificate Name: ", CName)

	age := widgets.NewQComboBox(nil)
	age.AddItems([]string{"3", "2", "1"})
	formLayout.AddRow3("Age (years): ", age)

	// domains := widgets.NewQTextEdit(nil)
	// domains.SetMaximumHeight(100)
	// formLayout.AddRow3("Domains: ", domains)

	// ips := widgets.NewQTextEdit(nil)
	// ips.SetMaximumHeight(100)
	// formLayout.AddRow3("IPs: ", ips)

	optionGroup := widgets.NewQHBoxLayout()

	//add button
	addButton := widgets.NewQPushButton(nil)
	addButton.SetText("Add")
	optionGroup.AddWidget(addButton, 0, 0)
	//cancel button
	cancelButton := widgets.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	optionGroup.AddWidget(cancelButton, 0, 0)

	formLayout.InsertRow6(10, optionGroup)

	cancelButton.ConnectClicked(func(checked bool) {
		window.Close()
		GlobalForm = ""
		mkgui(app, window)
	})
	addButton.ConnectClicked(func(checked bool) {
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
		_, err = userCertificate(issuer, CName.Text(), i, o, ou, config)

		if err != nil {
			fmt.Println("Error creating Cert", err)
		} else {
			GlobalForm = ""
			window.Close()
			vlayout2 := mkgui(app, window)
			vlayout2.AddLayout(listCerts(app, window), 0)
			GlobalForm = "LC"
		}
		// }
	})

	return formLayout
}
