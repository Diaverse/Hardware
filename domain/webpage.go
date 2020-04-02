package domain

import "encoding/json"

//this file contains the data structure used in the serving of html to the user

type WebPage struct {
	Title           string
	Scripts         []TestScript
	ScriptsJSON     []json.RawMessage
	Loggedin        bool
	LoggedInUser    string
	LoggedInHWToken string
	Content         string //The actual html
}
