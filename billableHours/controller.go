package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// upload handles file uploads via multipart/form-data
func upload(w http.ResponseWriter, r *http.Request) {
	// limit request body to 10MB to avoid huge uploads
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB

	// parse form to get file
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		RespondWithError(w, 400, "failed to parse multipart form", err.Error())
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		RespondWithError(w, 400, "failed to parse file", err.Error())
		return
	}
	defer file.Close()

	//fmt.Printf("uploaded File: %+v\n", header)

	// process csv from upload
	cm, err := readCSV(file)
	if err != nil {
		RespondWithError(w, 400, "failed to read file", err.Error())
		return
	}

	err = RespondWithJSON(w, 200, "file uploaded successfully", cm)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

// download generates and serves an invoice PDF for the specified company
func download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyName := vars["companyName"]
	if companyName == "" {
		RespondWithError(w, 400, "invalid companyName", "Please provide a valid company name")
		return
	}

	filename, err := generateInvoice(companyName)
	if err != nil {
		RespondWithError(w, 400, "failed to generate invoice", err.Error())
		return
	}

	// serve file as download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, filename)
}
