package models

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"proto-parking-services/pkg/services/api/v1/global"
	"time"
)

// Set global environment variable
var conf *global.Configuration
var messageError map[int]global.MessageError
var level, cases, fatal string

// Function initialization
func init() {
	conf = global.New()
	messageError = global.GetMessageError()
}

// Function Generate random number
func getRandomString() string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	const charset = "1234567890" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

// Function for register account
func ParkingRegistration(platNo, timeRequest string, connection *sql.DB, ctx context.Context) (code, message, statusMessage, invoice, waktu, status string) {
	var dataParking int

	// Checking platNo
	checkPlatNo := global.GenerateQueryParking(map[string]string{
		"platNo": platNo,
	})

	rows := connection.QueryRowContext(ctx, checkPlatNo)
	err := rows.Scan(&dataParking)
	if err != nil {
		fmt.Println(err)
	}

	if dataParking >= 1 {
		code = "04"
		statusMessage = "Data sudah ada di database"
		message = platNo + " sudah melakukan parkir pada " + timeRequest
	} else {
		invoice = getRandomString()
		message = "Transaksi berhasil diproses"
		code = "00"
		status = "PENDING"
		statusMessage = "Data berhasil diproses"
		waktu = timeRequest

		go func() {
			sql := `INSERT INTO "dataParking" ("invoiceId", "merchantId", "platNo", "enteredDate") VALUES ($1, $2, $3, $4)`

			_, err := connection.Query(sql, invoice, "MerchantId", platNo, waktu)

			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	return code, message, statusMessage, invoice, waktu, status
}

// Function for parking validation
// func ValidationParking(platNo, )
