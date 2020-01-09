package Handler

import (
	"Helper"
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//Die Upload-Funktion
func Uploader(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {
		Pfad := "Files/" + cookie.Value + "/"        //jeder Benutzer hat seinen eigenen Dateispeicherort
		file, fileheader, err := r.FormFile("datei") //nimmt sich die Datei aus dem HTTP Request
		if err != nil {
			fmt.Fprintf(w, "<div>%s</div>", err)
		}
		filename := fileheader.Filename     //der Dateiname wird aus dem Header gesucht
		e, _ := Helper.FilePathExists(Pfad) //Existiert der Ort, an dem die Daten gespeichert werden sollen schon
		if !e {
			Helper.CreateFolders(Pfad) //Wenn die Ordner nicht existieren, sollen diese erstellt werden
		}
		defer file.Close()
		fileBytes, _ := ioutil.ReadAll(file)     //Liest die Datei aus
		newFile, _ := os.Create(Pfad + filename) //Erstellt die Zieldatei
		newFile.Write(fileBytes)                 //Beschreibt die Zieldatei mit den Daten der Ursprungsdatei
		newFile.Close()
		detectedFileType := http.DetectContentType(fileBytes) //Ermittelt die Dateiendung
		datei := filename
		if detectedFileType == "application/zip" { //Falls es sich um eine .zip handelt, muss diese noch entpackt werden
			datei, err = unZip(Pfad+filename, Pfad) //Entpacken
			if err != nil {
				fmt.Fprintf(w, "<div>%s<div>", "Fehler beim entpacken: "+err.Error())
			}
		} else {
			datei = Pfad + datei
		}

	//Speichern der Metadaten zu der hochgeladenen Datei
	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")
	activity := r.FormValue("types")  //liest den Aktivitätstypen aus dem http-Request
	comment := r.FormValue("comment") //liest den Benutzer-Kommentar
	empData := [][]string{
		{"uploaddate", "type", "comment"},
		{date, activity, comment}} //Die Informationen, die gespeichert werden müssen
	infofile, err := os.Create(datei + ".csv") //Erstellen der Infodatei
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(infofile) //Beschreiben der CSV-Datei
	for _, empRow := range empData {
		csvwriter.Write(empRow)
	}
	csvwriter.Flush()
	infofile.Close()
	fmt.Println(Helper.CalculateEverything(datei)) //Auslesen der GPX-Datei, muss eventuell verschoben werden, hat hier keinen Sinn
	http.Redirect(w, r, "/MainPage", 301)
}

/* Source = https://golangcode.com/unzip-files-in-go/ */
func unZip(src string, dir string) (string, error) {
	r, err := zip.OpenReader(src)
	defer r.Close()
	fpath := ""
	for _, f := range r.File {
		fpath = dir + f.Name
		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
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
		outFile.Close()
		rc.Close()
		if err != nil {
			return "", err
		}
	}
	return fpath, nil
}
