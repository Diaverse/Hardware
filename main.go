package main

import (
	"./controller"
	"github.com/prometheus/common/log"
	"net/http"
)

//This file simply starts the local server
func main() {
	StartServer()
}

func StartServer() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/", controller.ServeWebpage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/user", controller.ServeUsersWebPage)
	mux.HandleFunc("/stopScript", controller.StopTestScriptHandler)
	mux.HandleFunc("/executeScript", controller.ExecuteTestScriptHandler)
	mux.HandleFunc("/authenticateToken", controller.CheckForExsistingHardwareToken)
	mux.HandleFunc("/logout", controller.ServeLogoutPage)
	controller.InitilizeStructs()
	controller.InitilizeScriptStopChan()
	log.Info("Starting server on port 7080")
	log.Fatal(http.ListenAndServe(":7080", mux))

}
