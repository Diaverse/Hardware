package domain

//this file contains the data structure used in the serving of html to the user

type WebPage struct {
	Title     string
	AuthToken string
	Scripts   []TestScript
	Loggedin  bool
	Content   string //The actual html
}
