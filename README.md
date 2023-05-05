# Seisma Test

## Assumptions
- The payment period is always per calendar month.
- Since the request payload only has details about paymentMonth, without specifying which month. The from and to date returned by the response payload will be the first and last day of the current month.

## How to use
- Using Postman to send a POST request to AWS Elastic Beanstalk endpoint "http://seismatest-env.eba-3vssjefi.ap-southeast-2.elasticbeanstalk.com/upload"
- In Request Body, select `form-data` and attach CSV or JSON txt file.
- Depends on the file types, the API will return calculated tax and informations

