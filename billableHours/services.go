package main

import (
	"codeberg.org/go-pdf/fpdf"
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

var companyMap CompanyMap

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

		// format data from file
		eid, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil {
			return nil, fmt.Errorf("invalid employee id %q: %w", id, err)
		}

		bRate, err := strconv.ParseFloat(strings.TrimSpace(rate), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid billable rate %q: %w", rate, err)
		}

		cName := strings.TrimSpace(strings.ToLower(companyName))

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
		company, exists := cm[cName]
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
			cm[cName] = em
		}

	}

	companyMap = cm
	return cm, nil
}

func generateInvoice(companyName string) (string, error) {
	cName := strings.TrimSpace(strings.ToLower(companyName))

	employees, exists := companyMap[cName]
	if !exists {
		return "", fmt.Errorf("company %s not found", cName)
	}

	// create new pdf file
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// file header
	pdf.SetFont("Arial", "", 16)
	pdf.Cell(40, 10, "Company: "+companyName)
	pdf.Ln(20)

	// table header
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(0, 102, 204)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(40, 10, "Employee ID", "TBR", 0, "", true, 0, "")
	pdf.CellFormat(40, 10, "Number of Hours", "1", 0, "", true, 0, "")
	pdf.CellFormat(40, 10, "Unit Price", "1", 0, "", true, 0, "")
	pdf.CellFormat(40, 10, "Cost", "TBL", 0, "", true, 0, "")
	pdf.Ln(-1)

	// table body
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0)
	totalCost := 0.0

	for i, emp := range employees {
		cost := emp.TotalHours * emp.BillableRate
		totalCost += cost

		pdf.CellFormat(40, 7, ""+strconv.Itoa(i), "TBR", 0, "R", false, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("%.2f", emp.TotalHours), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("%.2f", emp.BillableRate), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("%.2f", cost), "TBL", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}

	// totals
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 7, "", "TBR", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, "", "1", 0, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 7, "Total", "1", 0, "", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 7, fmt.Sprintf("%.2f", totalCost), "TBL", 0, "R", false, 0, "")
	pdf.Ln(-1)

	filename := cName + "_invoice.pdf"
	return filename, pdf.OutputFileAndClose(filename)
}
