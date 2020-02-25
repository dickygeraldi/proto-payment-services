package validation

import (
	"proto-parking-services/pkg/services/api/v1/logging"
	"regexp"
)

var re = regexp.MustCompile("select|insert|update|alter|delete")
var digitCheck = regexp.MustCompile(`^[0-9]+$`)
var checkUsername = regexp.MustCompile("^[a-z0-9]+(?:_[a-z0-9]+)*$")

// Function for validation user services request
func ParkingRegistration(api, platNo, dataIp, dateTime string) (string, bool) {
	if api != "" && platNo != "" {
		if re.MatchString(api) == true || re.MatchString(platNo) == true {
			go logging.SetLogging("Parking Request", dataIp, "SQL Injection", "Warning message", "warning", "SQL Injection in this request", dateTime)
			return "Coba lagi nanti, transaksi di pending", false
		} else {
			return "", true
		}
	} else {
		return "Semua data harus terisi", false
	}
}

// invoiceId
// status
// message
// trxId
// amount
