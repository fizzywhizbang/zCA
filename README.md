# ZCA
the Zoe Certificate Authority  

## About
I just started building this thing using code from:  
https://github.com/jsha/minica  
https://github.com/grantae/certinfo  

Really just throwing a graphical interface on tools that exist in go (lipstick on a pig if you will)  

## How it works
Install go  
Install therecipe/qt  
Download this source  
type go run .  
If using a mac download the latest release (I will do my best to update this every time I release new code)
  
The program will create a root directory and a crt directory
The program will create a config file that can be modified within the program or with a text editor  
Your root and intermediate certs will be stored under root and any other certs you create will be stored under crt  
You can have more than one root CA if you deisre  
  
I am still working on this and this program will be able to create multiple types of certificates  
OCSP repsonder will be separate but built to read from this program if you need one
