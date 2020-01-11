package main

// MAtrikel-Nr 3736476, 8721083
import (
	"./Handler"
	"net/http"
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
	HttpTests(t)
	LoginTest(t)
}
