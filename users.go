package main

import (
	"github.com/sjsafranek/simpleusers"
)

var users_file = "users.json"
var db simpleusers.Users

func CreateUser(email, password string) {
	user := simpleusers.User{Username: email, Email: email}
	user.SetPassword(password)
	db.Add(&user)
	db.Save(users_file)
}

func init() {
	db = simpleusers.Users{}
	db.Fetch(users_file)
	user := simpleusers.User{Username: "admin", Email: "admin@email.com"}
	user.SetPassword("dev")
	db.Add(&user)
	db.Save(users_file)
}
