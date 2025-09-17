package main

import (
	"log"
	"net/http"
)

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

func download(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
}
