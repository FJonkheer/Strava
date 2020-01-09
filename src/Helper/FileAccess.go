package Helper

import (
	"encoding/csv"
	"net/http"
	"os"
)

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
