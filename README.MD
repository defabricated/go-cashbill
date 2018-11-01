# go-cashbill

[![GoDoc](https://godoc.org/github.com/defabricated/go-cashbill?status.svg)](http://godoc.org/github.com/defabricated/go-cashbill)

Go API for Cashbill.pl which can be used to get SMS code data

NOT TESTED

## Example

```go
code := "123456"
token := "CASHBILL_TOKEN"

resp, err := cashbill.CheckSMSCode(code, token)
if err != nil {
    panic(err)
    return
}

fmt.Printf("The \"%s\" code was sent to %s and it is worth %d PLN.\n", resp.Code, resp.Number, resp.Value)
```

## Response fields

```go
resp.Code //SMS code without illegal characters

resp.Active //Is code active or not

resp.Number //Phone number to which the SMS was sent

resp.ClientNumber //Sender's phone number

resp.ActiveFrom //Time when the code was activated

resp.CodeValidityTime //How long the SMS code stays active

resp.TimeRemaining //For how long the SMS code will stay active

resp.Value //How much is the SMS code worth. To calculate how much sender had to pay multiply the value by 2.46
```
