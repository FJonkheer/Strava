package Handler

import (
	"Helper"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
)

type User struct { //Die Struktur eines Benutzers
	username string
	password string
}

func Handling(w http.ResponseWriter, r *http.Request) { //Wurde "Login" oder "Registrieren" gedrückt?
	if r.FormValue("login") == "Login" {
		login(w, r)
	} else {
		register(w, r)
	}
}

//Ich hab die Register-Funktion mal ein wenig verschönert
/*func register(w http.ResponseWriter, r *http.Request) { //Der Ablauf für Registrieren
	uname := r.FormValue("uname")
	password := r.FormValue("pword")
	salt := "15967" //Der Salt, welcher zum Verschlüsseln genutzt wird
	isGone := false
	if len(uname) > 0 && len(password) > 0 { //Wurde ein Benutzername und ein Passwort eingegeben?
		password = password + salt //Salting des Passworts
		pword := Helper.GetMD5Hash(password) //Hashing des Passworts mit Salt
		t, _ := Helper.FilePathExists("Strava/data/userdata") //gibt es den Dateispeicherpfad schon?
		if t {
			if Helper.FileExists("Strava/data/userdata/Test.csv") { //existiert die Speicherdatei schon?
				lines, err := Helper.ReadCsv("Strava/data/userdata/Test.csv") //CSV-Auslesen
				if err != nil {
					panic(err)
				}
				for _, line := range lines { //Prüfen ob es den Benutzernamen schon in der Speicherdatei gibt
					data := User{
						username: line[0],
						password: line[1],
					}
					if data.username == uname {
						isGone = true
					}
				}
				if isGone {
					fmt.Fprintf(w, "<div>%s</div>", "Username schon vorhanden")
				} else { //Wenn es den Benutzernamen noch nicht gibt
					os.Open("Strava/data/userdata/Test.csv")
					empData := [][]string{
						{uname, pword}}
					csvFile, err := os.OpenFile("Strava/data/userdata/Test.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
					if err != nil {
						log.Fatalf("failed opening file: %s", err)
					}
					csvwriter := csv.NewWriter(csvFile)
					for _, empRow := range empData {
						_ = csvwriter.Write(empRow) //schreibe die neuen Logindaten in die Speicherdatei
					}
					csvwriter.Flush()
					csvFile.Close()
				}
			} else { //Wenn die Speicherdatei noch nicht existiert
				empData := [][]string{
					{"uname", "pword"},
					{uname, pword}}
				csvFile, err := os.Create("Strava/data/userdata/Test.csv") //Wird eine Speicherdatei erstellt
				if err != nil {
					log.Fatalf("failed creating file: %s", err)
				}
				csvwriter := csv.NewWriter(csvFile)
				for _, empRow := range empData {
					_ = csvwriter.Write(empRow) //Und die neuen Logindaten in die Datei geschrieben
				}
				csvwriter.Flush()
				csvFile.Close()
			}
		} else { //Falls es die Ordner noch nicht gibt
			Helper.CreateFolders("Strava/data/userdata") //Werden erst noch die neuen Ordner erstellt

			empData := [][]string{
				{"uname", "pword"},
				{uname, pword}}
			csvFile, err := os.Create("Strava/data/userdata/Test.csv") //wenn es noch keinen Ordner gibt, wird auch eine neue Speicherdatei erstellt
			if err != nil {
				log.Fatalf("failed creating file: %s", err)
			}
			csvwriter := csv.NewWriter(csvFile)
			for _, empRow := range empData {
				_ = csvwriter.Write(empRow)
			}
			csvwriter.Flush()
			csvFile.Close()
		}

		http.Redirect(w, r, "/MainPage", 301)
	} else {
		fmt.Fprintf(w, "<div>%s</div>", "Username oder Passwort zu kurz")
	}
}
*/

func register(w http.ResponseWriter, r *http.Request) {
	uname := r.FormValue("uname")
	password := r.FormValue("pword")
	salt := "15967" //Der Salt, welcher zum Verschlüsseln genutzt wird
	isGone := false
	if len(uname) > 0 && len(password) > 0 { //Wurde ein Benutzername und ein Passwort eingegeben?
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
			if data.username == uname {
				isGone = true
			}
		}
		if isGone {
			fmt.Fprintf(w, "<div>%s</div>", "Username schon vorhanden")
		} else { //Wenn es den Benutzernamen noch nicht gibt
			os.Open("data/userdata/Test.csv")
			empData := [][]string{
				{uname, pword}}
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
		}
	} else {
		fmt.Fprintf(w, "<div>%s</div>", "Username oder Passwort zu kurz")
	}
	http.Redirect(w, r, "/MainPage", 301) //Nach dem Registrieren geht es zur MainPage
}

func login(w http.ResponseWriter, r *http.Request) {
	uname := r.FormValue("uname")
	pword := r.FormValue("pword")
	salt := "15967" //Wird zum Saling des Passwords genutzt
	if !Helper.FileExists("data/userdata/Test.csv") {
		fmt.Fprintf(w, "<div>%s</div>", "Keine User Vorhanden")
	} else {
		if len(uname) > 0 && len(pword) > 0 {
			pword = pword + salt                                   //Salting
			pword = Helper.GetMD5Hash(pword)                       //Hashing
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
					http.Redirect(w, r, "/MainPage", 301) //Wird auf die Startseite weitergeleitet
				}
			}
			fmt.Fprintf(w, "<div>%s</div>", "Keine passenden Einloggdaten gefunden") //Falls es keine passende Kombination gibt
		} else {
			fmt.Fprintf(w, "<div>%s</div>", "Username oder Passwort zu kurz") //Wenn kein Passwort und/oder Benutzername eingegeben wurde
		}
	}
}
