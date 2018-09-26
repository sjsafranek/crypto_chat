package main

func CreateUser(email, password string) {
	user := User{Username: email, Email: email}
	user.SetPassword(password)
	db.Add(&user)
	db.Save(users_file)
}
