package Handler

import (
	"Helper"
	"html/template"
	"net/http"
	"os"
	"reflect"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) { //Download einer Datei
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + Uname    //Benutzername muss abgefragt werden
		file := r.FormValue("File") //Das Feld, wo die Datei ausgewählt wurde
		path = path + file
		Helper.DownloadFile(w, r, path)
		http.Redirect(w, r, "/Review", 301)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) { //Löschen eines Eintrags
	//TODO: Abfrage ob wirklich gelöscht werden soll
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + Uname    //Benutzername muss abgefragt werden
		file := r.FormValue("File") //Das Feld, wo die Datei ausgewählt wurde
		path = path + file
		Helper.DeleteFiles(path)
		http.Redirect(w, r, "/Review", 301)
	}
}

func ChangeHandler(w http.ResponseWriter, r *http.Request) { //Ändern der InfoPage
	cookie, _ := r.Cookie(Uname)
	if cookie == nil {
		http.Redirect(w, r, "/Login", 301)
	} else {
		path := "Files/" + Uname    //Benutzername muss abgefragt werden
		file := r.FormValue("File") //Das Feld, wo die Datei ausgewählt wurde
		path = path + file
		Helper.ChangeInfoFile(w, r, path)
		http.Redirect(w, r, "/Review", 301)
	}
}

//Hier wird noch ausprobiert
//die Review-Page muss wahrscheinlich über so ein Template-Zeug gemacht werden

type Node struct {
	Contact_id  int
	Employer_id int
	First_name  string
	Middle_name string
	Last_name   string
}

var templateFuncs = template.FuncMap{"rangeStruct": RangeStructer}

// In the template, we use rangeStruct to turn our struct values
// into a slice we can iterate over
var htmlTemplate = `{{range .}}<tr>
{{range rangeStruct .}} <td>{{.}}</td>
{{end}}</tr>
{{end}}`

func Reviewer(w http.ResponseWriter, r *http.Request) {
	container := []Node{
		{1, 12, "Accipiter", "ANisus", "Nisus"},
		{2, 42, "Hello", "my", "World"},
	}

	// We create the template and register out template function
	t := template.New("t").Funcs(templateFuncs)
	t, err := t.Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, container)
	if err != nil {
		panic(err)
	}

}

// RangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}
