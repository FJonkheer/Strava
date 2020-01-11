package main

import (
	"Handler"
	"Helper"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

/*
TODO:
	- Die Review-Page
		- Sortieren der Einträge nach Datum
		- alle Dateien gleichzeitig auswählbar (nicht gut)
	- Caching
	- Tests
	- Matrikelnummer zu allen src-Dateien hinzufügen (Meine ist: 3736476 ^^)

Error-Handling
Kommentare

Optional:
	- Abfrage ob Datei bereits vorhanden ist beim Upload (überschreibt das vorhandene File, will man das?)
	- Bei Fehler im Einloggen - zurück zum Einloggen mit Fehlermeldung (keine weiße Seite nur mit Fehlermeldung)
	- Upload-Konventionen (Muss ein Kommentar eingegeben werden? Wurde überhaupt eine Datei ausgewählt?)
	- Volltextsuche für Kommentare (Bessere Ausgabe?)
*/

type Page struct { //Die Struktur einer Website
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := ioutil.ReadFile("Sites/" + filename) //liest den body der aufzurufenden Seite aus
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil //Gibt die Seite mit Titel und Body zurück
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):] //nimmt den Titel aus dem Request
	p, err := loadPage(title)      //lädt den Body mit dem Titel
	if err != nil {                //wenn die gesuchte Seite nicht existiert, soll man auf die Login-Seite kommen
		title = "Login"
		p, _ = loadPage(title)
	}
	_, err = fmt.Fprintf(w, "<div>%s</div>", p.Body) //gibt den Body dem Benutzer aus
	if err != nil {
		fmt.Println(err)
	}
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/MainPage", 301) //Zurück zur Startseite
}

var paths = []string{"Sites/Review.html"}
var t = template.Must(template.New("Review.html").ParseFiles(paths...))

func renderReview(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(Handler.Uname)
	user := Helper.Parsecsvtostruct(cookie.Value)
	err := t.Execute(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	portFlag := flag.String("port", ":9090", "choose Port")
	flag.Parse()
	http.HandleFunc("/", handler)                         //der generelle Seitenaufruf
	http.HandleFunc("/submit.php", Handler.Handling)      //der Login-Aufruf
	http.HandleFunc("/redirect.php", Handler.Redirecting) //die Weiterleitung von der Startseite
	http.HandleFunc("/back.php", back)                    //das Zurückleiten auf die MainPage
	http.HandleFunc("/upload.php", Handler.Uploader)      //der Dateiupload
	http.HandleFunc("/review.php", Handler.HandleReview)
	http.HandleFunc("/Review", renderReview)
	log.Fatal(http.ListenAndServeTLS(*portFlag, "src/Auth/cert.pem", "src/Auth/key.pem", nil)) //der "Webserver"
}
