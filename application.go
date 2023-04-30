package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Employee struct {
	FirstName    string
	LastName     string
	AnnualSalary float64
	PaymentStart string
	SuperRate    float64
}

type PaySlip struct {
	Name           string
	AnnualSalary   float64
	PayPeriod      string
	IncomeTax      float64
	NetIncome      float64
	Superannuation float64
}

type TaxRates struct {
	MaxValue float64
	Rate     float64
	Bracket  float64
	FixValue float64
}

type IncomeTaxCalculator struct {
	TaxRates []TaxRates
}

// CalculateTax method calculates the income tax for the given annual salary
func (itc IncomeTaxCalculator) CalculateTax(annualSalary float64) float64 {
	var tax float64
	for i, tr := range itc.TaxRates {
		if i == 0 {
			continue
		}
		if annualSalary > tr.MaxValue {
			continue
		} else {
			tax = (annualSalary-tr.Bracket)*tr.Rate + tr.FixValue
			break
		}
	}
	return math.Round(tax/12.0 + 0.5)
}

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
		newEmployee, err := ParseEmployee(line)
		payslip := GeneratePayslip(newEmployee)
		fmt.Println(payslip)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

// CalculateSuper calculates the superannuation for the given super rate and gross income
func CalculateSuper(superRate float64, grossIncome float64) float64 {
	return math.Round(superRate * grossIncome / 100)
}

// ParseEmployee parses a string array of employee details and returns an Employee object
func ParseEmployee(record []string) (Employee, error) {
	annualSalary, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return Employee{}, err
	}
	superRate, err := strconv.ParseFloat(strings.TrimRight(record[3], "%"), 64)
	if err != nil {
		return Employee{}, err
	}

	if err != nil {
		return Employee{}, err
	}
	return Employee{
		FirstName:    record[0],
		LastName:     record[1],
		AnnualSalary: annualSalary,
		SuperRate:    superRate,
		PaymentStart: record[4],
	}, nil
}

// GeneratePayslip method returns the payslip for given employee
func GeneratePayslip(employee Employee) PaySlip {
	var ps PaySlip
	ps.Name = employee.FirstName + " " + employee.LastName
	ps.AnnualSalary = employee.AnnualSalary
	ps.IncomeTax = IncomeTaxCalculator{TaxRates: []TaxRates{
		{0, 0, 0, 0},
		{18200, 0, 0, 0},
		{37000, 0.19, 18200, 0},
		{87000, 0.325, 37000, 3572},
		{180000, 0.37, 87000, 19822},
		{math.MaxFloat64, 0.45, 180000, 54232},
	}}.CalculateTax(employee.AnnualSalary)
	ps.NetIncome = employee.AnnualSalary - ps.IncomeTax
	ps.PayPeriod = employee.PaymentStart
	ps.Superannuation = CalculateSuper(employee.SuperRate, ps.AnnualSalary)
	return ps
}
