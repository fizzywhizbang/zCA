package main

import (
	"fmt"
	"strconv"

	"github.com/therecipe/qt/widgets"
)

func showCAIntForm(app *widgets.QApplication, window *widgets.QMainWindow) *widgets.QFormLayout {
	config := ConfigParser()
	formLayout := widgets.NewQFormLayout(nil)
	formLayout.SetFieldGrowthPolicy(widgets.QFormLayout__ExpandingFieldsGrow)
	label := widgets.NewQLabel2("Create New Intermediate Certificate Authority", nil, 0)
	formLayout.AddWidget(label)
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

	CN := widgets.NewQLineEdit(nil)
	CN.SetPlaceholderText("Widgets International")
	formLayout.AddRow3("Common Name: ", CN)

	CName := widgets.NewQLineEdit(nil)
	CName.SetPlaceholderText("intermediate.widgets.com")
	formLayout.AddRow3("Certificate Name: ", CName)

	age := widgets.NewQComboBox(nil)
	age.AddItems([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"})
	formLayout.AddRow3("Age (years): ", age)

	optionGroup := widgets.NewQHBoxLayout()

	//add button
	addButton := widgets.NewQPushButton(nil)
	addButton.SetText("Add")
	optionGroup.AddWidget(addButton, 0, 0)
	//cancel button
	cancelButton := widgets.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	optionGroup.AddWidget(cancelButton, 0, 0)

	formLayout.InsertRow6(8, optionGroup)

	cancelButton.ConnectClicked(func(checked bool) {
		window.Close()
		GlobalForm = ""
		mkgui(app, window)
	})
	addButton.ConnectClicked(func(checked bool) {
		if len(CName.Text()) >= 1 && len(CN.Text()) >= 1 {
			//check if file exists
			f := config.RootDIR + "/" + CName.Text() + ".pem"
			if fileExists(f) {
				widgets.QMessageBox_Critical(nil, "error", "A CA by that name already exists", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			} else {
				//create new ca
				caKey := config.RootDIR + "/" + GlobalCert + "-key.pem"
				caCert := config.RootDIR + "/" + GlobalCert + ".pem"

				issuer, err := getIssuer(caKey, caCert)
				if err != nil {
					fmt.Println(err)
				}

				//if key created properly then create the certificate
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
				_, err = intermediate(issuer, CName.Text(), i, c, s, l, o, ou, config)
				if err != nil {
					fmt.Println("Error creating Intermediate Cert", caKey)
				} else {
					GlobalForm = ""
					window.Close()
					mkgui(app, window)
				}

			}
		} else {
			widgets.QMessageBox_Critical(nil, "error", "You must supply a common name and a Intermediate CA name", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
	})

	return formLayout
}
