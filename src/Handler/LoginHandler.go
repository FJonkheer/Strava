package Handler

import (
	"Helper"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var Uname string

type User struct { //Die Struktur eines Benutzers
	username string
	password string
}

func Handling(w http.ResponseWriter, r *http.Request) { //Wurde "Login" oder "Registrieren" gedrückt?
	if r.FormValue("login") == "Login" {
		Login(w, r)
	} else {
		register(w, r)
	}
}

func register(w http.ResponseWriter, r *http.Request) {

	Uname := r.FormValue("uname")
	password := r.FormValue("pword")
	salt := "15967" //Der Salt, welcher zum Verschlüsseln genutzt wird
	isGone := false
	if len(Uname) > 0 && len(password) > 0 { //Wurde ein Benutzername und ein Passwort eingegeben?
		password = password + salt           //Salting des Passworts
		pword := Helper.GetMD5Hash(password) //Hashing des Passworts mit Salt

		t, _ := Helper.FilePathExists("data/userdata") //gibt es den Dateispeicherpfad noch nicht?
		if !t {
			Helper.CreateFolders("data/userdata") //Dann werden erst noch die neuen Ordner erstellt
		}

		if !Helper.FileExists("data/userdata/Test.csv") { //existiert die Speicherdatei schon?
			empData := [][]string{
				{"uname", "pword"}}
			csvFile, err := os.Create("data/userdata/Test.csv") //Wird eine Speicherdatei erstellt
			if err != nil {
				log.Fatalf("failed creating file: %s", err)
			}
			csvwriter := csv.NewWriter(csvFile)
			for _, empRow := range empData {
				_ = csvwriter.Write(empRow) //Und die Indexe geschrieben
			}
			csvwriter.Flush()
			csvFile.Close()
		}
		lines, err := Helper.ReadCsv("data/userdata/Test.csv") //CSV-Auslesen
		if err != nil {
			panic(err)
		}
		for _, line := range lines { //Prüfen ob es den Benutzernamen schon in der Speicherdatei gibt
			data := User{
				username: line[0],
				password: line[1],
			}
			if data.username == Uname {
				isGone = true
			}
		}
		if isGone {
			fmt.Fprintf(w, "<div>%s</div>", "Username schon vorhanden")
		} else { //Wenn es den Benutzernamen noch nicht gibt
			os.Open("data/userdata/Test.csv")
			empData := [][]string{
				{Uname, pword}}
			csvFile, err := os.OpenFile("data/userdata/Test.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Fatalf("failed opening file: %s", err)
			}
			csvwriter := csv.NewWriter(csvFile)
			for _, empRow := range empData {
				_ = csvwriter.Write(empRow) //schreibe die neuen Logindaten in die Speicherdatei
			}
			csvwriter.Flush()
			csvFile.Close()
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "Name", Value: Uname, Expires: expiration}
			http.SetCookie(w, &cookie)
			cookie2 := http.Cookie{Name: "Password", Value: pword, Expires: expiration}
			http.SetCookie(w, &cookie2)

		}
	} else {
		fmt.Fprintf(w, "<div>%s</div>", "Username oder Passwort zu kurz")
	}
	http.Redirect(w, r, "/MainPage", 301) //Nach dem Registrieren geht es zur MainPage
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "Name",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour), // Set expires for older versions of IE
		Path:    "/"})
	http.SetCookie(w, &http.Cookie{
		Name:    "Password",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour), // Set expires for older versions of IE
		Path:    "/"})
	http.Redirect(w, r, "/Login", 301)
}

func Login(w http.ResponseWriter, r *http.Request) {
	Uname := r.FormValue("uname")
	pword := r.FormValue("pword")
	salt := "15967" //Wird zum Saling des Passwords genutzt
	if !Helper.FileExists("data/userdata/Test.csv") {
		fmt.Fprintf(w, "<div>%s</div>", "Keine User Vorhanden")
	} else {
		if len(Uname) > 0 && len(pword) > 0 {
			pword = pword + salt             //Salting
			pword = Helper.GetMD5Hash(pword) //Hashing
			if validateUser(Uname, pword) {  //Wenn es eine passende User/Passwort-Kombination gibt
				expiration := time.Now().Add(24 * time.Hour)
				cookie := http.Cookie{Name: "Name", Value: Uname, Expires: expiration}
				http.SetCookie(w, &cookie)
				cookie2 := http.Cookie{Name: "Password", Value: pword, Expires: expiration}
				http.SetCookie(w, &cookie2)
				http.Redirect(w, r, "/MainPage", 301) //Wird auf die Startseite weitergeleitet
			} else {
				fmt.Fprintf(w, "<div>%s</div>", "Keine passenden Einloggdaten gefunden") //Falls es keine passende Kombination gibt
			}

		} else {
			fmt.Fprintf(w, "<div>%s</div>", "Username oder Passwort zu kurz")
		} //Wenn kein Passwort und/oder Benutzername eingegeben wurde
	}
}
func validateUser(uname string, pword string) bool {
	lines, err := Helper.ReadCsv("data/userdata/Test.csv") //Auslesen aller Einloggdaten
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		data := User{
			username: line[0],
			password: line[1],
		}
		if data.username == uname && data.password == pword { //Wenn es eine passende User/Passwort-Kombination gibt
			return true
		}
	}
	return false
}
