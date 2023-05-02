package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"main/utils"
	"net/http"
	"strconv"
	"strings"
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
	// r.POST("/uploadJSON", HandleJSONUpload)
	r.Run(":5000")

	http.HandleFunc("/uploadCSV", HandleCSV)
	http.ListenAndServe(":5000", nil)
}

func HandleCSVUpload(c *gin.Context) {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open the file
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

	var payslips []utils.PaySlip

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

	// Send the response
	c.JSON(http.StatusOK, payload)
}

// Handle CSV file uploaded from POST request
func HandleCSV(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("csv_file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Process the CSV data
	csvData, err := csv.NewReader(file).ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a response to confirm that the CSV file has been successfully processed
	response := make(map[string]interface{})
	response["message"] = "CSV file successfully processed"
	response["data"] = csvData

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

// Handle request to calculate tax
func HandleRequest(c *gin.Context) {
	var employees []utils.Employee
	if err := c.ShouldBindJSON(&employees); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var payslips []utils.PayslipResponse
	for _, employee := range employees {
		payslip := GenerateRESTPayslip(employee)
		payslips = append(payslips, payslip)
	}

	payload := gin.H{"message": "Calculated tax successfully", "payslips": payslips}

	c.JSON(http.StatusOK, payload)
}

// Generate payslip for given employee
func GenerateRESTPayslip(employee utils.Employee) utils.PayslipResponse {
	var payslip utils.PayslipResponse
	payslip.Employee = employee
	payslip.GrossIncome = int(employee.AnnualSalary)
	incomeTax := utils.CalculateTax(employee.AnnualSalary)
	payslip.IncomeTax = int(incomeTax)
	payslip.NetIncome = int(employee.AnnualSalary - incomeTax)
	payslip.Superannuation = int(utils.CalculateSuper(employee.SuperRate, employee.AnnualSalary))
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
func ParseEmployeeCSV(line []string) (utils.CSVEmployee, error) {
	annualSalary, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return utils.CSVEmployee{}, err
	}

	superRate, err := strconv.ParseFloat(strings.TrimRight(line[3], "%"), 64)
	if err != nil {
		return utils.CSVEmployee{}, err
	}

	return utils.CSVEmployee{FirstName: line[0], LastName: line[1], AnnualSalary: annualSalary, PaymentStart: line[4], SuperRate: superRate / 100}, nil
}

// Generate Payslip for given employee
func GenerateCSVPayslip(employee utils.CSVEmployee) utils.PaySlip {
	var payslip utils.PaySlip
	payslip.Name = employee.FirstName + " " + employee.LastName
	payslip.AnnualSalary = employee.AnnualSalary
	payslip.IncomeTax = utils.CalculateTax(employee.AnnualSalary)
	payslip.NetIncome = employee.AnnualSalary - payslip.IncomeTax
	payslip.PayPeriod = employee.PaymentStart
	payslip.Superannuation = utils.CalculateSuper(employee.SuperRate, employee.AnnualSalary)
	return payslip
}
