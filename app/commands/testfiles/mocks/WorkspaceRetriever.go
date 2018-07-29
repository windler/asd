// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"

// WorkspaceRetriever is an autogenerated mock type for the WorkspaceRetriever type
type WorkspaceRetriever struct {
	mock.Mock
}

// GetCurrentWorkspace provides a mock function with given fields: root
func (_m *WorkspaceRetriever) GetCurrentWorkspace(root string) string {
	ret := _m.Called(root)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(root)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetWorkspaceByPattern provides a mock function with given fields: root, pattern
func (_m *WorkspaceRetriever) GetWorkspaceByPattern(root string, pattern string) string {
	ret := _m.Called(root, pattern)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(root, pattern)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetWorkspacesIn provides a mock function with given fields: root
func (_m *WorkspaceRetriever) GetWorkspacesIn(root string) []string {
	ret := _m.Called(root)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(root)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
