package cashbill

import (
	"net/http"
	"fmt"
	"encoding/json"
	"time"
	"strconv"
	"regexp"
)

type cashbillResponse struct {
	Active           bool   `json:"active"`
	Number           string `json:"number"`
	ClientNumber     string `json:"clientNumber"`
	ActiveFrom       int64  `json:"activeFrom"`
	CodeValidityTime int64  `json:"codeValidityTime"`
	TimeRemaining    int64  `json:"timeRemaining"`
}

type CodeStatus struct {
	Code             string
	Active           bool
	Number           string
	ClientNumber     string
	ActiveFrom       time.Time
	CodeValidityTime time.Duration
	TimeRemaining    time.Duration
	Value            float32
}

func CheckSMSCode(code string, token string) (status *CodeStatus, err error) {
	//Remove illegal characters from SMS code
	if code, err = FormatCode(code); err != nil {
		return
	}

	//Call Cashbill API
	var resp *http.Response
	if resp, err = http.Get(fmt.Sprintf("https://sms.cashbill.pl/code/%s/%s/", token, code)); err != nil {
		return
	}

	defer resp.Body.Close()

	//Checks if status code is 200, if not return an error
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Status code: %v", resp.StatusCode)
		return
	}

	var cashbillResp cashbillResponse
	if err = json.NewDecoder(resp.Body).Decode(cashbillResp); err != nil {
		return
	}

	var value float32
	var intValue int

	//Get value based on number to which the SMS was sent
	if cashbillResp.Number[0:1] == "7" {
		if intValue, err = strconv.Atoi(cashbillResp.Number[1:2]); err != nil {
			return
		}

		if value = float32(intValue) / 2; value == 0 {
			value = 0.25
		}
	} else {
		if intValue, err = strconv.Atoi(cashbillResp.Number[1:3]); err != nil {
			return
		}

		value = float32(intValue) / 2
	}

	//Prepare the response
	status = &CodeStatus{
		code,
		cashbillResp.Active,
		cashbillResp.Number,
		cashbillResp.ClientNumber,
		time.Unix(cashbillResp.ActiveFrom, 0),
		time.Duration(cashbillResp.CodeValidityTime) * time.Second,
		time.Duration(cashbillResp.TimeRemaining) * time.Second,
		value,
	}
	return
}

//FormatCode function replaces illegal characters in SMS code
func FormatCode(code string) (string, error) {
	re, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		return "", err
	}

	return re.ReplaceAllString(code, ""), nil
}

