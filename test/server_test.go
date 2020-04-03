package test

import (
	controller "../controller"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.FailNow()
	}

	requestRecorder := httptest.NewRecorder()
	webpageHandler := http.HandlerFunc(controller.ServeWebpage)
	webpageHandler.ServeHTTP(requestRecorder, req)

	status := requestRecorder.Code
	if status != http.StatusOK {
		t.FailNow()
	}

	resp := requestRecorder.Result()
	c, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	t.Log(string(c))
}

func TestDashboard(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.FailNow()
	}
	req.ParseForm()
	req.Form.Set("loginUsr", "harrison")
	req.Form.Set("loginPass", "token1")

	reqRecorder := httptest.NewRecorder()
	webpageHandler := http.HandlerFunc(controller.ServeWebpage)
	webpageHandler.ServeHTTP(reqRecorder, req)
	status := reqRecorder.Code

	if status != http.StatusOK {
		t.FailNow()
	} else {
		resp := reqRecorder.Result()
		c, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		t.Log(string(c))
	}
}

func TestUsersPage(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)
	if err != nil {
		t.FailNow()
	}
	req.ParseForm()
	req.Form.Set("loginUsr", "harrison")
	req.Form.Set("loginPass", "token1")

	reqRecorder := httptest.NewRecorder()
	webpageHandler := http.HandlerFunc(controller.ServeWebpage)
	webpageHandler.ServeHTTP(reqRecorder, req)
	status := reqRecorder.Code

	if status != http.StatusOK {
		t.FailNow()
	} else {
		usersPageHandler := http.HandlerFunc(controller.ServeUsersWebPage)
		usersPageHandler.ServeHTTP(reqRecorder, req)
		fmt.Println(reqRecorder.Code)
		resp := reqRecorder.Result()
		c, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		t.Log(string(c))
	}
}

func TestCheckForFile(t *testing.T) {
	//if service.CheckForFile("test/server_test.go") == false {
	//	log.Println("Utility function TestCheckForFile has failed its unit test.")
	//	t.FailNow()
	//}
}
