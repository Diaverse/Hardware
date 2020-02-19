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
func ServeLogin(w http.ResponseWriter, r *http.Request) {
	//todo; define a way for a user to remain logged in across sessions
	//create a new login page struct
	page := domain.LoginPage{
		Title:     "Diaverse Login Screen",
		AuthToken: "",
		Content:   "<html><p>Hello</p></html>",
	}
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Fatal(err)
	}

	//parse the login template and serve
	t.Execute(w, page)
}

func ServeScriptListView(w http.ResponseWriter, r *http.Request) {
	page := domain.ListPage{
		Title: "Diaverse Script View",
		ScriptList: []domain.TestScript{
			domain.TestScript{
				Cases: []domain.TestCase{
					domain.TestCase{
						Responses:      []string{"Hello, how are you"},
						ExpectedOutput: []string{"I am fine"},
					},
					domain.TestCase{
						Responses:      []string{"what are you doing?"},
						ExpectedOutput: []string{"absolutely nothing."},
					},
				},
				Result: true,
			},
			domain.TestScript{
				Cases:  []domain.TestCase{},
				Result: true,
			},
		},
		SelectedScript: "Script One",
	}

	t, err := template.ParseFiles("templates/ScriptList.html")
	if err != nil {
		log.Fatal(err)
	}

	//parse the login template and serve
	t.Execute(w, page)
}

func ExecuteTestScript(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	scriptJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Cannot read json body of requested test script")
		w.Write([]byte("Cannot read json body of requested test script"))
		r.Body.Close()
		return
	}

	testScript := domain.TestScript{}
	err = json.Unmarshal(scriptJSON, &testScript)
	if err != nil {
		log.Error("Could not unmarshal test script json")
		r.Body.Close()
		w.Write([]byte("Could not unmarshal test script json"))
	}

	//caseResults := make([]bool, len(testScript.Cases))
	for _, e := range testScript.Cases {

		//generate all required TTS
		for j, k := range e.Responses {

			//the test hardware always speaks first.
			service.SpeakAloud(k)

			//try and listen three times before we continue
			for n := 0; n < 2; n++ {
				resp, conf, err := service.Recognize()
				if err != nil {
					log.Fatal("Failed to begin recognition")
				}

				if conf >= 0.8 {
					if service.TranscriptionConfidence(resp, e.ExpectedOutput[j]) >= 0.8 {
						//only continue until we have both an ASR confidence and NLU confidence of 80% or more.
						break
					}
				} else {
					service.SpeakAloud(k)
				}
			}
		}
	}
}
