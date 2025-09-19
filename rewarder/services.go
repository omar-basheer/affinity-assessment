package main

import (
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

// database:
// - Voucher:
//	- id (string, primary key)
//	- customer name
//	- order value
//	- amount
//	- creation date
//	- expiry date (creation date + validity)
//
// insertion:
// - loop through csv file
// - create voucher entry:
//	- if order value >= 1000 and order value < 5000:
//		- createVoucher(amount:100, validity:1)
//	- else if order value >= 5000 and order value < 10000:
//		- createVoucher(amount:500, validity:5)
//	- else if order value >= 10000:
//		- createVoucher(amount:1000, validity:10)
//	- else:
//		- no voucher
//
// retrieval:
// - lookup voucher by id
// - if current date > expiry date:
//	- return voucher expired
// - else:
//	- return voucher details (customer name, voucher amount, expiry date?)

type Voucher struct {
	ID           string
	CustomerID   int
	CustomerName string
	OrderValue   float64
	Amount       float64
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

func createVoucher(amount float64, validity int, customerID int, customerName string, orderValue float64) Voucher {
	// generate unique id
	id := uuid.New().String()
	createdAt := time.Now()
	expiresAt := createdAt.AddDate(0, 0, validity)

	v := Voucher{
		ID:           id,
		CustomerID:   customerID,
		CustomerName: customerName,
		OrderValue:   orderValue,
		Amount:       amount,
		CreatedAt:    createdAt,
		ExpiresAt:    expiresAt,
	}

	// store voucher in database
	_, err := DB.Exec(`
		INSERT INTO vouchers (id, customer_id, customer_name, order_value, amount, created_at, expires_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, v.ID, v.CustomerID, v.CustomerName, v.OrderValue, v.Amount, v.CreatedAt, v.ExpiresAt)

	if err != nil {
		log.Printf("failed to create voucher: %v\n", err)
	}

	//fmt.Printf("voucher created: %v\n", v)
	return v
}

func readCSV(reader io.Reader) ([]Voucher, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	// skip header
	if _, err := csvReader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	var vouchers []Voucher

	// read lines
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(record) < 3 {
			return nil, fmt.Errorf("invalid record: expected >=3 fields, got %d", len(record))
		}

		id, fName, orderValue := record[0], record[1], record[2]

		// format data from file
		// todo: maybe check the order value is greater than 999 (so we'll create a voucher) before formatting all the values?
		cid, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil {
			return nil, fmt.Errorf("invalid customer ID: %s", id)
		}

		fName = strings.TrimSpace(strings.ToLower(fName))

		value, err := strconv.ParseFloat(strings.TrimSpace(orderValue), 64)

		// create vouchers
		if value < 1000 {
			continue
		} else if value >= 1000 && value < 5000 {
			vouchers = append(vouchers, createVoucher(100, 1, cid, fName, value))
		} else if value >= 5000 && value < 10000 {
			vouchers = append(vouchers, createVoucher(500, 5, cid, fName, value))
		} else if value >= 10000 {
			vouchers = append(vouchers, createVoucher(1000, 10, cid, fName, value))
		}
	}

	return vouchers, nil
}
