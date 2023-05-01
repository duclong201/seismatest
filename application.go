package main

import (
	"fmt"
	"main/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// func main() {
// 	csvFile, err := os.Open("employee.csv")
// 	if err != nil {
// 		fmt.Println("Failed to read csv file.", err)
// 		return
// 	}
// 	fmt.Println("Successfully Opened CSV file")
// 	defer csvFile.Close()

// 	csvLines, err := csv.NewReader(csvFile).ReadAll()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	for i, line := range csvLines {
// 		if i == 0 {
// 			continue
// 		}
// 		newEmployee, err := ParseEmployeeJSON(line)
// 		payslip := GeneratePayslipCSV(newEmployee)
// 		fmt.Println(payslip)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 	}
// }

// // ParseEmployee parses a string array of employee details and returns an Employee object
// func ParseEmployeeJSON(record []string) (utils.CSVEmployee, error) {
// 	annualSalary, err := strconv.ParseFloat(record[2], 64)
// 	if err != nil {
// 		return utils.CSVEmployee{}, err
// 	}
// 	superRate, err := strconv.ParseFloat(strings.TrimRight(record[3], "%"), 64)
// 	if err != nil {
// 		return utils.CSVEmployee{}, err
// 	}

// 	if err != nil {
// 		return utils.CSVEmployee{}, err
// 	}
// 	return utils.CSVEmployee{
// 		FirstName:    record[0],
// 		LastName:     record[1],
// 		AnnualSalary: annualSalary,
// 		SuperRate:    superRate,
// 		PaymentStart: record[4],
// 	}, nil
// }

// // GeneratePayslip method returns the payslip for given employee
// func GeneratePayslipCSV(employee utils.CSVEmployee) utils.PaySlip {
// 	var ps utils.PaySlip
// 	ps.Name = employee.FirstName + " " + employee.LastName
// 	ps.AnnualSalary = employee.AnnualSalary
// 	ps.IncomeTax = utils.CalculateTax(employee.AnnualSalary)
// 	ps.NetIncome = employee.AnnualSalary - ps.IncomeTax
// 	ps.PayPeriod = employee.PaymentStart
// 	ps.Superannuation = utils.CalculateSuper(employee.SuperRate, ps.AnnualSalary)
// 	return ps
// }

func main() {
	fmt.Println("Start Application")

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Tested!",
		})
	})
	r.POST("/calculateTax", HandleRequest)
	r.Run(":5000")

	fmt.Println("Handle Request with gin")
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
		payslip := GeneratePayslip(employee)
		payslips = append(payslips, payslip)
	}

	payload := gin.H{"message": "Calculated tax successfully", "payslips": payslips}

	c.JSON(http.StatusOK, payload)
}

// Generate payslip for given employee
func GeneratePayslip(employee utils.Employee) utils.PayslipResponse {
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
