package Helper

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Filename  string
	Filedate  string
	Activity  string
	Comment   string
	Duration  string
	Distance  string
	Maxspeed  string
	Avgspeed  string
	Standtime string
}

type UserFiles struct {
	Username string
	Files    []File
}

func FileExists(filename string) bool { //Abfrage, ob eine Datei bereits existiert
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FilePathExists(path string) (bool, error) { //Abfrage, ob der Pfad bereits existiert
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func CreateFolders(path string) error { //Alle Ordner des Pfads erstellen
	return os.MkdirAll(path, os.ModePerm)
}

func ReadCsv(filename string) ([][]string, error) { //Eine CSV-Datei auslesen

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer func() {
		err = f.Close()
	}()
	if err != nil {
		fmt.Println(err)
	}

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func DeleteFiles(path string) error { //Löschen von GPX-Datei und zugehörigen Dateien

	path = strings.Replace(path, ".csv", "", -1)
	err := os.Remove(path) //Löschen der GPX-Datei
	if err != nil {
		return err
	}
	err = os.Remove(path + ".csv") //Die Infodatei soll auch entfernt werden
	if err != nil {
		return err
	}
	exists := FileExists(path + ".zip")
	if exists {
		err = os.Remove(path + ".zip") //wenn eine zugehörige ZIP-Datei existiert soll auch diese gelöscht werden
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadFile(w http.ResponseWriter, r *http.Request, path string) { //Herunterladen einer Datei
	path = strings.Replace(path, ".csv", "", -1)
	exists := FileExists(path + ".zip") //Falls eine zugehörige ZIP existiert, soll diese heruntergeladen werden
	if exists {
		path = path + ".zip"
	}
	file := strings.Split(path, "/")[2]
	w.Header().Set("Content-Disposition", "attachment; filename="+file)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	e := strings.NewReader(path)
	_, err := io.Copy(w, e)
	if err != nil {
		fmt.Println(err)
	}
}

func ChangeInfoFile(w http.ResponseWriter, r *http.Request, file string) {
	//Speichern der Metadaten zu der hochgeladenen Datei
	content, err := ReadCsv(file)
	var onefile File
	onefile.Filedate = content[1][0]
	onefile.Duration = content[1][3]
	onefile.Distance = content[1][4]
	onefile.Maxspeed = content[1][5]
	onefile.Avgspeed = content[1][6]
	onefile.Standtime = content[1][7]
	onefile.Activity = r.FormValue("types")  //liest den Aktivitätstypen aus dem http-Request
	onefile.Comment = r.FormValue("comment") //liest den Benutzer-Kommentar
	empData := [][]string{
		{"date", "type", "comment", "duration", "distance", "maxspeed", "avgspeed", "standtime"},
		{onefile.Filedate, onefile.Activity, onefile.Comment, onefile.Duration, onefile.Distance, onefile.Maxspeed, onefile.Avgspeed, onefile.Standtime}} //Die Informationen, die gespeichert werden müssen
	infofile, err := os.Create(file) //Erstellen der Infodatei
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(infofile) //Beschreiben der CSV-Datei
	for _, empRow := range empData {
		err = csvwriter.Write(empRow)
		if err != nil {
			fmt.Println(err)
		}
	}
	csvwriter.Flush()
	err = infofile.Close()
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/Review", 301)
}

func Scanforcsvfiles(path string) []string {
	var csvfiles []string
	/*err := filepath.Walk(path, func(filepath string, info os.FileInfo, err error) error {
		files = append(files, filepath)
		return nil
	})*/
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".csv" {
			csvfiles = append(csvfiles, file.Name())
		}
	}
	return csvfiles
}

func Parsecsvtostruct(path string) UserFiles {
	var user UserFiles
	csvfiles := Scanforcsvfiles(path)
	var onefile File
	for _, file := range csvfiles {
		content, _ := ReadCsv(path + "/" + file)
		onefile.Filename = file
		onefile.Filedate = content[1][0]
		onefile.Activity = content[1][1]
		onefile.Comment = content[1][2]
		onefile.Duration = content[1][3]
		onefile.Distance = content[1][4]
		onefile.Maxspeed = content[1][5]
		onefile.Avgspeed = content[1][6]
		onefile.Standtime = content[1][7]
		user.Files = append(user.Files, onefile)
	}
	return user
}

func GetfileName(r *http.Request) string {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	file := ""
	for key := range r.Form {
		if strings.Contains(key, ".gpx") {
			file = key
			break
		}
	}
	return file
}

/* Source = https://golangcode.com/unzip-files-in-go/ */
func UnZip(src string, dir string) (string, error) {
	r, err := zip.OpenReader(src)
	defer func() {
		err = r.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	fpath := ""
	for _, f := range r.File {
		fpath = dir + f.Name
		if f.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		_, err = io.Copy(outFile, rc)
		err = outFile.Close()
		if err != nil {
			fmt.Println(err)
		}
		err = rc.Close()
		if err != nil {
			fmt.Println(err)
		}
		if err != nil {
			return "", err
		}
	}
	return fpath, nil
}
