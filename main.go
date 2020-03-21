package main

import (
	"./controller"
	"github.com/prometheus/common/log"
	"net/http"
)

//This file simply starts the local server
func main() {
	StartServer()
	//service.Recognize()
}

func StartServer() {
	mux := http.DefaultServeMux
	fs := http.FileServer(http.Dir("./templates"))
	mux.Handle("/", fs)

	mux.HandleFunc("/listScripts", controller.ServeScriptListView)
	mux.HandleFunc("/register", controller.RegisterHardware)
	mux.HandleFunc("/executeScript", controller.ExecuteTestScriptHandler)
	controller.InitilizeStructs()
	log.Info("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
