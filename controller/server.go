package controller

import (
	"../domain"
	service "../service"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

//This is the server controller. It contains all server handlers.
const checkForExsistingHardwareTokenURL = "http://ec2-54-82-98-123.compute-1.amazonaws.com/CheckForExistingHardwareToken"
const getScriptsByHardwareTokenURL = "http://ec2-54-82-98-123.compute-1.amazonaws.com/GetScriptsByHardwareToken"
const getUserInfo = "http://ec2-54-82-98-123.compute-1.amazonaws.com/GetUserInfoByUsername"

var currentWebPage = domain.WebPage{
	Title:             "Diaverse",
	Scripts:           nil,
	Loggedin:          false,
	ScriptMap:         make(map[int]domain.TestScript),
	LoggedInUser:      "",
	LoggedInHWToken:   "",
	LoggedInEmail:     "",
	LoggedInBio:       "",
	LoggedInFirstName: "",
	LoggedInLastName:  "",
	Content:           "",
}
var scriptStopChan chan bool

func InitilizeScriptStopChan() {
	scriptStopChan = make(chan bool)
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

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	filedir := ""
	if strings.Contains(dir, "test") {
		//we are in test mode, the point of execution is different and therefore
		//so are the relative paths.

		filedir = "../templates/user.html"
	} else {
		filedir = "templates/user.html"
	}

	t, err := template.ParseFiles(filedir)
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

	} else if currentWebPage.Loggedin || (r.Method == http.MethodGet && r.FormValue("loginUsr") != "") {

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
			log.Debug(u)
			log.Debug(hw)
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
				log.Info(u + " Has Logged On.")
			}
		}

		form = url.Values{
			"username": {u},
		}

		//update user info
		resp, err = http.PostForm(getUserInfo, form)
		if err != nil {
			log.Fatal(err)
		}

		form = url.Values{
			"hardwareToken": {hw},
		}
		content, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		log.Warn(string(content))
		usrInfo := struct {
			Username  string `json:"Username"`
			Email     string `json:"Email"`
			FirstName string `json:"FirstName"`
			LastName  string `json:"LastName"`
			UserBio   string `json:"UserBio"`
		}{}

		err = json.Unmarshal(content, &usrInfo)
		currentWebPage.LoggedInEmail = usrInfo.Email
		currentWebPage.LoggedInUser = u
		currentWebPage.LoggedInFirstName = usrInfo.FirstName
		currentWebPage.LoggedInLastName = usrInfo.LastName
		currentWebPage.LoggedInBio = usrInfo.UserBio

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
			log.Error(e)
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
			scrp.ScriptID = v.ScriptID
			j, _ := json.Marshal(scrp)
			currentWebPage.ScriptsJSON = append(currentWebPage.ScriptsJSON, j)
			currentWebPage.ScriptMap[v.ScriptID] = scrp
			scripts[i] = scrp
			i++
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
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dir)
		filedir := ""
		if strings.Contains(dir, "test") {
			//we are in test mode, the point of execution is different and therefore
			//so are the relative paths.

			filedir = "../templates/login.html"
		} else {
			filedir = "templates/login.html"
		}

		t, err := template.ParseFiles(filedir)
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
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(dir)
		filedir := ""

		if strings.Contains(dir, "test") {
			//we are in test mode, the point of execution is different and therefore
			//so are the relative paths.

			filedir = "../templates/login.html"
		} else {
			filedir = "templates/login.html"
		}

		t, err := template.ParseFiles(filedir)
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

func ServeLogoutPage(w http.ResponseWriter, r *http.Request) {
	currentWebPage.Loggedin = false
	currentWebPage.LoggedInUser = ""
	currentWebPage.LoggedInHWToken = ""
	ServeWebpage(w, r)
}

var scriptInProgress struct {
	sync.RWMutex
	bool
}

func InitilizeStructs() {
	scriptInProgress.bool = false
}
func StopTestScriptHandler(w http.ResponseWriter, r *http.Request) {
	log.Warn("Sending signal")
	scriptStopChan <- true
	log.Warn("sent Signal down channel")
	scriptInProgress.Lock()
	scriptInProgress.bool = false
	scriptInProgress.Unlock()
	log.Warn("Stopped Test")
	ServeWebpage(w, r)
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

		scriptId, e := strconv.Atoi(r.URL.Query().Get("script"))
		if e != nil {
			log.Fatal(e)
		}
		r.ParseForm()

		script := currentWebPage.Scripts[scriptId]

		if r.FormValue("script") == "" {

			w.Write([]byte("Invalid Script format."))

			scriptInProgress.Lock()
			scriptInProgress.bool = false
			scriptInProgress.Unlock()
			return
		}

		scriptError := ExecuteTestScript(script)
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

//ExecuteTestScript is the function which executes the core logic of the project, it utilizes the code within service directory of the project to do the required audio I/O.
type TestCase struct {
	HardwareOutput []string `json:"hardwareOutput"`
	HardwareInput  []string `json:"hardwareInput"`
	VUIResponses   []string `json:"omitempty"`
	Result         float64  `json:"-"`
	TotalPassed    int      `json:"totalPass, omitempty"`
	TotalFailed    int      `json:"totalFail, omitempty"`
}

type TestScript struct {
	TestCases   []TestCase `json:"testCases"`
	PassPercent float64    `json:"passPercent, omitempty"`
}

func ExecuteTestScript(script domain.TestScript) error {
	for _, e := range script.TestCases {
		for i := 0; i < len(e.HardwareInput); i++ {
			service.PrepareAudioFiles(e.HardwareOutput)
			select {
			case <-scriptStopChan:
				log.Warn("Stopping Test")
				return nil
			default:
				service.SpeakAloud(e.HardwareOutput[i])
				response, _, err := service.Recognize()
				if err != nil {
					return err
				}
				log.Warn("I HEARD " + response)
				e.HardwareInput[i] = response
			}
		}
	}

	log.Debug(script.TestCases)
	j, _ := json.Marshal(script)
	fmt.Println(string(j))
	return nil
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
