package main

import (
	"encoding/csv"
	"fmt"
	"main/utils"
	"os"
	"strconv"
)

func main() {
	csvFile, err := os.Open("../csv/employee.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	var payslips []utils.PaySlip

	for i, line := range csvLines {
		if i == 0 {
			continue
		}
		employee, err := ParseEmployeeCSV(line)
		if err != nil {
			fmt.Println(err)
			return
		}
		payslip := GeneratePayslip(employee)
		fmt.Println(payslip)
		payslips = append(payslips, payslip)
	}
}

// Parse employee for given line from the csv file
func ParseEmployeeCSV(line []string) (utils.CSVEmployee, error) {
	annualSalary, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return utils.CSVEmployee{}, err
	}

	superRate, err := strconv.ParseFloat(line[3], 64)
	if err != nil {
		return utils.CSVEmployee{}, err
	}

	return utils.CSVEmployee{FirstName: line[0], LastName: line[1], AnnualSalary: annualSalary, PaymentStart: line[4], SuperRate: superRate}, nil
}

// Generate Payslip for given employee
func GeneratePayslip(employee utils.CSVEmployee) utils.PaySlip {
	var payslip utils.PaySlip
	payslip.Name = employee.FirstName + " " + employee.LastName
	payslip.AnnualSalary = employee.AnnualSalary
	payslip.IncomeTax = utils.CalculateTax(employee.AnnualSalary)
	payslip.NetIncome = employee.AnnualSalary - payslip.IncomeTax
	payslip.PayPeriod = employee.PaymentStart
	payslip.Superannuation = utils.CalculateSuper(employee.SuperRate, employee.AnnualSalary)
	return payslip
}
