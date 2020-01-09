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

func Uploader(w http.ResponseWriter, r *http.Request) {
	Pfad := "Files/Username/"
	file, fileheader, err := r.FormFile("datei")
	if err != nil {
		fmt.Fprintf(w, "<div>%s</div>", err)
	}
	filename := fileheader.Filename
	e, _ := Helper.FilePathExists(Pfad)
	if !e {
		Helper.CreateFolders(Pfad)
	}
	defer file.Close()
	fileBytes, _ := ioutil.ReadAll(file)
	newFile, _ := os.Create(Pfad + filename)
	newFile.Write(fileBytes)
	newFile.Close()
	detectedFileType := http.DetectContentType(fileBytes)
	datei := filename
	if detectedFileType == "application/zip" {
		datei, err = unZip(Pfad+filename, Pfad)
		if err != nil {
			fmt.Fprintf(w, "<div>%s<div>", "Fehler beim entpacken: "+err.Error())
		}
	}

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")
	activity := r.FormValue("types")
	comment := r.FormValue("comment")
	empData := [][]string{
		{"uploaddate", "type", "comment"},
		{date, activity, comment}}
	infofile, err := os.Create(Pfad + datei + ".csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(infofile)
	for _, empRow := range empData {
		csvwriter.Write(empRow)
	}
	csvwriter.Flush()
	infofile.Close()

	http.Redirect(w, r, "/MainPage", 301)
}

/* Source = https://golangcode.com/unzip-files-in-go/ */
func unZip(src string, dir string) (string, error) {
	r, err := zip.OpenReader(src)
	defer r.Close()
	fpath := ""
	for _, f := range r.File {
		fpath = filepath.Join(dir, f.Name)
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
