package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Tested!",
		})
	})
	r.POST("/calculateTax", HandleRequest)
	r.POST("/uploadCSV", HandleCSVUpload)
	r.POST("/uploadJSON", HandleJSONUpload)
	r.Run(":5000")
}

// HandleCSVUpload method handles the csv file uploaded with POST request and return processed payslips
func HandleCSVUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentType := file.Header.Get("Content-Type")
	fmt.Println(contentType)

	csvFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer csvFile.Close()
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var payslips []PaySlip

	for i, line := range csvLines {
		if i == 0 {
			continue
		}
		employee, err := ParseEmployeeCSV(line)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		payslip := GenerateCSVPayslip(employee)
		payslips = append(payslips, payslip)
	}
	payload := gin.H{"message": "Calculated tax successfully", "payslips": payslips}
	c.JSON(http.StatusOK, payload)
}

// Handle request to calculate tax
func HandleRequest(c *gin.Context) {
	var employees []Employee
	if err := c.ShouldBindJSON(&employees); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var payslips []PayslipResponse
	for _, employee := range employees {
		payslip := GenerateRESTPayslip(employee)
		payslips = append(payslips, payslip)
	}

	payload := gin.H{"message": "Calculated tax successfully", "payslips": payslips}

	c.JSON(http.StatusOK, payload)
}

// Handle JSON Upload from POST request and return processed payslips
func HandleJSONUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentType := file.Header.Get("Content-Type")
	fmt.Println(contentType)

	jsonFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to open file",
		})
		return
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid JSON data",
		})
		return
	}

	var jsonData []map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid JSON data",
		})
		return
	}

	var payslips []PayslipResponse

	for _, obj := range jsonData {
		fmt.Println(obj)
		employee := Employee{FirstName: obj["firstName"].(string),
			LastName:     obj["lastName"].(string),
			AnnualSalary: obj["annualSalary"].(float64),
			PaymentMonth: int(obj["paymentMonth"].(float64)),
			SuperRate:    obj["superRate"].(float64)}
		payslip := GenerateRESTPayslip(employee)
		payslips = append(payslips, payslip)
	}

	payload := gin.H{"message": "Calculated tax successfully", "payslips": payslips}

	c.JSON(http.StatusOK, payload)
}

// Generate payslip for given employee
func GenerateRESTPayslip(employee Employee) PayslipResponse {
	var payslip PayslipResponse
	payslip.Employee = employee
	payslip.GrossIncome = int(employee.AnnualSalary)
	incomeTax := CalculateTax(employee.AnnualSalary)
	payslip.IncomeTax = int(incomeTax)
	payslip.NetIncome = int(employee.AnnualSalary - incomeTax)
	payslip.Superannuation = int(CalculateSuper(employee.SuperRate, employee.AnnualSalary))
	currentMonth := time.Now().Month().String()
	payslip.FromDate = "01 " + currentMonth
	payslip.ToDate = lastDayOfCurrentMonth() + " " + currentMonth
	return payslip
}

// Return the number of days of current month
func lastDayOfCurrentMonth() string {
	now := time.Now()
	firstDayOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	// Subtract one day from it to get the last day of the given month
	lastDay := fmt.Sprintf("%d", firstDayOfNextMonth.AddDate(0, 0, -1).Day())
	return lastDay
}

// Parse employee for given line from the csv file
func ParseEmployeeCSV(line []string) (CSVEmployee, error) {
	annualSalary, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return CSVEmployee{}, err
	}

	superRate, err := strconv.ParseFloat(line[3], 64)
	if err != nil {
		return CSVEmployee{}, err
	}

	return CSVEmployee{FirstName: line[0], LastName: line[1], AnnualSalary: annualSalary, PaymentStart: line[4], SuperRate: superRate}, nil
}

// Generate Payslip for given employee
func GenerateCSVPayslip(employee CSVEmployee) PaySlip {
	var payslip PaySlip
	payslip.Name = employee.FirstName + " " + employee.LastName
	payslip.AnnualSalary = employee.AnnualSalary
	payslip.IncomeTax = CalculateTax(employee.AnnualSalary)
	payslip.NetIncome = employee.AnnualSalary - payslip.IncomeTax
	payslip.PayPeriod = employee.PaymentStart
	payslip.Superannuation = CalculateSuper(employee.SuperRate, employee.AnnualSalary)
	return payslip
}

// Load TaxRates from csv file
func GetTaxRates() []TaxRate {
	// Provide default tax rates
	defaultTaxRates := []TaxRate{{0, 18200, 0, 0},
		{18200, 37000, 0, 0.19},
		{37000, 87000, 3572, 0.325},
		{87000, 180000, 19822, 0.37},
		{180000, math.MaxFloat64, 54232, 0.45}}

	// Read tax rates in CSV file
	csvFile, err := os.Open("../csv/taxRates.csv")
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
	return math.Round(superRate * grossIncome)
}
