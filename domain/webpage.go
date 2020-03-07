package domain

//this file contains the data structure used in the serving of html to the user

type LoginPage struct {
	Title     string
	AuthToken string
	Content   string //The actual html
}

type ListPage struct {
	Title          string
	ScriptList     []TestScript
	SelectedScript string
}
