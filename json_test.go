package json

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/rkusa/web"
)

type testType struct {
	Foo string `json:"foo"`
}

func TestRead(t *testing.T) {
	app := web.New()
	app.Use(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var test testType
		if err := Read(r, &test); err != nil {
			t.Error(err)
		}

		if test.Foo != "bar" {
			t.Errorf("decode failed")
		}

		rw.WriteHeader(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	payload := `{"foo":"bar"}`

	r, err := http.NewRequest("POST", "/", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	app.ServeHTTP(rec, r)

	if rec.Code != http.StatusNoContent {
		t.Errorf("request failed")
	}
}

func TestWrite(t *testing.T) {
	app := web.New()
	app.Use(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		Write(rw, testType{"bar"})
	})

	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, (*http.Request)(nil))

	if rec.Code != http.StatusOK {
		t.Errorf("response failed")
	}

	if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("content-type not set correctly, got: %s", rec.Header().Get("Content-Type"))
	}

	if rec.Body.String() != `{"foo":"bar"}` {
		t.Errorf("wrong body, got: %s", rec.Body)
	}
}
