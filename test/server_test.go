package test

import (
	"log"
	"net/http"
	"testing"
	"time"
)
var gochan = make(chan bool)

func StartTestServer()***REMOVED***


***REMOVED***

func TestIndex(t *testing.T)***REMOVED***


	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8085/", nil)
	page, _ := http.DefaultClient.Do(req)
	log.Print(page.StatusCode)
	time.Sleep(200*time.Second)
***REMOVED***
