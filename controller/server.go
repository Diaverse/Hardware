package controller

import (
	"../domain"
	"github.com/prometheus/common/log"
	"html/template"
	"net/http"
)

//This is the server controller. It contains all server handlers.

//ServeLogin listens on the root directory of the hosted web page and serves the needed html to the user
func ServeLogin(w http.ResponseWriter, r *http.Request) ***REMOVED***
	//todo; define a way for a user to remain logged in across sessions
	//create a new login page struct
	page := domain.LoginPage***REMOVED***
		Title:     "Diaverse Login Screen",
		AuthToken: "",
		Content:   "<html><p>Hello</p></html>",
	***REMOVED***
	t, err := template.ParseFiles("templates/login.html")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	//parse the login template and serve
	t.Execute(w, page)
***REMOVED***

func ServeScriptListView(w http.ResponseWriter, r *http.Request) ***REMOVED***
	page := domain.ListPage***REMOVED***
		Title:          "Diaverse Script View",
		ScriptList:     []string***REMOVED***"Script 1", "Script 2", "Script 3"***REMOVED***,
		SelectedScript: "This one",
	***REMOVED***

	t, err := template.ParseFiles("templates/ScriptList.html")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	//parse the login template and serve
	t.Execute(w, page)
***REMOVED***

//GatherTestScripts is a listener attached to the web UI. It queries the
//page for the authentication token input text area and then displays the list of test scripts for that user
func GatherTestScripts(w http.ResponseWriter, r *http.Request) ***REMOVED***

***REMOVED***
