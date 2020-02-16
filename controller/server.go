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

	scriptJson := r.Form.Get("testScript")
	script := domain.TestScript{}
	err := json.Unmarshal([]byte(scriptJson), &script)
	if err != nil {
		log.Error("Invalid script form passed to executor")
		log.Error(err)
		return
	}
	log.Info(script)

}
