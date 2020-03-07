package controller

import (
	"../domain"
	service "../service"
	"encoding/json"
	"github.com/prometheus/common/log"
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

	scriptJSON, err := ioutil.ReadAll(r.Body)
	if err != nil ***REMOVED***
		log.Error("Cannot read json body of requested test script")
		w.Write([]byte("Cannot read json body of requested test script"))
		r.Body.Close()
		return
	***REMOVED***

	testScript := domain.TestScript***REMOVED******REMOVED***
	err = json.Unmarshal(scriptJSON, &testScript)
	if err != nil ***REMOVED***
		log.Error("Could not unmarshal test script json")
		r.Body.Close()
		w.Write([]byte("Could not unmarshal test script json"))
	***REMOVED***

	//caseResults := make([]bool, len(testScript.Cases))
	for _, e := range testScript.Cases ***REMOVED***

		//generate all required TTS
		for j, k := range e.Responses ***REMOVED***

			//the test hardware always speaks first.
			service.SpeakAloud(k)

			//try and listen three times before we continue
			for n := 0; n < 2; n++ ***REMOVED***
				resp, conf, err := service.Recognize()
				if err != nil ***REMOVED***
					log.Fatal("Failed to begin recognition")
				***REMOVED***

				if conf >= 0.8 ***REMOVED***
					if service.TranscriptionConfidence(resp, e.ExpectedOutput[j]) >= 0.8 ***REMOVED***
						//only continue until we have both an ASR confidence and NLU confidence of 80% or more.
						break
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					service.SpeakAloud(k)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
