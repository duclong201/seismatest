# Seisma Test

## Intro

This Go application is deployed to Elastic Beanstalk and is hosted on GitHub. Continuous integration and deployment (CI/CD) is being handled by GitHub actions, which automatically deploys the latest version of the application whenever changes are made to the main branch. Users can interact with the application by sending POST requests to Elastic Beanstalk, which responds with JSON data.

## Assumptions
- The payment period is always per calendar month.
- Since the request payload only has details about paymentMonth, without specifying which month. The from and to date returned by the response payload will be the first and last day of the current month.

## How to use
- Use Postman to send POST Request to AWS Elastic Beanstalk.
- There are 2 endpoints (base URL: `http://seismatest-env.eba-3vssjefi.ap-southeast-2.elasticbeanstalk.com`)
 1. `/calculateTax`: send POST request with JSON body. E.g:
 ```
 [{
    "firstName": "Long",
    "lastName" : "Nguyen",
    "annualSalary": 150000,
    "paymentMonth": 1,
    "superRate": 0.09
}, {
    "firstName": "Duc",
    "lastName" : "Tran",
    "annualSalary": 180000,
    "paymentMonth": 1,
    "superRate": 0.1
}]
 ```
 
 2. `/upload`: send POST request with file attached. In Postman request body, select `form-data` and attach CSV or JSON text file.
 
 ```
 {
    "file": "employee.csv"
 }
 ```
- Depends on the file types, the API will return calculated tax and other informations.

- For CSV: 
```
{
    "message": "Calculated tax successfully",
    "payslips": [
        {
            "Name": "Monica Tan",
            "AnnualSalary": 60050,
            "PayPeriod": "01 March – 31 March",
            "IncomeTax": 11063,
            "NetIncome": 48987,
            "Superannuation": 5405
        },
        {
            "Name": "Brend Tulu",
            "AnnualSalary": 120000,
            "PayPeriod": "01 March – 31 March",
            "IncomeTax": 32032,
            "NetIncome": 87968,
            "Superannuation": 12000
        },
        {
            "Name": "Long Nguyen",
            "AnnualSalary": 100000,
            "PayPeriod": "01 March - 31 March",
            "IncomeTax": 24632,
            "NetIncome": 75368,
            "Superannuation": 10000
        }
    ]
}
```

- For JSON text file:
```
{
    "message": "Calculated tax successfully",
    "payslips": [
        {
            "Employee": {
                "FirstName": "Long",
                "LastName": "Nguyen",
                "AnnualSalary": 150000,
                "PaymentMonth": 1,
                "SuperRate": 0.09
            },
            "FromDate": "01 May",
            "ToDate": "31 May",
            "GrossIncome": 150000,
            "IncomeTax": 43132,
            "Superannuation": 13500,
            "NetIncome": 106868
        },
        {
            "Employee": {
                "FirstName": "Duc",
                "LastName": "Tran",
                "AnnualSalary": 180000,
                "PaymentMonth": 1,
                "SuperRate": 0.1
            },
            "FromDate": "01 May",
            "ToDate": "31 May",
            "GrossIncome": 180000,
            "IncomeTax": 54232,
            "Superannuation": 18000,
            "NetIncome": 125768
        }
    ]
}
```
