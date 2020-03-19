package test

import (
	"../service"
	"log"
	"net/http"
	"testing"
	"time"
)

var gochan = make(chan bool)

func StartTestServer() ***REMOVED***

***REMOVED***

func TestIndex(t *testing.T) ***REMOVED***

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8085/", nil)
	page, _ := http.DefaultClient.Do(req)
	log.Print(page.StatusCode)
	time.Sleep(200 * time.Second)
***REMOVED***

func TestCheckForFile(t *testing.T) ***REMOVED***
	if service.CheckForFile("test/server_test.go") == false ***REMOVED***
		log.Println("Utility function TestCheckForFile has failed its unit test.")
		t.FailNow()
	***REMOVED***
***REMOVED***
