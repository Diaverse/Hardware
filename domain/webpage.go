package domain

//this file contains the data structure used in the serving of html to the user

type LoginPage struct ***REMOVED***
	Title     string
	AuthToken string
	Content   string //The actual html
***REMOVED***

type ListPage struct ***REMOVED***
	Title          string
	ScriptList     []TestScript
	SelectedScript string
***REMOVED***
