package validation

import (
	"proto-parking-services/pkg/services/api/v1/logging"
	"regexp"
	"time"
)

var re = regexp.MustCompile("select|insert|update|alter|delete")

// Function for validation user services request
func ParkingRegistration(api, platNo, dataIp string, dateTime time.Time) (string, bool) {
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
