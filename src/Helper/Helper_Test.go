package Helper

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test(t *testing.T) {
	MD5_Test(t)
	FileAccess_Test(t)
	GPX_Test(t)
	CalculationTest(t)
}
func MD5_Test(t *testing.T) {
	assert.Equal(t, GetMD5Hash("Hello"), "8b1a9953c4611296a827abf8c47804d7", "")
	assert.NotEqual(t, GetMD5Hash("Test"), "8b1a9953c4611296a827abf8c47804d7", "")

}
func GPX_Test(t *testing.T) {
	assert.Equal(t, CreateFolders("testfolder"), nil)
	var testfolder []string
	testfolder = append(testfolder, "existantfile.csv")
	testdatei := "testfolder/existantfile.csv"
	type Metadata struct {
		Date        string  `xml:"metadata>time"`
		Trackpoints []trkpt `xml:"trk>trkseg>trkpt"`
	}
	var TestMeta Metadata
	assert.Equal(t, GpxRead(testdatei), TestMeta)

}

func FileAccess_Test(t *testing.T) {
	//Vorbereitung
	assert.Equal(t, CreateFolders("testfolder"), nil)
	var testfolder []string
	testfolder = append(testfolder, "existantfile.csv")
	testdatei := "testfolder/existantfile.csv"
	empData := [][]string{
		{"test1", "test2"},
		{"value1", "value2"},
	}
	csvdatei, _ := os.Create(testdatei)
	csvwriter := csv.NewWriter(csvdatei)
	for _, empRow := range empData {
		_ = csvwriter.Write(empRow) //Und die Indexe geschrieben
	}
	csvwriter.Flush()
	csvdatei.Close()
	var teststruct UserFiles
	assert.Equal(t, FileExists("testfolder/notexistentfile"), false)
	//assert.Equal(t, interface{}(ReadCsv(testdatei)), empData)
	assert.Equal(t, FileExists(testdatei), true)
	assert.Equal(t, Scanforcsvfiles("testfolder"), testfolder)
	assert.Equal(t, Parsecsvtostruct(testdatei), teststruct)
	DeleteFiles(testdatei)
	assert.Equal(t, FileExists(testdatei), false)
	testfolder = append(testfolder, "test.zip")
	UnZip("test.zip", "testfolder")

}
func CalculationTest(t *testing.T) {
	assert.Equal(t, Validation(20, 20, 20), "f", "")
	assert.Equal(t, Validation(6, 10, 20), "l", "")
	assert.Equal(t, Latlongtodistance(0, 0, 0, 0, 0, 0), float64(0), "")
	assert.Equal(t, Latlongtodistance(10, 10, 20, 20, 10, 0), float64(1.5464880483491938e+06), "")

	type Metadata struct {
		Date        string  `xml:"metadata>time"`
		Trackpoints []trkpt `xml:"trk>trkseg>trkpt"`
	}
	assert.Equal(t, CreateFolders("testfolder"), nil)
	var testfolder []string
	testfolder = append(testfolder, "existantfile.csv")
	testdatei := "testfolder/existantfile.csv"
	a, b, c, d, e := calculateEverything(testdatei)
	assert.Equal(t, a, "0")
	assert.Equal(t, b, "0")
	assert.Equal(t, c, "0")
	assert.Equal(t, d, "0")
	assert.Equal(t, e, "0")
	assert.Equal(t, getDate(testdatei), "")
	f, g, h, i, j, k := GetInfo(testdatei)
	assert.Equal(t, f, "0")
	assert.Equal(t, g, "0")
	assert.Equal(t, h, "0")
	assert.Equal(t, i, "0")
	assert.Equal(t, j, "0")
	assert.Equal(t, k, "0")
}
