package Helper

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//Trackpoint-Struktur
type trkpt struct {
	Latitude  string `xml:"lat,attr"`
	Longitude string `xml:"lon,attr"`
	Time      string `xml:"time"`
}

//Metadaten und alle Trackpoints der gpx-Datei
type metadata struct {
	Date        string  `xml:"metadata>time"`
	Trackpoints []trkpt `xml:"trk>trkseg>trkpt"`
}

//einlesen der Gpx-Datei
func Gpxread(path string) {
	xmlFile, err := os.Open(path) //öffnen des GPX-Files
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("XML successfully opened")  //Ausgabe, kann weg
	defer xmlFile.Close()                   //schließen des GPX-Files
	byteValue, _ := ioutil.ReadAll(xmlFile) //den Inhalt des GPX-Files auslesen

	var Run metadata //ein Objekt der Struktur metadata erstellen

	err = xml.Unmarshal(byteValue, &Run)        //unmarshal liest den Inhalt des GPX-Files aus und sortiert die Einträge in das mitgegebene Objekt
	Run.Date = strings.Split(Run.Date, "T")[0]  //Das Datum von der Uhrzeit trennen
	fmt.Println("Date: " + Run.Date)            //Ausgabe, kann weg
	for i := 0; i < len(Run.Trackpoints); i++ { //von jedem Trackpoint die Zeit formatieren
		Run.Trackpoints[i].Time = strings.Split(Run.Trackpoints[i].Time, "T")[1]
		Run.Trackpoints[i].Time = strings.Replace(Run.Trackpoints[i].Time, "Z", "", -1)

		//Ausgabe, kann dann irgendwann weg
		fmt.Println("Trackpoint Time: " + Run.Trackpoints[i].Time)
		fmt.Println("Trackpoint longitude: " + Run.Trackpoints[i].Longitude)
		fmt.Println("Trackpoint latitude: " + Run.Trackpoints[i].Latitude + "\n")
	}
}