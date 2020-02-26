package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"proto-parking-services/pkg/services/api/v1/global"
	"runtime"
	"strconv"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

// Set global environment variable
var conf *global.Configuration
var level, cases, fatal string
var connection *gosocketio.Client

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func HandleSocket() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c, err := gosocketio.Dial(
		gosocketio.GetUrl(os.Getenv("SOCKET_HOST"), 80, false),
		transport.GetDefaultWebsocketTransport())

	if err != nil {
		log.Fatal("Error 1: ", err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal("Error 2 :", err)
	}

	connection = c

	defer c.Close()
}

func getDataFromChannel(channel string) bool {

	fmt.Println("Listening to channel: ", channel)

	data := false
	counting := 0

	go func() {
		for {
			err := connection.On(channel, func(h *gosocketio.Channel, args interface{}) {
				mResult := args.(map[string]interface{})
				if mResult["invoice"] != "" {
					log.Println("Data Invoice ", mResult["invoice"])
					log.Println("Data merchantApi : ", mResult["merchantApiKey"])
					log.Println("Data Status : ", mResult["status"])
					data = true
				}
			})

			if err != nil {
				log.Fatal(err)
			}

			if data == true && counting <= 3 {
				break
			} else if data == false && counting <= 3 {
				time.Sleep(15 * time.Second)
				counting = counting + 1
			} else if counting > 3 {
				break
			}
		}
	}()

	return true
}

// Function initialization
func init() {
	conf = global.New()
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
	location, _ := time.LoadLocation("Asia/Jakarta")

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
		message = platNo + " sudah terdaftar parkir pada " + timeRequest.In(location).Format("2006-01-02 15:04")
	} else {
		invoice = getRandomString()
		message = "Transaksi berhasil diproses"
		code = "00"
		status = "PENDING"
		statusMessage = "Data berhasil diproses"
		waktu = timeRequest.In(location).Format("2006-01-02 15:04")

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
func ValidationParking(platNo string, timeRequest time.Time, connection *sql.DB, ctx context.Context) (code, message, status, qrContent, jamMasuk, jamKeluar, totalJam, amount string) {
	var invoiceId string
	var timeDiff string
	var enteredDate time.Time
	location, _ := time.LoadLocation("Asia/Jakarta")

	// Checking get invoice and enteredDate
	checkInvoice := global.GenerateQueryParkingData(map[string]string{
		"platNo":   platNo,
		"dateTime": timeRequest.In(location).Format("2006-01-02 15:04"),
	})

	rows := connection.QueryRowContext(ctx, checkInvoice)
	err := rows.Scan(&invoiceId, &timeDiff, &enteredDate)

	if err != nil {
		fmt.Println(err)
		code = "05"
		message = "Plat Nomor tidak ditemukan"
		status = "Gagal"
		qrContent = ""
		jamKeluar = ""
		jamMasuk = ""
		totalJam = ""
		amount = ""
	} else {
		var nominalTransaction int
		url := os.Getenv("URL_QREN")
		timeTransaction, _ := strconv.Atoi(timeDiff)

		if timeTransaction >= 2 {
			nominalTransaction = (timeTransaction-1)*1000 + 2000
		} else {
			nominalTransaction = 2000
		}

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
		req.Header.Add("Authorization", "Basic "+os.Getenv("AUTH_QREN"))

		client := &http.Client{}
		res, e := client.Do(req)
		if e != nil {
			fmt.Println(e)
		}

		defer res.Body.Close()

		if res.StatusCode == 200 {
			c := make(map[string]interface{})

			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal([]byte(string(body)), &c)

			code = "00"
			message = "Generate QR Content berhasil"
			status = "Transaksi Berhasil"
			qrContent = fmt.Sprintf("%v", c["content"])
			jamKeluar = timeRequest.In(location).Format("2006-01-02 15:04")
			jamMasuk = enteredDate.Format("2006-01-02 15:04")
			amount = transaksi

			if timeDiff == "0" {
				timeDiff = "1"
			}

			totalJam = timeDiff

			go func() {

				sql := fmt.Sprintf(`UPDATE "dataParking" set "qreninvoiceid" = $1, "amount" = $2, "exitDate" = $3 where "invoiceId" = $4`)

				_, err := connection.Query(sql, c["invoiceId"], transaksi, jamKeluar, invoiceId)

				if err != nil {
					fmt.Println(err)
				}
			}()

			go func() {
				x := getDataFromChannel(fmt.Sprintf("%v", c["invoiceId"]))
				if x == true {
					fmt.Println("Transaksi berhasil diupdate untuk invoice ", fmt.Sprintf("%v", c["invoiceId"]))
				}
			}()

		} else {
			code = "10"
			message = "Error sistem pada QREN"
			status = "Transaksi gagal"
			qrContent = ""
		}
	}

	return code, message, status, qrContent, jamMasuk, jamKeluar, totalJam, amount
}
