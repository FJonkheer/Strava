package Handler

import (
	"Helper"
	"net/http"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) { //Download einer Datei
	path := "Files/Username"    //Benutzername muss abgefragt werden
	file := r.FormValue("File") //Das Feld, wo die Datei ausgewählt wurde
	path = path + file
	Helper.DownloadFile(w, r, path)
	http.Redirect(w, r, "/Review", 301)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) { //Löschen eines Eintrags
	//TODO: Abfrage ob wirklich gelöscht werden soll
	path := "Files/Username"    //Benutzername muss abgefragt werden
	file := r.FormValue("File") //Das Feld, wo die Datei ausgewählt wurde
	path = path + file
	Helper.DeleteFiles(path)
	http.Redirect(w, r, "/Review", 301)
}

func ChangeHandler(w http.ResponseWriter, r *http.Request) { //Ändern der InfoPage
	path := "Files/Username"    //Benutzername muss abgefragt werden
	file := r.FormValue("File") //Das Feld, wo die Datei ausgewählt wurde
	path = path + file
	Helper.ChangeInfoFile(w, r, path)
	http.Redirect(w, r, "/Review", 301)
}
