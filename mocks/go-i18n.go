package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// MockLocalizer is a mock of Localizer interface.
type MockLocalizer struct {
	ctrl     *gomock.Controller
	recorder *MockLocalizerMockRecorder
}

// MockLocalizerMockRecorder is the mock recorder for MockLocalizer.
type MockLocalizerMockRecorder struct {
	mock *MockLocalizer
}

// NewMockLocalizer creates a new mock instance.
func NewMockLocalizer(ctrl *gomock.Controller) *MockLocalizer {
	mock := &MockLocalizer{ctrl: ctrl}
	mock.recorder = &MockLocalizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocalizer) EXPECT() *MockLocalizerMockRecorder {
	return m.recorder
}

// MustLocalize mocks base method.
func (m *MockLocalizer) MustLocalize(config *i18n.LocalizeConfig) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MustLocalize", config)
	ret0, _ := ret[0].(string)
	return ret0
}

// MustLocalize indicates an expected call of MustLocalize.
func (mr *MockLocalizerMockRecorder) MustLocalize(chatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustLocalize", reflect.TypeOf((*MockLocalizer)(nil).MustLocalize), chatID)
}
