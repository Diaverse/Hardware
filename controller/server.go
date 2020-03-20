package controller

import (
	"../domain"
	"encoding/json"
	"github.com/prometheus/common/log"
	"sync"

	"html/template"
	"io/ioutil"
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
	//page := domain.ListPage***REMOVED***
	//	Title: "Diaverse Script View",
	//	ScriptList: []domain.TestScript***REMOVED***
	//		domain.TestScript***REMOVED***
	//			Cases: []domain.TestCase***REMOVED***
	//				domain.TestCase***REMOVED***
	//					Responses:      []string***REMOVED***"Hello, how are you"***REMOVED***,
	//					ExpectedOutput: []string***REMOVED***"I am fine"***REMOVED***,
	//				***REMOVED***,
	//				domain.TestCase***REMOVED***
	//					Responses:      []string***REMOVED***"what are you doing?"***REMOVED***,
	//					ExpectedOutput: []string***REMOVED***"absolutely nothing."***REMOVED***,
	//				***REMOVED***,
	//			***REMOVED***,
	//			Result: true,
	//		***REMOVED***,
	//		domain.TestScript***REMOVED***
	//			Cases:  []domain.TestCase***REMOVED******REMOVED***,
	//			Result: true,
	//		***REMOVED***,
	//	***REMOVED***,
	//	SelectedScript: "Script One",
	//***REMOVED***

	//t, err := template.ParseFiles("templates/ScriptList.html")
	//if err != nil ***REMOVED***
	//	log.Fatal(err)
	//***REMOVED***
	//
	////parse the login template and serve
	//t.Execute(w, page)
***REMOVED***

var scriptInProgress struct ***REMOVED***
	sync.RWMutex
	bool
***REMOVED***

func InitilizeStructs() ***REMOVED***
	scriptInProgress.bool = false
***REMOVED***

func ExecuteTestScriptHandler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	scriptInProgress.Lock()
	curState := scriptInProgress.bool
	scriptInProgress.Unlock()

	if !curState ***REMOVED***
		scriptInProgress.Lock()
		scriptInProgress.bool = true
		scriptInProgress.Unlock()

		defer r.Body.Close()
		log.Info("Got test script")

		scriptJSON, err := ioutil.ReadAll(r.Body)
		if err != nil ***REMOVED***
			log.Error("Cannot read json body of requested test script")
			w.Write([]byte("Cannot read json body of requested test script"))
			r.Body.Close()
			scriptInProgress.Lock()
			scriptInProgress.bool = false
			scriptInProgress.Unlock()
			return
		***REMOVED***

		script := domain.TestScript***REMOVED******REMOVED***
		e := json.Unmarshal(scriptJSON, &script)
		if e != nil ***REMOVED***

			w.Write([]byte("Invalid Script format."))
			scriptInProgress.Lock()
			scriptInProgress.bool = false
			scriptInProgress.Unlock()
			return
		***REMOVED***

		scriptError := ExecuteTestScript(&script)
		if scriptError != nil ***REMOVED***

			w.Write([]byte("Error executing test script, view service logs for additional information."))
			scriptInProgress.Lock()
			scriptInProgress.bool = false
			scriptInProgress.Unlock()
			return
		***REMOVED*** else ***REMOVED***
			if script.Result == true ***REMOVED***
				w.Write([]byte("Test script completed successfully!"))
				scriptInProgress.Lock()
				scriptInProgress.bool = false
				scriptInProgress.Unlock()
				return
			***REMOVED*** else ***REMOVED***
				w.Write([]byte("Test script failed. Check service logs for more information."))
				scriptInProgress.Lock()
				scriptInProgress.bool = false
				scriptInProgress.Unlock()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		w.Write([]byte("Test script already in progress, resend script after current script completes, or cancel current script."))
	***REMOVED***
***REMOVED***

func RegisterHardware(w http.ResponseWriter, r *http.Request) ***REMOVED***

***REMOVED***
