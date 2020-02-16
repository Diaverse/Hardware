package main

import (
	"./controller"
	"github.com/prometheus/common/log"
	"net/http"
)

//This file simply starts the local server
func main() ***REMOVED***
	mux := http.DefaultServeMux

	mux.HandleFunc("/", controller.ServeLogin)
	mux.HandleFunc("/listScripts", controller.ServeScriptListView)
	mux.HandleFunc("/executeScript", controller.ExecuteTestScript)

	log.Fatal(http.ListenAndServe(":8080", mux))
***REMOVED***
