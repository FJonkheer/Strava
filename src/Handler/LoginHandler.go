package Handler

import (
	"Helper"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
)

type User struct {
	username string
	password string
}

func Handling(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("login") == "Login" {
		login(w, r)
	} else {
		register(w, r)
	}
}
func register(w http.ResponseWriter, r *http.Request) {
	uname := r.FormValue("uname")
	password := r.FormValue("pword")
	salt := "15967"
	isGone := false
	if len(uname) > 0 && len(password) > 0 {
		password = password + salt
		pword := Helper.GetMD5Hash(password)
		t, _ := Helper.FilePathExists("Strava/data/userdata")
		if t {
			if Helper.FileExists("Strava/data/userdata/Test.csv") {
				lines, err := Helper.ReadCsv("Strava/data/userdata/Test.csv")
				if err != nil {
					panic(err)
				}
				for _, line := range lines {
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
				} else {
					os.Open("Strava/data/userdata/Test.csv")
					empData := [][]string{
						{uname, pword}}
					csvFile, err := os.OpenFile("Strava/data/userdataTest.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
					if err != nil {
						log.Fatalf("failed opening file: %s", err)
					}
					csvwriter := csv.NewWriter(csvFile)
					for _, empRow := range empData {
						_ = csvwriter.Write(empRow)
					}
					csvwriter.Flush()
					csvFile.Close()
				}
			} else {
				empData := [][]string{
					{"uname", "pword"},
					{uname, pword}}
				csvFile, err := os.Create("Strava/data/userdata/Test.csv")
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
		} else {
			Helper.CreateFolders("Strava/data/userdata")

			empData := [][]string{
				{"uname", "pword"},
				{uname, pword}}
			csvFile, err := os.Create("Strava/data/userdata/Test.csv")
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
func login(w http.ResponseWriter, r *http.Request) {
	uname := r.FormValue("uname")
	pword := r.FormValue("pword")
	salt := "15967"
	success := false
	if Helper.FileExists("Strava/data/userdata/Test.csv") {
		if len(uname) > 0 && len(pword) > 0 {
			pword = pword + salt
			pword = Helper.GetMD5Hash(pword)
			lines, err := Helper.ReadCsv("data/userdata/Test.csv")
			if err != nil {
				panic(err)
			}
			for _, line := range lines {
				data := User{
					username: line[0],
					password: line[1],
				}
				if data.username == uname && data.password == pword {
					//Redirect
					http.Redirect(w, r, "/MainPage", 301)
					break
				}
			}
			if !success {
				fmt.Fprintf(w, "<div>%s</div>", "Keine passenden Einloggdaten gefunden")
			}
		} else {
			fmt.Fprintf(w, "<div>%s</div>", "Username oder Passwort zu kurz")
		}
	} else {
		fmt.Fprintf(w, "<div>%s</div>", "Keine User Vorhanden")
	}
}
