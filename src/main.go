package main

import (
	"Handler"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/*
TODO:
- Die Review-Page
- Erweiterte Funktionalitäten der Upload-Seite (Kommentar, Laufen/Fahrrad und Uploaddatum abspeichern
- Ordner pro hochladenem Nutzer erstellen und Dateien dort abspeichern
- Abfrage ob Datei bereits vorhanden ist
- Was passiert mit den "geuploadeten" Dateien - regeln wo die hinmüssen, wie die abgespeichert werden
- Bei Fehler im Einloggen - zurück zum Einloggen mit Fehlermeldung
*/

type Page struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := ioutil.ReadFile("Sites/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	p, err := loadPage(title)
	if err != nil {
		title = "Login"
		p, _ = loadPage(title)
	}
	fmt.Fprintf(w, "<div>%s</div>", p.Body)
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/MainPage", 301)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/submit.php", Handler.Handling)
	http.HandleFunc("/redirect.php", Handler.Redirecting)
	http.HandleFunc("/back.php", back)
	http.HandleFunc("/upload.php", Handler.Uploader)
	log.Fatal(http.ListenAndServeTLS(":9090", "src/Auth/cert.pem", "src/Auth/key.pem", nil))
}
