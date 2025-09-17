package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// company map: key: companyName, value: employee-map
// employee map: key employee id, value: {billableRate, totalHours}

// insertion:
// - loop through csv file
// - update map with entry:
//		- if company exists
//			- if employee has worked for company
//				- update employee hours
//			- else:
//				add employee details
//		- else:
//			- add company and employee details

// retrieval
// - access map by company name
// - loop though employees
// - update output pdf with employee data
// - return

type CompanyMap map[string]EmployeeMap
type EmployeeMap map[int]Employee

type Employee struct {
	BillableRate float64
	TotalHours   float64
}

func readCSV(reader io.Reader) (CompanyMap, error) {
	cm := make(CompanyMap)
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = -1
	layout := "15:04"

	// skip header
	if _, err := csvReader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// read lines
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(record) < 6 {
			return nil, fmt.Errorf("invalid record: expected >=6 fields, got %d", len(record))
		}

		id, rate, companyName, _, startTime, endTime := record[0], record[1], record[2], record[3], record[4], record[5]

		// convert strings to required types
		eid, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil {
			return nil, fmt.Errorf("invalid employee id %q: %w", id, err)
		}

		bRate, err := strconv.ParseFloat(strings.TrimSpace(rate), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid billable rate %q: %w", rate, err)
		}

		start, err := time.Parse(layout, strings.TrimSpace(startTime))
		if err != nil {
			return nil, fmt.Errorf("invalid start time %q: %w", startTime, err)
		}
		end, err := time.Parse(layout, strings.TrimSpace(endTime))
		if err != nil {
			return nil, fmt.Errorf("invalid end time %q: %w", endTime, err)
		}

		diff := end.Sub(start).Hours()

		// update maps
		company, exists := cm[companyName]
		if exists {
			employee, exists := company[eid]
			if exists {
				employee.TotalHours += diff
				company[eid] = employee
			} else {
				employee = Employee{BillableRate: bRate, TotalHours: diff}
				company[eid] = employee
			}
		} else {
			em := make(EmployeeMap)
			em[eid] = Employee{BillableRate: bRate, TotalHours: diff}
			cm[companyName] = em
		}

	}

	return cm, nil
}
