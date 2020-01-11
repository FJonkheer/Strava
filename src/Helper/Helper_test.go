package Helper

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test(t *testing.T) {
	MD5_Test(t)
	FileAccess_Test(t)
}
func MD5_Test(t *testing.T) {
	assert.Equal(t, GetMD5Hash("Hello"), "8b1a9953c4611296a827abf8c47804d7", "")
	assert.NotEqual(t, GetMD5Hash("Test"), "8b1a9953c4611296a827abf8c47804d7", "")
}

func FileAccess_Test(t *testing.T) {
	//Vorbereitung
	assert.Equal(t, CreateFolders("testfolder"), nil)
	testdatei := "testfolder/existantfile.csv"
	os.Create(testdatei)

	assert.Equal(t, Validation(20, 20, 20), "f", "")
	assert.Equal(t, Validation(6, 10, 20), "l", "")
	assert.Equal(t, Latlongtodistance(0, 0, 0, 0, 0, 0), float64(0), "")
	assert.Equal(t, Latlongtodistance(10, 10, 20, 20, 10, 0), float64(1.5464880483491938e+06), "")
	assert.Equal(t, FileExists("testfolder/notexistentfile"), false)

	//Aufräumen
	assert.Equal(t, FileExists(testdatei), true)
	assert.Equal(t, DeleteFiles(testdatei), nil)
}
