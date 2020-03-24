package controller

import (
	"../domain"
	"encoding/json"
	"github.com/prometheus/common/log"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
)

//This is the server controller. It contains all server handlers.

//ServeLogin listens on the root directory of the hosted web page and serves the needed html to the user
func ServeLogin(w http.ResponseWriter, r *http.Request) {
	//todo; define a way for a user to remain logged in across sessions
	//create a new login page struct
	page := domain.LoginPage{
		Title:     "Diaverse Login Screen",
		AuthToken: "",
		Content:   "<html><p>Hello</p></html>",
	}
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	//parse the login template and serve
	t.Execute(w, page)
}

func ServeScriptListView(w http.ResponseWriter, r *http.Request) {
	//page := domain.ListPage{
	//	Title: "Diaverse Script View",
	//	ScriptList: []domain.TestScript{
	//		domain.TestScript{
	//			Cases: []domain.TestCase{
	//				domain.TestCase{
	//					Responses:      []string{"Hello, how are you"},
	//					ExpectedOutput: []string{"I am fine"},
	//				},
	//				domain.TestCase{
	//					Responses:      []string{"what are you doing?"},
	//					ExpectedOutput: []string{"absolutely nothing."},
	//				},
	//			},
	//			Result: true,
	//		},
	//		domain.TestScript{
	//			Cases:  []domain.TestCase{},
	//			Result: true,
	//		},
	//	},
	//	SelectedScript: "Script One",
	//}

	//t, err := template.ParseFiles("templates/ScriptList.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	////parse the login template and serve
	//t.Execute(w, page)
}

var scriptInProgress struct {
	sync.RWMutex
	bool
}

func InitilizeStructs() {
	scriptInProgress.bool = false
}

func ExecuteTestScriptHandler(w http.ResponseWriter, r *http.Request) {
	scriptInProgress.Lock()
	curState := scriptInProgress.bool
	scriptInProgress.Unlock()

	if !curState {
		scriptInProgress.Lock()
		scriptInProgress.bool = true
		scriptInProgress.Unlock()

		defer r.Body.Close()
		log.Info("Got test script")

		scriptJSON, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Cannot read json body of requested test script")
			w.Write([]byte("Cannot read json body of requested test script"))
			r.Body.Close()
			scriptInProgress.Lock()
			scriptInProgress.bool = false
			scriptInProgress.Unlock()
			return
		}

		script := domain.TestScript{}
		e := json.Unmarshal(scriptJSON, &script)
		if e != nil {

			w.Write([]byte("Invalid Script format."))
			scriptInProgress.Lock()
			scriptInProgress.bool = false
			scriptInProgress.Unlock()
			return
		}

		scriptError := ExecuteTestScript(&script)
		if scriptError != nil {
			w.Write([]byte("Error executing test script, view service logs for additional information."))
			scriptInProgress.Lock()
			scriptInProgress.bool = false
			scriptInProgress.Unlock()
			return
		} else {
			j, e := json.Marshal(script)
			if e != nil {
				log.Fatal("cannot create result JSON for current test, fatal")
			}

			w.Write(j)
		}
	} else {
		w.Write([]byte("Test script already in progress, resend script after current script completes, or cancel current script."))
	}
}

func RegisterHardware(w http.ResponseWriter, r *http.Request) {

}
