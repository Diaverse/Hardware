package controller

import (
	"../domain"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

//This is the server controller. It contains all server handlers.
const checkForExsistingHardwareTokenURL = "http://ec2-54-82-98-123.compute-1.amazonaws.com/CheckForExistingHardwareToken"
const getScriptsByHardwareTokenURL = "http://ec2-54-82-98-123.compute-1.amazonaws.com/GetScriptsByHardwareToken"
const getUserInfo = "http://ec2-54-82-98-123.compute-1.amazonaws.com/GetUserInfo"

var currentWebPage = domain.WebPage{
	Title:           "Diaverse",
	Scripts:         nil,
	Loggedin:        false,
	LoggedInUser:    "",
	LoggedInHWToken: "",
	Content:         "",
}

func ServeUsersWebPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	req, err := http.NewRequest(http.MethodGet, getUserInfo, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.ParseForm()
	req.Form.Set("hwtoken", currentWebPage.LoggedInHWToken)
	req.Form.Set("user", currentWebPage.LoggedInUser)

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles("templates/user.html")
	if err != nil {
		log.Fatal(err)
	}

	//parse the login template and serve
	e := t.Execute(w, currentWebPage)
	if e != nil {
		log.Fatal(e)
	}
}

//ServeWebpage listens on the root directory of the hosted web page and serves the needed html to the user
func ServeWebpage(w http.ResponseWriter, r *http.Request) {
	//todo; define a way for a user to remain logged in across sessions
	//create a new login page struct
	r.ParseForm()
	if strings.Contains(r.URL.String(), "user") {
		ServeUsersWebPage(w, r)
		return

	} else if (r.Method == http.MethodGet && r.FormValue("loginUsr") != "") || currentWebPage.Loggedin {
		u := ""
		hw := ""
		if !currentWebPage.Loggedin {
			u = r.FormValue("loginUsr")
			hw = r.FormValue("loginPass")
			currentWebPage.LoggedInUser = u
			currentWebPage.LoggedInHWToken = hw
		} else {
			u = currentWebPage.LoggedInUser
			hw = currentWebPage.LoggedInHWToken
		}

		form := url.Values{
			"username": {u},
			"hwtoken":  {hw},
		}

		resp, err := http.PostForm(checkForExsistingHardwareTokenURL, form)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != 202 {
			log.Info("Detected Invalid Login Attempt via Hardware UI ")
		} else {
			if r.FormValue("loginUsr") != "" {
				log.Info(r.FormValue("loginUsr") + " Has Logged On.")
			}
		}

		form = url.Values{
			"hardwareToken": {hw},
		}
		//get scripts
		resp, err = http.PostForm(getScriptsByHardwareTokenURL, form)
		if err != nil {
			log.Fatal(err)
		}

		type row struct {
			HardwareToken string `json:"Hardwaretoken"`
			Script        string `json:"Script"`
			ScriptID      int    `json:"ScriptID"`
		}

		s := make(map[string]row)

		c, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			log.Info(string(c))
			log.Fatal(e)
		}

		e = json.Unmarshal(c, &s)
		if e != nil {
			log.Fatal(e)
		}

		scriptsContents := make(map[int]string)
		i := 0
		for range s {
			i++
		}
		scripts := make([]domain.TestScript, i)
		i = 0
		for _, v := range s {
			scriptsContents[v.ScriptID] = v.Script
			scrp := domain.TestScript{}
			e = json.Unmarshal([]byte(v.Script), &scrp)
			scripts[i] = scrp
		}

		//sample data
		//sampleScripts := []domain.TestScript{
		//	domain.TestScript{
		//		TestCases: []domain.TestCase{
		//			domain.TestCase{
		//				HardwareOutput: []string{"Test Case one"},
		//				HardwareInput:  []string{"two"},
		//				Result:         0,
		//				TotalPassed:    0,
		//				TotalFailed:    0,
		//			},
		//			domain.TestCase{
		//				HardwareOutput: []string{"Test Case Two"},
		//				HardwareInput:  []string{"two"},
		//				Result:         0,
		//				TotalPassed:    0,
		//				TotalFailed:    0,
		//			},
		//			domain.TestCase{
		//				HardwareOutput: []string{"Test Case three"},
		//				HardwareInput:  []string{"two"},
		//				Result:         0,
		//				TotalPassed:    0,
		//				TotalFailed:    0,
		//			},
		//		},
		//		PassPercent: 0,
		//	},
		//	domain.TestScript{
		//		TestCases: []domain.TestCase{
		//			domain.TestCase{
		//				HardwareOutput: []string{"Test Case one"},
		//				HardwareInput:  []string{"two"},
		//				Result:         0,
		//				TotalPassed:    0,
		//				TotalFailed:    0,
		//			},
		//			domain.TestCase{
		//				HardwareOutput: []string{"Test Case Two"},
		//				HardwareInput:  []string{"two"},
		//				Result:         0,
		//				TotalPassed:    0,
		//				TotalFailed:    0,
		//			},
		//			domain.TestCase{
		//				HardwareOutput: []string{"Test Case three"},
		//				HardwareInput:  []string{"two"},
		//				Result:         0,
		//				TotalPassed:    0,
		//				TotalFailed:    0,
		//			},
		//		},
		//		PassPercent: 0,
		//	},
		//}

		//process template
		currentWebPage.Scripts = scripts

		t, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Fatal(err)
		}
		currentWebPage.Loggedin = true
		//parse the login template and serve
		e = t.Execute(w, currentWebPage)
		if e != nil {
			log.Fatal(e)
		}

	} else {

		t, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Fatal(err)
		}

		//parse the login template and serve
		e := t.Execute(w, currentWebPage)
		if e != nil {
			log.Fatal(e)
		}

	}
}

func ServeUsersPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

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
		r.ParseForm()

		script := domain.TestScript{}
		e := json.Unmarshal(scriptJSON, &script)
		if e != nil && r.FormValue("script") == "" {

			w.Write([]byte("Invalid Script format."))

			scriptInProgress.Lock()
			scriptInProgress.bool = false
			scriptInProgress.Unlock()
			return
		} else {
			scriptJSON = []byte(r.Form.Get("script"))
			fmt.Printf(string(scriptJSON))
			e := json.Unmarshal(scriptJSON, &script)
			if e != nil {
				log.Info("Invalid script format.")
			}
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
			fmt.Println(string(j))
			w.Write(j)
			scriptInProgress.Lock()
			scriptInProgress.bool = false
			scriptInProgress.Unlock()
		}
	} else {
		w.Write([]byte("Test script already in progress, resend script after current script completes, or cancel current script."))
	}
}

func RegisterHardware(w http.ResponseWriter, r *http.Request) {

}

func CheckForExsistingHardwareToken(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	token := r.FormValue("hardwareToken")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	req, err := http.NewRequest(http.MethodGet, checkForExsistingHardwareTokenURL, strings.NewReader(token))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	log.Info(resp.StatusCode)

	type ServerResp struct {
		user          string
		hardwareToken string
	}
}

func AuthorizeToken(token string) {

}
