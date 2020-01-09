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
	- Die Review-Page-Funktionalitäten
		- Anzeige der GPX-Dateien
		- Anzeige der gespeicherten Informationen zu jeder GPX-Datei
		- Download von GPX/ZIP-Dateien
		- Löschen von GPX-Dateien (mit Bestätigung)
		- Die Verarbeitung der .gpx-Dateien
		- Art der Aktivität und Kommentar müssen bearbeitet werden können
	- Ordner pro hochladenem Nutzer erstellen und Dateien dort abspeichern
	- Abfrage ob Datei bereits vorhanden ist beim Upload (überschreibt das vorhandene File, will man das?)
	- Bei Fehler im Einloggen - zurück zum Einloggen mit Fehlermeldung (keine weiße Seite nur wmit Fehlermeldung)
	- Upload-Konventionen (Muss ein Kommentar eingegeben werden? Wurde überhaupt eine Datei ausgewählt?)
	- Logging
	- Tests

Error-Handling
Kommentare
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
	fmt.Fprintf(w, "<div>%s</div>", p.Body) //gibt den Body dem Benutzer aus
}

func back(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/MainPage", 301) //Zurück zur Startseite
}

func main() {
	http.HandleFunc("/", handler)                             //der generelle Seitenaufruf
	http.HandleFunc("/submit.php", Handler.Handling)          //der Login-Aufruf
	http.HandleFunc("/redirect.php", Handler.Redirecting)     //die Weiterleitung von der Startseite
	http.HandleFunc("/back.php", back)                        //das Zurückleiten auf die MainPage
	http.HandleFunc("/upload.php", Handler.Uploader)          //der Dateiupload
	http.HandleFunc("/download.php", Handler.DownloadHandler) //Download einer Datei
	http.HandleFunc("/delete.php", Handler.DeleteHandler)     //Löschen von Dateien
	http.HandleFunc("/change.php", Handler.ChangeHandler)
	log.Fatal(http.ListenAndServeTLS(":9090", "src/Auth/cert.pem", "src/Auth/key.pem", nil)) //der "Webserver"
}
