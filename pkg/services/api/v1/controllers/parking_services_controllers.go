package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"os"
	v1 "proto-parking-services/pkg/api/v1"
	"proto-parking-services/pkg/services/api/v1/models"
	"proto-parking-services/pkg/services/api/v1/validation"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UserServices implemented on version 1 proto interface
type parkingServices struct {
	db *sql.DB
}

// socket connection handle
func socketHandle() {
	l, err := net.Listen("tcp", os.Getenv("SOCKET_HOST"))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()
	for {
		// listen all incoming message
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Socket handle request
func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	data := make(map[string]interface{})

	json.Unmarshal(buf, &data)
	fmt.Println(data)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}

// New sending otp services create sending otp service
func NewUserServicesService(db *sql.DB) v1.ParkingServicesServer {
	return &parkingServices{db: db}
}

// Checking Api version
func (s *parkingServices) CheckApi(api string) error {
	if len(api) > 0 {
		if os.Getenv("API_VERSION") != api {
			return status.Errorf(codes.Unimplemented, "Unsupported API Version: Service API implement using '%s', but asked for '%s'", os.Getenv("API_VERSION"), api)
		}
	}
	return nil
}

// Func percobaan header
func (s *parkingServices) RegisterParkingServices(ctx context.Context, req *v1.RegisterParkingRequest) (*v1.RegisterParkingResponse, error) {
	timeRequest := time.Now()
	data, _ := peer.FromContext(ctx)
	var code, status, message, invoiceNo, enteredDate, statusParking string

	message, statusValidation := validation.ParkingRegistration(req.Api, req.PlatNo, data.Addr.String(), timeRequest)

	if statusValidation == false {
		code = "05"
		status = "Validasi gagal"
		message = "Data harus diisi"
	} else {
		if err := s.CheckApi(req.Api); err != nil {
			return nil, err
		} else {
			status = "Transaksi berhasil di proses"
			code, message, status, invoiceNo, enteredDate, statusParking = models.ParkingRegistration(req.PlatNo, timeRequest, s.db, ctx)
		}
	}

	return &v1.RegisterParkingResponse{
		Message: message,
		Code:    code,
		Status:  status,
		Data: &v1.RegisterParkingData{
			StatusParking: statusParking,
			InvoiceNo:     invoiceNo,
			EnteredDate:   enteredDate,
		},
	}, nil
}

// Func percobaan header
func (s *parkingServices) ParkingValidationServices(ctx context.Context, req *v1.RegisterParkingRequest) (*v1.ValidationParkingResponse, error) {
	timeRequest := time.Now()
	data, _ := peer.FromContext(ctx)
	var code, status, message, qrContent string

	message, statusValidation := validation.ParkingRegistration(req.Api, req.PlatNo, data.Addr.String(), timeRequest)

	if statusValidation == false {
		code = "05"
		status = "Validasi gagal"
		message = "Data harus diisi"
	} else {
		if err := s.CheckApi(req.Api); err != nil {
			return nil, err
		} else {
			status = "Transaksi berhasil di proses"
			code, message, status, qrContent = models.ValidationParking(req.PlatNo, timeRequest, s.db, ctx)
		}
	}

	return &v1.ValidationParkingResponse{
		Message: message,
		Code:    code,
		Status:  status,
		Data: &v1.ValidationParkingData{
			QrContent: qrContent,
		},
	}, nil
}
