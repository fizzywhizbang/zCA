package main

import (
	"crypto"
	"crypto/x509"
	"flag"
	"fmt"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type issuer struct {
	key  crypto.Signer
	cert *x509.Certificate
}

var GlobalForm = ""
var GlobalCert = ""

func main() {
	//first run check an create directory structure
	makeDirectory("root") //for the root certificates
	makeDirectory("crt")  //for signed certificates
	app := widgets.NewQApplication(len(os.Args), os.Args)

	gui(app)
	app.Exec()
}

func gui(app *widgets.QApplication) {

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(800, 600)
	window.SetWindowTitle("ZCA")

	centralWidget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(centralWidget)

	verticalLayout := widgets.NewQVBoxLayout()

	toolbar := toolbarInit(app, window, widgets.NewQToolBar2(nil), verticalLayout)

	verticalLayout.SetMenuBar(toolbar)

	centralWidget.SetLayout(verticalLayout)

	// make the window visible
	window.Show()
}

func toolbarInit(app *widgets.QApplication, window *widgets.QMainWindow, toolbar *widgets.QToolBar, vlayout *widgets.QVBoxLayout) *widgets.QToolBar {
	toolbar.SetToolButtonStyle(core.Qt__ToolButtonTextOnly)
	toolbar.SetMovable(true)

	label := widgets.NewQLabel2("Select CA", nil, 0)
	toolbar.AddWidget(label)

	selector := widgets.NewQComboBox(nil)
	selector.AddItems(getCas())
	toolbar.AddWidget(selector)

	toolbar.AddSeparator()

	goButton := widgets.NewQPushButton2("New CA", nil)
	toolbar.AddWidget(goButton)
	goButton.ConnectClicked(func(checked bool) {
		if GlobalForm != "CA" {
			vlayout.AddLayout(showCAForm(app, window), 0)
			GlobalForm = "CA"
		}

		window.Show()
	})

	serverCert := widgets.NewQPushButton2("New Server Cert", nil)
	toolbar.AddWidget(serverCert)
	serverCert.ConnectClicked(func(checked bool) {
		GlobalCert = selector.CurrentText()
		if GlobalForm != "SC" {
			vlayout.AddLayout(showServerForm(app, window), 0)
			GlobalForm = "SC"
		}
	})

	return toolbar
}

func main2() {

	//define flag variables
	var caName = flag.String("ca-name", "zca", "Root CA Name that will look for matching key and cert")
	// var domains = flag.String("domains", "", "Comma separated domain names to include as Server Alternative Names.")
	// var ipAddresses = flag.String("ip-addresses", "", "Comma separated IP addresses to include as Server Alternative Names.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		// fmt.Fprintf(os.Stderr, `
		// Minica is a simple CA intended for use in situations where the CA operator
		// also operates each host where a certificate will be used. It automatically
		// generates both a key and a certificate when asked to produce a certificate.
		// It does not offer OCSP or CRL services. Minica is appropriate, for instance,
		// for generating certificates for RPC systems or microservices.

		// On first run, minica will generate a keypair and a root certificate in the
		// current directory, and will reuse that same keypair and root certificate
		// unless they are deleted.

		// On each run, minica will generate a new keypair and sign an end-entity (leaf)
		// certificate for that keypair. The certificate will contain a list of DNS names
		// and/or IP addresses from the command line flags. The key and certificate are
		// placed in a new directory whose name is chosen as the first domain name from
		// the certificate, or the first IP address if no domain names are present. It
		// will not overwrite existing keys or certificates.

		// `)
		flag.PrintDefaults()
	}
	flag.Parse()

	//check for CA key and cert
	caKey := "root/" + *caName + "-key.pem"
	caCert := "root/" + *caName + ".pem"

	//if default notify
	if *caName == "zca" {
		if fileExists(caKey) && fileExists(caCert) {
			fmt.Println("Using", *caName, "as default root ca")
		} else {
			fmt.Println("You did not supply a custom name for the root CA therefore, zca will be the default")
		}

	}

	//this will create issuer if not exist and return issuer data
	// issuer, err := getIssuer(caKey, caCert, *caName)
	// if err != nil {
	// 	fmt.Println(err)
	// 	// os.Exit(1)
	// }
	// fmt.Println(issuer)

	// // if no domain or ip address quit
	// if *domains == "" && *ipAddresses == "" {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

	// //if more args than expected quit
	// if len(flag.Args()) > 0 {
	// 	fmt.Printf("Extra arguments: %s (maybe there are spaces in your domain list?)\n", flag.Args())
	// 	os.Exit(1)
	// }

	// //parse domains
	// domainSlice := split(*domains)
	// domainRe := regexp.MustCompile("^[A-Za-z0-9.*-]+$")
	// for _, d := range domainSlice {
	// 	if !domainRe.MatchString(d) {
	// 		fmt.Printf("Invalid domain name %q\n", d)
	// 		os.Exit(1)
	// 	}
	// }

	// //parse ip addresses
	// ipSlice := split(*ipAddresses)
	// for _, ip := range ipSlice {
	// 	if net.ParseIP(ip) == nil {
	// 		fmt.Printf("Invalid IP address %q\n", ip)
	// 		os.Exit(1)
	// 	}
	// }

	// //now create request and sign with issuer data
	// _, err = sign(issuer, domainSlice, ipSlice)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
