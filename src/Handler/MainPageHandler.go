package Handler

import (
	"net/http"
)

func Redirecting(w http.ResponseWriter, r *http.Request) { //Weiterleiten auf Upload oder Review
	name, _ := r.Cookie("Name")
	password, _ := r.Cookie("Password")

	if !validateUser(name.Value, password.Value) {
		http.Redirect(w, r, "/Login", 301)
	} else {
		if r.FormValue("redirect") == "Upload File" {
			upload(w, r)
		} else {
			review(w, r)
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request) {

	name, _ := r.Cookie("Name")
	password, _ := r.Cookie("Password")

	if !validateUser(name.Value, password.Value) {
		http.Redirect(w, r, "/Login", 301)
	} else {
		http.Redirect(w, r, "/Upload", 301) //Upload-Seite
	}
}

func review(w http.ResponseWriter, r *http.Request) {
	name, _ := r.Cookie("Name")
	password, _ := r.Cookie("Password")

	if !validateUser(name.Value, password.Value) {
		http.Redirect(w, r, "/Login", 301)
	} else {
		http.Redirect(w, r, "/Review", 301) //Review-Seite
	}
}
