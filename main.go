package main

import (
	"./controller"
	"github.com/prometheus/common/log"
	"net/http"
)

//This file simply starts the local server
func main() ***REMOVED***
	StartServer()
***REMOVED***

func StartServer() ***REMOVED***
	mux := http.DefaultServeMux
	mux.HandleFunc("/", controller.ServeLogin)
	mux.HandleFunc("/listScripts", controller.ServeScriptListView)
	mux.HandleFunc("/register", controller.RegisterHardware)
	mux.HandleFunc("/executeScript", controller.ExecuteTestScriptHandler)

	log.Info("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

***REMOVED***
