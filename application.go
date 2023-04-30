package main

import (
	"encoding/csv"
	"fmt"
	"main/utils"
	"os"
	"strconv"
	"strings"
)

func main() {
	csvFile, err := os.Open("employee.csv")
	if err != nil {
		fmt.Println("Failed to read csv file.", err)
		return
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for i, line := range csvLines {
		if i == 0 {
			continue
		}
		newEmployee, err := ParseEmployeeJSON(line)
		payslip := GeneratePayslipCSV(newEmployee)
		fmt.Println(payslip)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// ParseEmployee parses a string array of employee details and returns an Employee object
func ParseEmployeeJSON(record []string) (utils.CSVEmployee, error) {
	annualSalary, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return utils.CSVEmployee{}, err
	}
	superRate, err := strconv.ParseFloat(strings.TrimRight(record[3], "%"), 64)
	if err != nil {
		return utils.CSVEmployee{}, err
	}

	if err != nil {
		return utils.CSVEmployee{}, err
	}
	return utils.CSVEmployee{
		FirstName:    record[0],
		LastName:     record[1],
		AnnualSalary: annualSalary,
		SuperRate:    superRate,
		PaymentStart: record[4],
	}, nil
}

// GeneratePayslip method returns the payslip for given employee
func GeneratePayslipCSV(employee utils.CSVEmployee) utils.PaySlip {
	var ps utils.PaySlip
	ps.Name = employee.FirstName + " " + employee.LastName
	ps.AnnualSalary = employee.AnnualSalary
	ps.IncomeTax = utils.CalculateTax(employee.AnnualSalary)
	ps.NetIncome = employee.AnnualSalary - ps.IncomeTax
	ps.PayPeriod = employee.PaymentStart
	ps.Superannuation = utils.CalculateSuper(employee.SuperRate, ps.AnnualSalary)
	return ps
}
