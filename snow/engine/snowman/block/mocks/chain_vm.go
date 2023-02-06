// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lasthyphen/dijetsnodego/snow/engine/snowman/block (interfaces: ChainVM)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	manager "github.com/lasthyphen/dijetsnodego/database/manager"
	ids "github.com/lasthyphen/dijetsnodego/ids"
	snow "github.com/lasthyphen/dijetsnodego/snow"
	snowman "github.com/lasthyphen/dijetsnodego/snow/consensus/snowman"
	common "github.com/lasthyphen/dijetsnodego/snow/engine/common"
	version "github.com/lasthyphen/dijetsnodego/version"
	gomock "github.com/golang/mock/gomock"
)

// MockChainVM is a mock of ChainVM interface.
type MockChainVM struct {
	ctrl     *gomock.Controller
	recorder *MockChainVMMockRecorder
}

// MockChainVMMockRecorder is the mock recorder for MockChainVM.
type MockChainVMMockRecorder struct {
	mock *MockChainVM
}

// NewMockChainVM creates a new mock instance.
func NewMockChainVM(ctrl *gomock.Controller) *MockChainVM {
	mock := &MockChainVM{ctrl: ctrl}
	mock.recorder = &MockChainVMMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChainVM) EXPECT() *MockChainVMMockRecorder {
	return m.recorder
}

// AppGossip mocks base method.
func (m *MockChainVM) AppGossip(arg0 context.Context, arg1 ids.NodeID, arg2 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppGossip", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppGossip indicates an expected call of AppGossip.
func (mr *MockChainVMMockRecorder) AppGossip(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppGossip", reflect.TypeOf((*MockChainVM)(nil).AppGossip), arg0, arg1, arg2)
}

// AppRequest mocks base method.
func (m *MockChainVM) AppRequest(arg0 context.Context, arg1 ids.NodeID, arg2 uint32, arg3 time.Time, arg4 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppRequest", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppRequest indicates an expected call of AppRequest.
func (mr *MockChainVMMockRecorder) AppRequest(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppRequest", reflect.TypeOf((*MockChainVM)(nil).AppRequest), arg0, arg1, arg2, arg3, arg4)
}

// AppRequestFailed mocks base method.
func (m *MockChainVM) AppRequestFailed(arg0 context.Context, arg1 ids.NodeID, arg2 uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppRequestFailed", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppRequestFailed indicates an expected call of AppRequestFailed.
func (mr *MockChainVMMockRecorder) AppRequestFailed(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppRequestFailed", reflect.TypeOf((*MockChainVM)(nil).AppRequestFailed), arg0, arg1, arg2)
}

// AppResponse mocks base method.
func (m *MockChainVM) AppResponse(arg0 context.Context, arg1 ids.NodeID, arg2 uint32, arg3 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppResponse", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppResponse indicates an expected call of AppResponse.
func (mr *MockChainVMMockRecorder) AppResponse(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppResponse", reflect.TypeOf((*MockChainVM)(nil).AppResponse), arg0, arg1, arg2, arg3)
}

// BuildBlock mocks base method.
func (m *MockChainVM) BuildBlock(arg0 context.Context) (snowman.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildBlock", arg0)
	ret0, _ := ret[0].(snowman.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildBlock indicates an expected call of BuildBlock.
func (mr *MockChainVMMockRecorder) BuildBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildBlock", reflect.TypeOf((*MockChainVM)(nil).BuildBlock), arg0)
}

// Connected mocks base method.
func (m *MockChainVM) Connected(arg0 context.Context, arg1 ids.NodeID, arg2 *version.Application) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connected", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connected indicates an expected call of Connected.
func (mr *MockChainVMMockRecorder) Connected(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connected", reflect.TypeOf((*MockChainVM)(nil).Connected), arg0, arg1, arg2)
}

// CreateHandlers mocks base method.
func (m *MockChainVM) CreateHandlers(arg0 context.Context) (map[string]*common.HTTPHandler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHandlers", arg0)
	ret0, _ := ret[0].(map[string]*common.HTTPHandler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateHandlers indicates an expected call of CreateHandlers.
func (mr *MockChainVMMockRecorder) CreateHandlers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHandlers", reflect.TypeOf((*MockChainVM)(nil).CreateHandlers), arg0)
}

// CreateStaticHandlers mocks base method.
func (m *MockChainVM) CreateStaticHandlers(arg0 context.Context) (map[string]*common.HTTPHandler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStaticHandlers", arg0)
	ret0, _ := ret[0].(map[string]*common.HTTPHandler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStaticHandlers indicates an expected call of CreateStaticHandlers.
func (mr *MockChainVMMockRecorder) CreateStaticHandlers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStaticHandlers", reflect.TypeOf((*MockChainVM)(nil).CreateStaticHandlers), arg0)
}

// CrossChainAppRequest mocks base method.
func (m *MockChainVM) CrossChainAppRequest(arg0 context.Context, arg1 ids.ID, arg2 uint32, arg3 time.Time, arg4 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CrossChainAppRequest", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// CrossChainAppRequest indicates an expected call of CrossChainAppRequest.
func (mr *MockChainVMMockRecorder) CrossChainAppRequest(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CrossChainAppRequest", reflect.TypeOf((*MockChainVM)(nil).CrossChainAppRequest), arg0, arg1, arg2, arg3, arg4)
}

// CrossChainAppRequestFailed mocks base method.
func (m *MockChainVM) CrossChainAppRequestFailed(arg0 context.Context, arg1 ids.ID, arg2 uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CrossChainAppRequestFailed", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CrossChainAppRequestFailed indicates an expected call of CrossChainAppRequestFailed.
func (mr *MockChainVMMockRecorder) CrossChainAppRequestFailed(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CrossChainAppRequestFailed", reflect.TypeOf((*MockChainVM)(nil).CrossChainAppRequestFailed), arg0, arg1, arg2)
}

// CrossChainAppResponse mocks base method.
func (m *MockChainVM) CrossChainAppResponse(arg0 context.Context, arg1 ids.ID, arg2 uint32, arg3 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CrossChainAppResponse", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// CrossChainAppResponse indicates an expected call of CrossChainAppResponse.
func (mr *MockChainVMMockRecorder) CrossChainAppResponse(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CrossChainAppResponse", reflect.TypeOf((*MockChainVM)(nil).CrossChainAppResponse), arg0, arg1, arg2, arg3)
}

// Disconnected mocks base method.
func (m *MockChainVM) Disconnected(arg0 context.Context, arg1 ids.NodeID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnected", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnected indicates an expected call of Disconnected.
func (mr *MockChainVMMockRecorder) Disconnected(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnected", reflect.TypeOf((*MockChainVM)(nil).Disconnected), arg0, arg1)
}

// GetBlock mocks base method.
func (m *MockChainVM) GetBlock(arg0 context.Context, arg1 ids.ID) (snowman.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", arg0, arg1)
	ret0, _ := ret[0].(snowman.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock.
func (mr *MockChainVMMockRecorder) GetBlock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockChainVM)(nil).GetBlock), arg0, arg1)
}

// HealthCheck mocks base method.
func (m *MockChainVM) HealthCheck(arg0 context.Context) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HealthCheck indicates an expected call of HealthCheck.
func (mr *MockChainVMMockRecorder) HealthCheck(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockChainVM)(nil).HealthCheck), arg0)
}

// Initialize mocks base method.
func (m *MockChainVM) Initialize(arg0 context.Context, arg1 *snow.Context, arg2 manager.Manager, arg3, arg4, arg5 []byte, arg6 chan<- common.Message, arg7 []*common.Fx, arg8 common.AppSender) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockChainVMMockRecorder) Initialize(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockChainVM)(nil).Initialize), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
}

// LastAccepted mocks base method.
func (m *MockChainVM) LastAccepted(arg0 context.Context) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastAccepted", arg0)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastAccepted indicates an expected call of LastAccepted.
func (mr *MockChainVMMockRecorder) LastAccepted(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastAccepted", reflect.TypeOf((*MockChainVM)(nil).LastAccepted), arg0)
}

// ParseBlock mocks base method.
func (m *MockChainVM) ParseBlock(arg0 context.Context, arg1 []byte) (snowman.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseBlock", arg0, arg1)
	ret0, _ := ret[0].(snowman.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseBlock indicates an expected call of ParseBlock.
func (mr *MockChainVMMockRecorder) ParseBlock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseBlock", reflect.TypeOf((*MockChainVM)(nil).ParseBlock), arg0, arg1)
}

// SetPreference mocks base method.
func (m *MockChainVM) SetPreference(arg0 context.Context, arg1 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPreference", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPreference indicates an expected call of SetPreference.
func (mr *MockChainVMMockRecorder) SetPreference(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPreference", reflect.TypeOf((*MockChainVM)(nil).SetPreference), arg0, arg1)
}

// SetState mocks base method.
func (m *MockChainVM) SetState(arg0 context.Context, arg1 snow.State) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetState", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetState indicates an expected call of SetState.
func (mr *MockChainVMMockRecorder) SetState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetState", reflect.TypeOf((*MockChainVM)(nil).SetState), arg0, arg1)
}

// Shutdown mocks base method.
func (m *MockChainVM) Shutdown(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockChainVMMockRecorder) Shutdown(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockChainVM)(nil).Shutdown), arg0)
}

// Version mocks base method.
func (m *MockChainVM) Version(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Version indicates an expected call of Version.
func (mr *MockChainVMMockRecorder) Version(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockChainVM)(nil).Version), arg0)
}
