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
	mux.HandleFunc("/", controller.ServeWebpage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/user", controller.ServeUsersWebPage)
	mux.HandleFunc("/register", controller.RegisterHardware)
	mux.HandleFunc("/executeScript", controller.ExecuteTestScriptHandler)
	mux.HandleFunc("/authenticateToken", controller.CheckForExsistingHardwareToken)

	controller.InitilizeStructs()
	log.Info("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
