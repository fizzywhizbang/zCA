package main

import (
	"fmt"
	"strconv"

	qt "github.com/mappu/miqt/qt6"
)

func showCAForm(app *qt.QApplication, window *qt.QMainWindow) *qt.QFormLayout {
	config := ConfigParser()
	formLayout := qt.NewQFormLayout(nil)
	formLayout.SetFieldGrowthPolicy(qt.QFormLayout__ExpandingFieldsGrow)
	label := qt.NewQLabel3("Create New Root Certificate Authority")
	formLayout.AddWidget(label.QWidget)
	C := qt.NewQLineEdit(nil)
	C.SetPlaceholderText("US")
	formLayout.AddRow3("Country (two letter): ", C.QWidget)

	S := qt.NewQLineEdit(nil)
	S.SetPlaceholderText("CA")
	formLayout.AddRow3("Province/State (two letter): ", S.QWidget)

	L := qt.NewQLineEdit(nil)
	L.SetPlaceholderText("Los Angeles")
	formLayout.AddRow3("Locality: ", L.QWidget)

	O := qt.NewQLineEdit(nil)
	O.SetPlaceholderText("Widgets INTL")
	formLayout.AddRow3("Organization: ", O.QWidget)

	OU := qt.NewQLineEdit(nil)
	OU.SetPlaceholderText("Widgets Web Services")
	formLayout.AddRow3("Organizational Unit: ", OU.QWidget)

	CN := qt.NewQLineEdit(nil)
	CN.SetPlaceholderText("Widgets International")
	formLayout.AddRow3("Common Name: ", CN.QWidget)

	CName := qt.NewQLineEdit(nil)
	CName.SetPlaceholderText("root.qt.com")
	formLayout.AddRow3("Certificate Name: ", CName.QWidget)

	age := qt.NewQComboBox(nil)
	age.AddItems([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"})
	formLayout.AddRow3("Age (years): ", age.QWidget)

	optionGroup := qt.NewQHBoxLayout(nil)

	//add button
	addButton := qt.NewQPushButton(nil)
	addButton.SetText("Add")
	optionGroup.AddWidget(addButton.QWidget)
	//cancel button
	cancelButton := qt.NewQPushButton(nil)
	cancelButton.SetText("Cancel")
	optionGroup.AddWidget(cancelButton.QWidget)

	formLayout.InsertRow6(8, optionGroup.QLayout)

	cancelButton.OnClicked(func() {
		window.Close()
		GlobalForm = ""
		mkgui(app, window)
	})
	addButton.OnClicked(func() {
		if len(CName.Text()) >= 1 && len(CN.Text()) >= 1 {
			//check if file exists
			f := config.RootDIR + "/" + CName.Text() + ".pem"
			if fileExists(f) {
				qt.QMessageBox_Critical(nil, "error", "A CA by that name already exists")
			} else {
				//create new ca
				caKey := config.RootDIR + "/" + CName.Text() + "-key.pem"
				caCert := config.RootDIR + "/" + CName.Text() + ".pem"

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
						mkgui(app, window)
					}
				}

			}
		} else {
			qt.QMessageBox_Critical(nil, "error", "You must supply a common name and a root CA name")
		}
	})

	return formLayout
}
