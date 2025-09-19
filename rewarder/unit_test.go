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
	if err := InitDB(":memory:"); err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	
	csv := `Customer ID,Customer First Name,Order Value
1,Kweku,100
2,Abena,1200
3,Kojo,4800
4,Esi,7500
5,Yaw,15000
6,Akua,999
7,Mensah,10000
`
	vouchers, err := readCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("readCSV returned error: %v", err)
	}

	// assertions?
	// should get 5 vouchers (all except kweku=100 and akua=999)
	count := 5
	if len(vouchers) != count {
		t.Fatalf("expected %d vouchers, got %d", count, len(vouchers))
	}

expected := map[string][]any{
		"abena":  {100.0, 1},
		"kojo":   {100.0, 1},
		"esi":    {500.0, 5},
		"yaw":    {1000.0, 10},
		"mensah": {1000.0, 10},
	}

	for _, v := range vouchers {
		value, ok := expected[v.CustomerName]
		if !ok {
			t.Errorf("unexpected voucher for customer %s", v.CustomerName)
		}
		if v.Amount != value[0].(float64) {
			t.Errorf("expected voucher amount %.2f for customer %s, got %.2f", value, v.CustomerName, v.Amount)
		}
		// check validity? (expiry - creation)hours / 24 = days
		if int(v.ExpiresAt.Sub(v.CreatedAt).Hours()/24) != value[1].(int) {
			t.Errorf("expected voucher expiry %d days for customer %s, got %d days", value[1], v.CustomerName, int(v.ExpiresAt.Sub(v.CreatedAt).Hours()/24))
		}
		if v.ExpiresAt.Before(v.CreatedAt) {
			t.Errorf("voucher expiry is before creation date for customer %s", v.CustomerName)
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

	csv := `Customer ID,Customer First Name,Order Value
1,Kweku,100
2,Abena,1200
3,Kojo,4800
4,Esi,7500
5,Yaw,15000
6,Akua,999
7,Mensah,10000
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
