package handler

import (
	"net/http"
	"testing"
)

func TestHandler_deductionOfPoints(t *testing.T) {
	type fields struct {
		Service   *service.Service
		TokenAuth *jwtauth.JWTAuth
		log       *zerolog.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Service:   tt.fields.Service,
				TokenAuth: tt.fields.TokenAuth,
				log:       tt.fields.log,
			}
			h.deductionOfPoints(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_getCurrentBalance(t *testing.T) {
	type fields struct {
		Service   *service.Service
		TokenAuth *jwtauth.JWTAuth
		log       *zerolog.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Service:   tt.fields.Service,
				TokenAuth: tt.fields.TokenAuth,
				log:       tt.fields.log,
			}
			h.getCurrentBalance(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_getWithdrawalOfPoints(t *testing.T) {
	type fields struct {
		Service   *service.Service
		TokenAuth *jwtauth.JWTAuth
		log       *zerolog.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				Service:   tt.fields.Service,
				TokenAuth: tt.fields.TokenAuth,
				log:       tt.fields.log,
			}
			h.getWithdrawalOfPoints(tt.args.w, tt.args.r)
		})
	}
}
