# Seisma Test

## Assumptions
- The payment period is always per calendar month.
- Since the request payload only has details about paymentMonth, without specifying which month. The from and to date returned by the response payload will be the first and last day of the current month.

## How to use
- Using Postman to send a POST request to AWS Elastic Beanstalk endpoint "http://seismatest-env.eba-3vssjefi.ap-southeast-2.elasticbeanstalk.com/upload"
- In Request Body, select `form-data` and attach CSV or JSON txt file.
- Depends on the file types, the API will return calculated tax and informations

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
