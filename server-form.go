package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/therecipe/qt/widgets"
)

func showServerForm(app *widgets.QApplication, window *widgets.QMainWindow) *widgets.QFormLayout {
	formLayout := widgets.NewQFormLayout(nil)
	formLayout.SetFieldGrowthPolicy(widgets.QFormLayout__ExpandingFieldsGrow)

	C := widgets.NewQLineEdit(nil)
	C.SetPlaceholderText("US")
	formLayout.AddRow3("Country (two letter): ", C)

	S := widgets.NewQLineEdit(nil)
	S.SetPlaceholderText("CA")
	formLayout.AddRow3("Province/State (two letter): ", S)

	L := widgets.NewQLineEdit(nil)
	L.SetPlaceholderText("Los Angeles")
	formLayout.AddRow3("Locality: ", L)

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
	CName.SetPlaceholderText("root.wigets.com")
	formLayout.AddRow3("Certificate Name: ", CName)

	age := widgets.NewQComboBox(nil)
	age.AddItems([]string{"2", "1"})
	formLayout.AddRow3("Age (years): ", age)

	domains := widgets.NewQTextEdit(nil)
	domains.SetMaximumHeight(100)
	formLayout.AddRow3("Domains: ", domains)

	ips := widgets.NewQTextEdit(nil)
	ips.SetMaximumHeight(100)
	formLayout.AddRow3("IPs: ", ips)

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
		gui(app)
	})
	addButton.ConnectClicked(func(checked bool) {
		//check if file exists
		//create new ca
		caKey := "root/" + GlobalCert + "-key.pem"
		caCert := "root/" + GlobalCert + ".pem"

		issuer, err := getIssuer(caKey, caCert)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(issuer)

		c := []string{}
		c = append(c, C.Text())
		s := []string{}
		s = append(s, S.Text())
		l := []string{}
		l = append(l, L.Text())
		o := []string{}
		o = append(o, O.Text())
		ou := []string{}
		ou = append(ou, OU.Text())
		i, _ := strconv.Atoi(age.CurrentText())

		domainSlice := []string{}
		scanner := bufio.NewScanner(strings.NewReader(domains.ToPlainText()))
		for scanner.Scan() {
			domainSlice = append(domainSlice, scanner.Text())
		}
		ipSlice := []string{}
		scanner2 := bufio.NewScanner(strings.NewReader(ips.ToPlainText()))
		for scanner.Scan() {
			ipSlice = append(ipSlice, scanner2.Text())
		}

		//sign(iss *issuer, cn string, y int, domains, ipAddresses []string)
		_, err = sign(issuer, CName.Text(), i, domainSlice, ipSlice, c, s, l, o, ou)

		// 	_, err = makeRootCert(key, caCert, CN.Text(), c, s, l, o, ou, i)
		if err != nil {
			fmt.Println("Error creating Cert", err)
		} else {
			GlobalForm = ""
			window.Close()
			gui(app)
		}
		// }
	})

	return formLayout
}
