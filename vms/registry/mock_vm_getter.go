// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lasthyphen/dijetsnodego/vms/registry (interfaces: VMGetter)

// Package registry is a generated GoMock package.
package registry

import (
	reflect "reflect"

	ids "github.com/lasthyphen/dijetsnodego/ids"
	vms "github.com/lasthyphen/dijetsnodego/vms"
	gomock "github.com/golang/mock/gomock"
)

// MockVMGetter is a mock of VMGetter interface.
type MockVMGetter struct {
	ctrl     *gomock.Controller
	recorder *MockVMGetterMockRecorder
}

// MockVMGetterMockRecorder is the mock recorder for MockVMGetter.
type MockVMGetterMockRecorder struct {
	mock *MockVMGetter
}

// NewMockVMGetter creates a new mock instance.
func NewMockVMGetter(ctrl *gomock.Controller) *MockVMGetter {
	mock := &MockVMGetter{ctrl: ctrl}
	mock.recorder = &MockVMGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVMGetter) EXPECT() *MockVMGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockVMGetter) Get() (map[ids.ID]vms.Factory, map[ids.ID]vms.Factory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(map[ids.ID]vms.Factory)
	ret1, _ := ret[1].(map[ids.ID]vms.Factory)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockVMGetterMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVMGetter)(nil).Get))
}
