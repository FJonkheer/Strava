package Handler

// MAtrikel-Nr 3736476, 8721083
import (
	"encoding/csv"
	"fmt"
	"github.com/FJonkheer/Strava/src/Helper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//Die Upload-Funktion
func Uploader(w http.ResponseWriter, r *http.Request) {
	name, _ := r.Cookie("Name")
	password, _ := r.Cookie("Password")

	if !validateUser(name.Value, password.Value) {
		http.Redirect(w, r, "/Login", 301)
	} else {
		Pfad := "Files/" + name.Value + "/"          //jeder Benutzer hat seinen eigenen Dateispeicherort
		file, fileheader, err := r.FormFile("datei") //nimmt sich die Datei aus dem HTTP Request
		if err != nil {
			_, err := fmt.Fprintf(w, "<div>%s</div>", err)
			if err != nil {
				fmt.Println()
			}
		}
		filename := fileheader.Filename     //der Dateiname wird aus dem Header gesucht
		e, _ := Helper.FilePathExists(Pfad) //Existiert der Ort, an dem die Daten gespeichert werden sollen schon
		if !e {
			err := Helper.CreateFolders(Pfad) //Wenn die Ordner nicht existieren, sollen diese erstellt werden
			if err != nil {
				fmt.Println("Konnte Ordner nicht erstellen")
			}
		}
		defer func() {
			err = file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		fileBytes, _ := ioutil.ReadAll(file)     //Liest die Datei aus
		newFile, _ := os.Create(Pfad + filename) //Erstellt die Zieldatei
		_, err = newFile.Write(fileBytes)        //Beschreibt die Zieldatei mit den Daten der Ursprungsdatei
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
		err = newFile.Close()
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
		detectedFileType := http.DetectContentType(fileBytes) //Ermittelt die Dateiendung
		datei := filename
		if detectedFileType == "application/zip" { //Falls es sich um eine .zip handelt, muss diese noch entpackt werden
			datei, err = Helper.UnZip(Pfad+filename, Pfad) //Entpacken
			if err != nil {
				_, err = fmt.Fprintf(w, "<div>%s<div>", "Fehler beim entpacken: "+err.Error())
				err = file.Close()
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			datei = Pfad + datei
		}
		filedate, duration, distance, maxspeed, avgspeed, standtime := Helper.GetInfo(datei)
		idistance := fmt.Sprintf("%f", distance)
		imaxspeed := fmt.Sprintf("%f", maxspeed)
		iavgspeed := fmt.Sprintf("%f", avgspeed)
		vali := Helper.Validation(maxspeed, avgspeed, distance)
		if filedate == "" {
			currentTime := time.Now()
			filedate = currentTime.Format("2006-01-02")
		}
		activity := r.FormValue("types")  //liest den Aktivitätstypen aus dem http-Request
		comment := r.FormValue("comment") //liest den Benutzer-Kommentar
		switch vali {
		case "f":
			activity = "Fahrrad fahren"
			break
		case "l":
			activity = "Laufen"
			break
		default:
			break
		}
		//Speichern der Metadaten zu der hochgeladenen Datei
		empData := [][]string{
			{"date", "type", "comment", "duration", "distance", "maxspeed", "avgspeed", "standtime"},
			{filedate, activity, comment, duration.String(), idistance, imaxspeed, iavgspeed, standtime.String()}} //Die Informationen, die gespeichert werden müssen
		infofile, err := os.Create(datei + ".csv") //Erstellen der Infodatei
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		csvwriter := csv.NewWriter(infofile) //Beschreiben der CSV-Datei
		for _, empRow := range empData {
			err = csvwriter.Write(empRow)
			err = file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
		csvwriter.Flush()
		err = infofile.Close()
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(empData)
		http.Redirect(w, r, "/MainPage", 301)
	}
}
