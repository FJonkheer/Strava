package Handler

// MAtrikel-Nr 3736476, 8721083
import (
	"fmt"
	"github.com/FJonkheer/Strava/src/Helper"
	"net/http"
	"strings"
)

func HandleReview(w http.ResponseWriter, r *http.Request) {
	name, _ := r.Cookie("Name")
	password, _ := r.Cookie("Password")

	if !validateUser(name.Value, password.Value) {
		http.Redirect(w, r, "/Login", 301)
	} else {
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
}
func DownloadHandler(w http.ResponseWriter, r *http.Request) { //Download einer Datei
name, _ := r.Cookie("Name")
	password, _ := r.Cookie("Password")

	if !validateUser(name.Value, password.Value) {
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
	name, _ := r.Cookie("Name")
	password, _ := r.Cookie("Password")

	if !validateUser(name.Value, password.Value) {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + name.Value + "/" //Benutzername muss abgefragt werden
		file := Helper.GetfileName(r)       //Das Feld, wo die Datei ausgewählt wurde
		path = path + file
		err := Helper.DeleteFiles(path)
		if err != nil {
			fmt.Println("Konnte Datei nicht löschen")
		}
		http.Redirect(w, r, "/Review", 301)
	}
}

func changeHandler(w http.ResponseWriter, r *http.Request) { //Ändern der InfoPage
	name, _ := r.Cookie("Name")
	password, _ := r.Cookie("Password")

	if !validateUser(name.Value, password.Value) {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + name.Value + "/" //Benutzername muss abgefragt werden
		file := Helper.GetfileName(r)
		path = path + file
		Helper.ChangeInfoFile(w, r, path)
		http.Redirect(w, r, "/Review", 301)
	}
}

func searchcomment(w http.ResponseWriter, r *http.Request) {
	name, _ := r.Cookie("Name")
	password, _ := r.Cookie("Password")

	if !validateUser(name.Value, password.Value) {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + name.Value + "/" //Benutzername muss abgefragt werden
		comment := r.FormValue("searchcomment")
		var files []string
		csvfiles := Helper.Scanforcsvfiles(path)
		for _, file := range csvfiles {
			content, _ := Helper.ReadCsv(path + file)
			if strings.Contains(content[1][2], comment) {
				files = append(files, file)
			}
		}
		for _, file := range files {
			file = strings.Replace(file, ".csv", "", -1)
			_, err := fmt.Fprintf(w, "<div>%s</div><br>", file)
			if err != nil {
				fmt.Println("Fehler bei Ausgabe")
			}
		}
	}
}
