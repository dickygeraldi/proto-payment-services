package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"proto-parking-services/pkg/services/api/v1/global"
	"strconv"
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
func ParkingRegistration(platNo string, timeRequest time.Time, connection *sql.DB, ctx context.Context) (code, message, statusMessage, invoice, waktu, status string) {
	var dataParking int

	// Checking platNo
	checkPlatNo := global.GenerateQueryParking(map[string]string{
		"platNo": platNo,
	})
	fmt.Println(checkPlatNo)

	rows := connection.QueryRowContext(ctx, checkPlatNo)
	err := rows.Scan(&dataParking)
	if err != nil {
		fmt.Println(err)
	}

	if dataParking >= 1 {
		code = "04"
		statusMessage = "Data sudah ada di database"
		message = platNo + " sudah terdaftar parkir pada " + timeRequest.Format("2006-01-02 15:04")
	} else {
		invoice = getRandomString()
		message = "Transaksi berhasil diproses"
		code = "00"
		status = "PENDING"
		statusMessage = "Data berhasil diproses"
		waktu = timeRequest.Format("2006-01-02 15:04")

		go func() {
			sql := `INSERT INTO "dataParking" ("invoiceId", "merchantId", "platNo", "enteredDate", "status") VALUES ($1, $2, $3, $4, $5)`

			_, err := connection.Query(sql, invoice, "MerchantId", platNo, waktu, status)

			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	return code, message, statusMessage, invoice, waktu, status
}

// Function for parking validation
func ValidationParking(platNo string, timeRequest time.Time, connection *sql.DB, ctx context.Context) (code, message, status, qrContent string) {
	var invoiceId string
	var timeDiff string

	// Checking get invoice and enteredDate
	checkInvoice := global.GenerateQueryParkingData(map[string]string{
		"platNo":   platNo,
		"dateTime": timeRequest.Format("2006-01-02 15:04"),
	})

	fmt.Println(checkInvoice)

	rows := connection.QueryRowContext(ctx, checkInvoice)
	err := rows.Scan(&invoiceId, &timeDiff)
	fmt.Println(timeDiff)
	if err != nil {
		fmt.Println(err)
		code = "05"
		message = "Plat Nomor tidak ditemukan"
		status = "Gagal"
		qrContent = ""
	} else {
		url := os.Getenv("URL_QREN")
		timeTransaction, _ := strconv.Atoi(timeDiff)

		nominalTransaction := timeTransaction * 2000
		transaksi := strconv.Itoa(nominalTransaction)

		body := &global.Qren{
			MerchantApiKey: os.Getenv("API_KEY"),
			InvoiceName:    invoiceId,
			Nominal:        transaksi,
			StaticQR:       "0",
			QrGaruda:       "1",
		}

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		req, _ := http.NewRequest("POST", url, buf)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "basic "+os.Getenv("AUTH_QREN"))
		fmt.Println(buf)

		client := &http.Client{}
		res, e := client.Do(req)
		if e != nil {
			fmt.Println(e)
		}

		defer res.Body.Close()

		if res.StatusCode == 200 {
			body, _ := ioutil.ReadAll(res.Body)

			code = "00"
			message = "Generate QR Content berhasil"
			status = "Transaksi Berhasil"
			qrContent = string(body)
		} else {
			code = "10"
			message = "Error sistem pada QREN"
			status = "Transaksi gagal"
			qrContent = ""
		}
	}

	return code, message, status, qrContent
}
