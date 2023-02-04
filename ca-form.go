package main

import (
	"fmt"
	"strconv"

	"github.com/therecipe/qt/widgets"
)

func showCAForm(app *widgets.QApplication, window *widgets.QMainWindow) *widgets.QFormLayout {
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

	CN := widgets.NewQLineEdit(nil)
	CN.SetPlaceholderText("Widgets International")
	formLayout.AddRow3("Common Name: ", CN)

	CName := widgets.NewQLineEdit(nil)
	CName.SetPlaceholderText("root.wigets.com")
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
		gui(app, window)
	})
	addButton.ConnectClicked(func(checked bool) {
		if len(CName.Text()) >= 1 && len(CN.Text()) >= 1 {
			//check if file exists
			f := "root/" + CName.Text() + ".pem"
			if fileExists(f) {
				widgets.QMessageBox_Critical(nil, "error", "A CA by that name already exists", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			} else {
				//create new ca
				caKey := "root/" + CName.Text() + "-key.pem"
				caCert := "root/" + CName.Text() + ".pem"

				//step 1 make key
				key, err := makeKey(caKey)
				if err != nil {
					fmt.Println("Unable to generate key by the name", caKey)
				} else {
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
					_, err = makeRootCert(key, caCert, CN.Text(), c, s, l, o, ou, i)
					if err != nil {
						fmt.Println("Error creating Root Cert", caKey)
					} else {
						GlobalForm = ""
						window.Close()
						gui(app, window)
					}
				}

			}
		} else {
			widgets.QMessageBox_Critical(nil, "error", "You must supply a common name and a root CA name", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}
	})

	return formLayout
}
