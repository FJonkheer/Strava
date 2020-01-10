package Handler

import (
	"Helper"
	"fmt"
	"net/http"
	"strings"
)

func HandleReview(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("review") {
	case "Delete Record":
		deleteHandler(w, r)
		break
	case "Search":
		searchcomment(w, r)
		break
	case "Change Info":
		changeHandler(w, r)
		break
	default:
		DownloadHandler(w, r)
		break
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) { //Download einer Datei

	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + cookie.Value //Benutzername muss abgefragt werden
		file := Helper.GetfileName(r)   //Das Feld, wo die Datei ausgewählt wurde
		path = path + "/" + file
		Helper.DownloadFile(w, r, path)
		http.Redirect(w, r, "/Review", 301)
	}
}

//Löschen eines Eintrags
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: Abfrage ob wirklich gelöscht werden soll
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + cookie.Value + "/" //Benutzername muss abgefragt werden
		file := Helper.GetfileName(r)         //Das Feld, wo die Datei ausgewählt wurde
		path = path + file
		Helper.DeleteFiles(path)
		http.Redirect(w, r, "/Review", 301)
	}
}

func changeHandler(w http.ResponseWriter, r *http.Request) { //Ändern der InfoPage
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + cookie.Value + "/" //Benutzername muss abgefragt werden
		file := Helper.GetfileName(r)
		path = path + file
		Helper.ChangeInfoFile(w, r, path)
		http.Redirect(w, r, "/Review", 301)
	}
}

func searchcomment(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + cookie.Value + "/" //Benutzername muss abgefragt werden
		comment := r.FormValue("searchcomment")
		var files []string
		csvfiles := Helper.Scanforcsvfiles(path)
		for _, file := range csvfiles {
			content, _ := Helper.ReadCsv(file)
			if strings.Contains(content[1][2], comment) {
				files = append(files, file)
			}
		}
		fmt.Fprintf(w, "<div>%s</div>", files)
	}
}
