package controller

import (
	"../domain"
	"encoding/json"
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
		Title: "Diaverse Script View",
		ScriptList: []domain.TestScript***REMOVED***
			domain.TestScript***REMOVED***
				Cases: []domain.TestCase***REMOVED***
					domain.TestCase***REMOVED***
						Responses:      []string***REMOVED***"Hello, how are you"***REMOVED***,
						ExpectedOutput: []string***REMOVED***"I am fine"***REMOVED***,
					***REMOVED***,
					domain.TestCase***REMOVED***
						Responses:      []string***REMOVED***"what are you doing?"***REMOVED***,
						ExpectedOutput: []string***REMOVED***"absolutely nothing."***REMOVED***,
					***REMOVED***,
				***REMOVED***,
				Result: true,
			***REMOVED***,
			domain.TestScript***REMOVED***
				Cases:  []domain.TestCase***REMOVED******REMOVED***,
				Result: true,
			***REMOVED***,
		***REMOVED***,
		SelectedScript: "Script One",
	***REMOVED***

	t, err := template.ParseFiles("templates/ScriptList.html")
	if err != nil ***REMOVED***
		log.Fatal(err)
	***REMOVED***

	//parse the login template and serve
	t.Execute(w, page)
***REMOVED***

func ExecuteTestScript(w http.ResponseWriter, r *http.Request) ***REMOVED***
	defer r.Body.Close()

	scriptJson := r.Form.Get("testScript")
	script := domain.TestScript***REMOVED******REMOVED***
	err := json.Unmarshal([]byte(scriptJson), &script)
	if err != nil ***REMOVED***
		log.Error("Invalid script form passed to executor")
		log.Error(err)
		return
	***REMOVED***
	log.Info(script)

***REMOVED***
