package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

	em, ok := cm["acme"]
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

func TestUploadCSV(t *testing.T) {
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

	// test through router to make sure router + middleware are working
	rr := httptest.NewRecorder()
	router := Router()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", rr.Code, rr.Body.String())
	}
}

func TestGenerateInvoice_NoCompany(t *testing.T) {
	companyMap = CompanyMap{}
	_, err := generateInvoice("SomeCompany")
	if err == nil {
		t.Fatal("expected error for missing company, got nil")
	}
}

func TestGenerateInvoice_Success(t *testing.T) {
	companyMap = CompanyMap{
		"google": {
			1: {BillableRate: 100, TotalHours: 5},
		},
	}

	filename, err := generateInvoice("Google")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defer os.Remove(filename)

	// check that file exists when created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("expected file %v to exist", filename)
	}
}

func TestDownloadInvoice(t *testing.T) {
	companyMap = CompanyMap{
		"netflix": {
			1: {BillableRate: 100, TotalHours: 5},
		},
	}

	req := httptest.NewRequest("GET", "/api/download/Netflix", nil)

	// test through router to make sure router + middleware are working
	rr := httptest.NewRecorder()
	router := Router()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", rr.Code, rr.Body.String())
	}

	ct := rr.Header().Get("Content-Type")
	if ct != "application/pdf" {
		t.Fatalf("expected application/pdf, got %s", ct)
	}
}
