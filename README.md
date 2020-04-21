Diaverse Hardware
---
## About 
This repository contains the code used to power the testing hardware for the diaverse project. The hardware is responsible for interacting with the VUI under test and offers a small GUI to allow the user to start and stop test scripts. It uses the 
Google Cloud Platform API's for speech to text and text to speech and a low cost raspberry pi.  

## How to Run
   If you want to run this project you must have the following, 
+ A Raspberry Pi running Raspbian
+ A Logitech c270 webcam
+ a simple USB speaker or 3.5mm speaker


## Steps 

### API Setup


First you need to setup a google cloud platform account [here](https://cloud.google.com/)
 
Then you need to enable the API's for speech to text and text to speech.
 
After enabling the API's you need to create a new service account with access permission to the API's. 
Be sure to generate a .json credentials key while creating the service account as it will be needed later.  


### Hardware Setup
 
 --- 
 flash your RPi with the latest version of raspbian and install all updates. 

 Currently the project is configured for the default 'pi' user, if you wish to use another user please change 
 the paths within the codebase. 
 
 

Clone the repository in your home directory then run 

`cd Hardware && ./setup`

Then wait for the setup script to install all the required dependencies and GCP CLI. At the end of the installation script
you will be asked for your Google Cloud Platform login information. 


Then create a new directory in the home directory of the pi user named `credentials`.


Finally, edit your .bashrc to include the following 

`export GOOGLE_APPLICATION_CREDENTIALS=/home/pi/credentials/<NAME OF CREDENTIAL FILE>`


At this point you should be good to go, you have the API's enabled, the keys downloaded and exported, and the dependencies installed.  

To start the hardware server go the following commands `go build ./... && go build *.go && ./main`

At this point the server should be listening on port 8080, you can view the web UI by navigating to `http://<IP Address Of your PI>:8080`. 

## Logging in 
To use the test hardware login with your given username and harware token. You should then see a list of all test scripts assosiated with your user.
