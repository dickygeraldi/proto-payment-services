package controllers

import (
	"context"
	"database/sql"
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
	timeRequest := time.Now().Format("2006-01-02 15:04:05")
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
	// timeRequest := time.Now().Format("2006-01-02 15:04:05")
	// data, _ := peer.FromContext(ctx)
	// var code, status, message, qrContent string

	// message, statusValidation := validation.ParkingRegistration(req.Api, req.PlatNo, data.Addr.String(), timeRequest)

	// if statusValidation == false {
	// 	code = "05"
	// 	status = "Validasi gagal"
	// 	message = "Data harus diisi"
	// } else {
	// 	if err := s.CheckApi(req.Api); err != nil {
	// 		return nil, err
	// 	} else {
	// 		status = "Transaksi berhasil di proses"
	// 		code, message, status, qrContent = models.ValidationParking(req.PlatNo, timeRequest, s.db, ctx)
	// 	}
	// }

	return &v1.ValidationParkingResponse{
		Message: "message",
		Code:    "code",
		Status:  "status",
		Data: &v1.ValidationParkingData{
			QrContent: "qrContent",
		},
	}, nil
}
