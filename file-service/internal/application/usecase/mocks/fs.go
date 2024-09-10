// Code generated by MockGen. DO NOT EDIT.
// Source: internal/application/usecase/login-pam.go
//
// Generated by this command:
//
//	mockgen -source=internal/application/usecase/login-pam.go -destination=internal/application/usecase/mocks/fs.go -package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLoginPAMUseCase is a mock of LoginPAMUseCase interface.
type MockLoginPAMUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockLoginPAMUseCaseMockRecorder
}

// MockLoginPAMUseCaseMockRecorder is the mock recorder for MockLoginPAMUseCase.
type MockLoginPAMUseCaseMockRecorder struct {
	mock *MockLoginPAMUseCase
}

// NewMockLoginPAMUseCase creates a new mock instance.
func NewMockLoginPAMUseCase(ctrl *gomock.Controller) *MockLoginPAMUseCase {
	mock := &MockLoginPAMUseCase{ctrl: ctrl}
	mock.recorder = &MockLoginPAMUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoginPAMUseCase) EXPECT() *MockLoginPAMUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockLoginPAMUseCase) Execute(username, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", username, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockLoginPAMUseCaseMockRecorder) Execute(username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockLoginPAMUseCase)(nil).Execute), username, password)
}
