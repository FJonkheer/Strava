package main

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
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

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/action_page.php", login)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func login(w http.ResponseWriter, r *http.Request) {
	uname := r.FormValue("uname")
	password := r.FormValue("pword")
	salt := "15967"
	isGone := false
	if len(uname) > 0 && len(password) > 0 {
		lines, err := ReadCsv("Test.csv")
		if err != nil {
			panic(err)
		}

		for _, line := range lines {
			data := CsvLine{
				Column1: line[0],
				Column2: line[1],
			}
			if data.Column1 == uname {
				isGone = true
			}

		}
		if isGone {
			fmt.Fprintf(w, "<div>%s</div>", "Username schon vorhanden")
		} else {
			password = password + salt
			pword := GetMD5Hash(password)
			os.Open("Test.csv")
			if fileExists("Test.csv") {
				empData := [][]string{
					{uname, pword}}
				csvFile, err := os.OpenFile("Test.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					log.Fatalf("failed opening file: %s", err)
				}
				csvwriter := csv.NewWriter(csvFile)
				for _, empRow := range empData {
					_ = csvwriter.Write(empRow)
				}
				csvwriter.Flush()
				csvFile.Close()
			} else {
				password = password + salt
				pword := GetMD5Hash(password)
				empData := [][]string{
					{"uname", "pword"},
					{uname, pword}}
				csvFile, err := os.Create("Test.csv")
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
		}

	} else {
		fmt.Fprintf(w, "<div>%s</div>", "Username oder Passwort zu kurz")
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

type CsvLine struct {
	Column1 string
	Column2 string
}

func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
