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
  
The program will create a root directory and a crt directory
Your root certs will be stored under root and any other certs you create will be stored under crt  
You can have more than one root CA if you deisre  
  
I am still working on this and this program will be able to create multiple types of certificates  
as well as being able to have root and intermediate certificate authorities  
  
