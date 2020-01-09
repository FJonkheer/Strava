package Handler

import (
	"net/http"
)

func Redirecting(w http.ResponseWriter, r *http.Request) { //Weiterleiten auf Upload oder Review
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
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
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {
		http.Redirect(w, r, "/Upload", 301) //Upload-Seite
	}
}

func review(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {

		http.Redirect(w, r, "/Review", 301) //Review-Seite
	}
}
