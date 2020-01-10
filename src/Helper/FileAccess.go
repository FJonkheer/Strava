package Helper

import (
	"encoding/csv"
	"net/http"
	"os"
	"path/filepath"
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
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func DeleteFiles(path string) error { //Löschen von GPX-Datei und zugehörigen Dateien
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
	exists := FileExists(path + ".zip") //Falls eine zugehörige ZIP existiert, soll diese heruntergeladen werden
	if exists {
		path = path + ".zip"
	}
	http.ServeFile(w, r, path)
}

/*func ChangeInfoFile(w http.ResponseWriter, r *http.Request, file string) {
	path := "Files/" + Handler.Uname //Benutzernamenabfrage
	//Speichern der Metadaten zu der hochgeladenen Datei
	content, err := ReadCsv(path + file)

	date := content[1][0]
	activity := r.FormValue("types")  //liest den Aktivitätstypen aus dem http-Request
	comment := r.FormValue("comment") //liest den Benutzer-Kommentar
	empData := [][]string{
		{"uploaddate", "type", "comment"},
		{date, activity, comment}} //Die Informationen, die gespeichert werden müssen
	infofile, err := os.Create(path + file + ".csv") //Erstellen der Infodatei
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(infofile) //Beschreiben der CSV-Datei
	for _, empRow := range empData {
		csvwriter.Write(empRow)
	}
	csvwriter.Flush()
	infofile.Close()
	http.Redirect(w, r, "/Review", 301)
}
*/

func Scanforcsvfiles(path string) []string {
	var files []string
	var csvfiles []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if filepath.Ext(file) == ".csv" {
			csvfiles = append(csvfiles, file)
		}
	}
	return csvfiles
}

func Parsecsvtostruct(username string) UserFiles {
	var user UserFiles
	user.Username = username
	path := "/Files/" + username + "/"
	csvfiles := Scanforcsvfiles(path)
	for i, file := range csvfiles {
		content, _ := ReadCsv(file)
		user.Files[i].Filename = file
		user.Files[i].Filedate = content[1][0]
		user.Files[i].Activity = content[1][1]
		user.Files[i].Comment = content[1][2]
		user.Files[i].Duration = content[1][3]
		user.Files[i].Distance = content[1][4]
		user.Files[i].Maxspeed = content[1][5]
		user.Files[i].Avgspeed = content[1][6]
		user.Files[i].Standtime = content[1][7]
	}
	return user
}
