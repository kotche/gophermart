// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	jwtauth "github.com/go-chi/jwtauth/v5"
	gomock "github.com/golang/mock/gomock"
	model "github.com/kotche/gophermart/internal/model"
)

// MockAuthServiceContract is a mock of AuthServiceContract interface.
type MockAuthServiceContract struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceContractMockRecorder
}

// MockAuthServiceContractMockRecorder is the mock recorder for MockAuthServiceContract.
type MockAuthServiceContractMockRecorder struct {
	mock *MockAuthServiceContract
}

// NewMockAuthServiceContract creates a new mock instance.
func NewMockAuthServiceContract(ctrl *gomock.Controller) *MockAuthServiceContract {
	mock := &MockAuthServiceContract{ctrl: ctrl}
	mock.recorder = &MockAuthServiceContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServiceContract) EXPECT() *MockAuthServiceContractMockRecorder {
	return m.recorder
}

// AuthenticationUser mocks base method.
func (m *MockAuthServiceContract) AuthenticationUser(ctx context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticationUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// AuthenticationUser indicates an expected call of AuthenticationUser.
func (mr *MockAuthServiceContractMockRecorder) AuthenticationUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticationUser", reflect.TypeOf((*MockAuthServiceContract)(nil).AuthenticationUser), ctx, user)
}

// CreateUser mocks base method.
func (m *MockAuthServiceContract) CreateUser(ctx context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthServiceContractMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthServiceContract)(nil).CreateUser), ctx, user)
}

// GenerateToken mocks base method.
func (m *MockAuthServiceContract) GenerateToken(user *model.User, tokenAuth *jwtauth.JWTAuth) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", user, tokenAuth)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthServiceContractMockRecorder) GenerateToken(user, tokenAuth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuthServiceContract)(nil).GenerateToken), user, tokenAuth)
}

// MockAccrualOrderServiceContract is a mock of AccrualOrderServiceContract interface.
type MockAccrualOrderServiceContract struct {
	ctrl     *gomock.Controller
	recorder *MockAccrualOrderServiceContractMockRecorder
}

// MockAccrualOrderServiceContractMockRecorder is the mock recorder for MockAccrualOrderServiceContract.
type MockAccrualOrderServiceContractMockRecorder struct {
	mock *MockAccrualOrderServiceContract
}

// NewMockAccrualOrderServiceContract creates a new mock instance.
func NewMockAccrualOrderServiceContract(ctrl *gomock.Controller) *MockAccrualOrderServiceContract {
	mock := &MockAccrualOrderServiceContract{ctrl: ctrl}
	mock.recorder = &MockAccrualOrderServiceContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccrualOrderServiceContract) EXPECT() *MockAccrualOrderServiceContractMockRecorder {
	return m.recorder
}

// CheckLuhn mocks base method.
func (m *MockAccrualOrderServiceContract) CheckLuhn(number uint64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLuhn", number)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckLuhn indicates an expected call of CheckLuhn.
func (mr *MockAccrualOrderServiceContractMockRecorder) CheckLuhn(number interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLuhn", reflect.TypeOf((*MockAccrualOrderServiceContract)(nil).CheckLuhn), number)
}

// GetUploadedOrders mocks base method.
func (m *MockAccrualOrderServiceContract) GetUploadedOrders(ctx context.Context, userID int) ([]model.AccrualOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUploadedOrders", ctx, userID)
	ret0, _ := ret[0].([]model.AccrualOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUploadedOrders indicates an expected call of GetUploadedOrders.
func (mr *MockAccrualOrderServiceContractMockRecorder) GetUploadedOrders(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUploadedOrders", reflect.TypeOf((*MockAccrualOrderServiceContract)(nil).GetUploadedOrders), ctx, userID)
}

// LoadOrder mocks base method.
func (m *MockAccrualOrderServiceContract) LoadOrder(ctx context.Context, numOrder uint64, userID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadOrder", ctx, numOrder, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadOrder indicates an expected call of LoadOrder.
func (mr *MockAccrualOrderServiceContractMockRecorder) LoadOrder(ctx, numOrder, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadOrder", reflect.TypeOf((*MockAccrualOrderServiceContract)(nil).LoadOrder), ctx, numOrder, userID)
}

// MockWithdrawOrderServiceContract is a mock of WithdrawOrderServiceContract interface.
type MockWithdrawOrderServiceContract struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawOrderServiceContractMockRecorder
}

// MockWithdrawOrderServiceContractMockRecorder is the mock recorder for MockWithdrawOrderServiceContract.
type MockWithdrawOrderServiceContractMockRecorder struct {
	mock *MockWithdrawOrderServiceContract
}

// NewMockWithdrawOrderServiceContract creates a new mock instance.
func NewMockWithdrawOrderServiceContract(ctrl *gomock.Controller) *MockWithdrawOrderServiceContract {
	mock := &MockWithdrawOrderServiceContract{ctrl: ctrl}
	mock.recorder = &MockWithdrawOrderServiceContractMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdrawOrderServiceContract) EXPECT() *MockWithdrawOrderServiceContractMockRecorder {
	return m.recorder
}

// DeductionOfPoints mocks base method.
func (m *MockWithdrawOrderServiceContract) DeductionOfPoints(ctx context.Context, order *model.WithdrawOrder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeductionOfPoints", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeductionOfPoints indicates an expected call of DeductionOfPoints.
func (mr *MockWithdrawOrderServiceContractMockRecorder) DeductionOfPoints(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeductionOfPoints", reflect.TypeOf((*MockWithdrawOrderServiceContract)(nil).DeductionOfPoints), ctx, order)
}

// GetBalance mocks base method.
func (m *MockWithdrawOrderServiceContract) GetBalance(ctx context.Context, userID int) (float32, float32) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, userID)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(float32)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockWithdrawOrderServiceContractMockRecorder) GetBalance(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockWithdrawOrderServiceContract)(nil).GetBalance), ctx, userID)
}

// GetWithdrawalOfPoints mocks base method.
func (m *MockWithdrawOrderServiceContract) GetWithdrawalOfPoints(ctx context.Context, userID int) ([]model.WithdrawOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawalOfPoints", ctx, userID)
	ret0, _ := ret[0].([]model.WithdrawOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawalOfPoints indicates an expected call of GetWithdrawalOfPoints.
func (mr *MockWithdrawOrderServiceContractMockRecorder) GetWithdrawalOfPoints(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawalOfPoints", reflect.TypeOf((*MockWithdrawOrderServiceContract)(nil).GetWithdrawalOfPoints), ctx, userID)
}
