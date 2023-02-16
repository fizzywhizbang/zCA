read -p "Enter the name of the CA to add to your trust store:" CA
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "root/${CA}.pem"