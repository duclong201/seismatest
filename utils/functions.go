package utils

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

func GetTaxRates() []TaxRate {
	// Provide default tax rates
	defaultTaxRates := []TaxRate{{0, 18200, 0, 0},
		{18200, 37000, 0, 0.19},
		{37000, 87000, 3572, 0.325},
		{87000, 180000, 19822, 0.37},
		{180000, math.MaxFloat64, 54232, 0.45}}

	// Read tax rates in CSV file
	csvFile, err := os.Open("taxRates.csv")
	if err != nil {
		fmt.Println("Failed to read csv file. Using default calculator", err)
		return defaultTaxRates
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	var taxRates []TaxRate

	for i, line := range csvLines {
		if i == 0 {
			continue
		}
		taxRate, err := ParseTaxRates(line)
		if err != nil {
			fmt.Println(err)
			return defaultTaxRates
		}
		taxRates = append(taxRates, taxRate)
	}
	return taxRates
}

// Parse Tax rate from read csv line
func ParseTaxRates(line []string) (TaxRate, error) {
	defaultTaxRate := TaxRate{Bracket: 0, MaxValue: 0, FixValue: 0, Rate: 0}
	bracket, err := strconv.ParseFloat(line[0], 64)
	if err != nil {
		return defaultTaxRate, err
	}
	maxValue, err := strconv.ParseFloat(line[1], 64)
	if err != nil {
		if line[1] == "math.MaxFloat64" {
			maxValue = math.MaxFloat64
		} else {
			return defaultTaxRate, err
		}
	}
	fixValue, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return defaultTaxRate, err
	}
	rate, err := strconv.ParseFloat(line[3], 64)
	if err != nil {
		return defaultTaxRate, err
	}
	return TaxRate{Bracket: bracket, MaxValue: maxValue, FixValue: fixValue, Rate: rate}, nil
}

// CalculateTax method calculates the income tax for the given annual salary
func CalculateTax(annualSalary float64) float64 {
	var tax float64
	for i, tr := range GetTaxRates() {
		if i == 0 {
			continue
		}
		if annualSalary < tr.MaxValue {
			tax = (annualSalary-tr.Bracket)*tr.Rate + tr.FixValue
			break
		}
	}
	return math.Round(tax)
}

// CalculateSuper calculates the superannuation for the given super rate and gross income
func CalculateSuper(superRate float64, grossIncome float64) float64 {
	return math.Round(superRate * grossIncome / 100)
}
