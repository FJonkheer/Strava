package Tests

import (
	"Helper"
	"github.com/stretchr/testify/assert"
	"net/http"

	"Handler"
	"net/http/httptest"
	"testing"
)

func HttpTests(t *testing.T) {

	req, err := http.NewRequest("Post", "/redirect", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler.Redirecting)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
func LoginTest(t *testing.T) {
	req, err := http.NewRequest("Post", "/Login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("uname", "Test")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler.Login)

	handler.ServeHTTP(rr, req)
	if req.Header.Get("uname") != "Test" {
		t.Errorf("handler returned wrong header: got %v want %v",
			req.Header.Get("uname"), "Test")
	}

}

func Test(t *testing.T) {
	assert.Equal(t, Helper.GetMD5Hash("Hello"), "8b1a9953c4611296a827abf8c47804d7", "")
	assert.NotEqual(t, Helper.GetMD5Hash("Test"), "8b1a9953c4611296a827abf8c47804d7", "")
	assert.Equal(t, Helper.Validation(20, 20, 20), "f", "")
	assert.Equal(t, Helper.Validation(6, 10, 20), "l", "")
	assert.Equal(t, Helper.Latlongtodistance(0, 0, 0, 0), float64(0), "")
	assert.Equal(t, Helper.Latlongtodistance(10, 10, 20, 20), float64(1.5464880483491938e+06), "")
	HttpTests(t)
	LoginTest(t)
}
