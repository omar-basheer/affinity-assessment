package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadCSV(t *testing.T) {
	csv := `"Employee ID","Billable Rate (per hour)","Project","Date","Start time","End Time"
"1","100","Acme","7/1/19","09:00","11:00"
"1","100","Acme","7/1/19","12:00","14:00"
"2","150","Acme","7/1/19","10:00","15:00"
`
	cm, err := readCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("readCSV returned error: %v", err)
	}

	// assertions?
	if len(cm) != 1 {
		t.Fatalf("expected 1 company, got %d", len(cm))
	}

	em, ok := cm["Acme"]
	if !ok {
		t.Fatalf("company Acme missing")
	}

	if e, ok := em[1]; !ok {
		t.Fatalf("employee 1 missing")
	} else {
		// employee 1 worked 2 + 2 = 4 hours
		if int(e.TotalHours) != 4 {
			t.Fatalf("expected employee 1 total 4 hours, got %v", e.TotalHours)
		}
	}

	if e2, ok := em[2]; !ok {
		t.Fatalf("employee 2 missing")
	} else {
		if int(e2.TotalHours) != 5 {
			t.Fatalf("expected employee 2 total 5 hours, got %v", e2.TotalHours)
		}
	}
}

func TestUploadHandler(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}

	csv := `"Employee ID","Billable Rate (per hour)","Project","Date","Start time","End Time"
"1","100","Acme","7/1/19","09:00","11:00"
`
	if _, err := io.Copy(fw, bytes.NewBufferString(csv)); err != nil {
		t.Fatal(err)
	}
	w.Close()

	req := httptest.NewRequest("POST", "/api/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	rr := httptest.NewRecorder()
	router := Router() // make sure Router() is deterministic in tests
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", rr.Code, rr.Body.String())
	}
}
