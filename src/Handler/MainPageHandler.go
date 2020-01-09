package Handler

import "net/http"

func Redirecting(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("redirect") == "Upload File" {
		upload(w, r)
	} else {
		review(w, r)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/Upload", 301)
}

func review(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/Review", 301)
}
