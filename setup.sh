#!/bin/bash


#update the pi and its installed software
sudo apt-get update -y && sudo apt-get upgrade -y


#Install gstreamer and its plugins
sudo apt-get install libgstreamer1.0-0 gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-doc gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-pulseaudio


#install pulse audio
sudo apt-get install -y pulseaudio pulseaudio-utils

#install golang
sudo apt-get install -y golang

#install our audio player
sudo apt-get install -y mplayer


#install our ftp server
sudo apt-get install -y proftpd

#update everything
sudo apt-get update -y && sudo apt-get upgrade -y


#install google cloud command line tool
sudo echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list

sudo apt-get install -y apt-transport-https ca-certificates

sudo curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -

sudo apt-get update && sudo apt-get install -y google-cloud-sdk


#import the API's we need
clear
echo "installing text to speech api"
cd src/LPLibs

go get cloud.google.com/go/texttospeech/apiv1

go get google.golang.org/genproto/googleapis/cloud/texttospeech/v1


clear
echo "installing speech to text api"
cd ../scripts

go get cloud.google.com/go/speech/apiv1

go get google.golang.org/genproto/googleapis/cloud/speech/v1


clear


echo "building project"

go build livecaption.go



cd ../LPLibs

go build *.go

cd ../tests/ChatterBox

go build *.go

cd ../../..
clear

#setup gcloud

sudo gcloud init
