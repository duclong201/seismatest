package api

import (
	"encoding/json"
	"fmt"
	"main/utils"
	"net/http"
	"time"
)

type PayslipResponse struct {
	utils.Employee
	FromDate       string
	ToDate         string
	GrossIncome    int
	IncomeTax      int
	Superannuation int
	NetIncome      int
}

func main() {
	http.HandleFunc("/calculateTax", HandleRequest)
	http.ListenAndServe(":8080", nil)
}

// Handle REST request to calculate tax
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	var employees []utils.Employee
	err := json.NewDecoder(r.Body).Decode(&employees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var payslips []PayslipResponse
	for _, employee := range employees {
		payslip := GenerateJSONResponse(employee)
		payslips = append(payslips, payslip)
	}
	w.Header().Set("ContentType", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payslips)
}

// Generate payslip for given employee
func GenerateJSONResponse(employee utils.Employee) PayslipResponse {
	var payslip PayslipResponse
	payslip.Employee = employee
	payslip.AnnualSalary = employee.AnnualSalary
	incomeTax := utils.CalculateTax(employee.AnnualSalary)
	payslip.IncomeTax = int(incomeTax)
	payslip.NetIncome = int(employee.AnnualSalary - incomeTax)
	payslip.Superannuation = int(utils.CalculateSuper(employee.SuperRate, employee.AnnualSalary))
	currentMonth := time.Now().Month().String()
	payslip.FromDate = "01 " + currentMonth
	payslip.ToDate = lastDayOfCurrentMonth() + currentMonth
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
