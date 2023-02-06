// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lasthyphen/dijetsnodego/vms/platformvm/txs (interfaces: Staker)

// Package txs is a generated GoMock package.
package txs

import (
	reflect "reflect"
	time "time"

	ids "github.com/lasthyphen/dijetsnodego/ids"
	bls "github.com/lasthyphen/dijetsnodego/utils/crypto/bls"
	gomock "github.com/golang/mock/gomock"
)

// MockStaker is a mock of Staker interface.
type MockStaker struct {
	ctrl     *gomock.Controller
	recorder *MockStakerMockRecorder
}

// MockStakerMockRecorder is the mock recorder for MockStaker.
type MockStakerMockRecorder struct {
	mock *MockStaker
}

// NewMockStaker creates a new mock instance.
func NewMockStaker(ctrl *gomock.Controller) *MockStaker {
	mock := &MockStaker{ctrl: ctrl}
	mock.recorder = &MockStakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStaker) EXPECT() *MockStakerMockRecorder {
	return m.recorder
}

// CurrentPriority mocks base method.
func (m *MockStaker) CurrentPriority() Priority {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentPriority")
	ret0, _ := ret[0].(Priority)
	return ret0
}

// CurrentPriority indicates an expected call of CurrentPriority.
func (mr *MockStakerMockRecorder) CurrentPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentPriority", reflect.TypeOf((*MockStaker)(nil).CurrentPriority))
}

// EndTime mocks base method.
func (m *MockStaker) EndTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EndTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// EndTime indicates an expected call of EndTime.
func (mr *MockStakerMockRecorder) EndTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndTime", reflect.TypeOf((*MockStaker)(nil).EndTime))
}

// NodeID mocks base method.
func (m *MockStaker) NodeID() ids.NodeID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NodeID")
	ret0, _ := ret[0].(ids.NodeID)
	return ret0
}

// NodeID indicates an expected call of NodeID.
func (mr *MockStakerMockRecorder) NodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeID", reflect.TypeOf((*MockStaker)(nil).NodeID))
}

// PendingPriority mocks base method.
func (m *MockStaker) PendingPriority() Priority {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingPriority")
	ret0, _ := ret[0].(Priority)
	return ret0
}

// PendingPriority indicates an expected call of PendingPriority.
func (mr *MockStakerMockRecorder) PendingPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingPriority", reflect.TypeOf((*MockStaker)(nil).PendingPriority))
}

// PublicKey mocks base method.
func (m *MockStaker) PublicKey() (*bls.PublicKey, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicKey")
	ret0, _ := ret[0].(*bls.PublicKey)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PublicKey indicates an expected call of PublicKey.
func (mr *MockStakerMockRecorder) PublicKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicKey", reflect.TypeOf((*MockStaker)(nil).PublicKey))
}

// StartTime mocks base method.
func (m *MockStaker) StartTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// StartTime indicates an expected call of StartTime.
func (mr *MockStakerMockRecorder) StartTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTime", reflect.TypeOf((*MockStaker)(nil).StartTime))
}

// SubnetID mocks base method.
func (m *MockStaker) SubnetID() ids.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubnetID")
	ret0, _ := ret[0].(ids.ID)
	return ret0
}

// SubnetID indicates an expected call of SubnetID.
func (mr *MockStakerMockRecorder) SubnetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubnetID", reflect.TypeOf((*MockStaker)(nil).SubnetID))
}

// Weight mocks base method.
func (m *MockStaker) Weight() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Weight")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Weight indicates an expected call of Weight.
func (mr *MockStakerMockRecorder) Weight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Weight", reflect.TypeOf((*MockStaker)(nil).Weight))
}
