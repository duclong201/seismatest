package utils

type CSVEmployee struct {
	FirstName    string
	LastName     string
	AnnualSalary float64
	PaymentStart string
	SuperRate    float64
}

type Employee struct {
	FirstName    string
	LastName     string
	AnnualSalary float64
	PaymentMonth int
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

type TaxRate struct {
	MaxValue float64
	Rate     float64
	Bracket  float64
	FixValue float64
}

type Calculator struct {
	TaxRates []TaxRate
}

type PayslipResponse struct {
	Employee
	FromDate       string
	ToDate         string
	GrossIncome    int
	IncomeTax      int
	Superannuation int
	NetIncome      int
}
