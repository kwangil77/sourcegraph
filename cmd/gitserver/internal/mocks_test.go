// Code generated by go-mockgen 1.3.7; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the mockgen.yaml file in the root of this repository.

package internal

import (
	"context"
	"io"
	"sync"

	api "github.com/sourcegraph/sourcegraph/internal/api"
	protocol "github.com/sourcegraph/sourcegraph/internal/gitserver/protocol"
	trace "github.com/sourcegraph/sourcegraph/internal/trace"
)

// MockService is a mock implementation of the service interface (from the
// package github.com/sourcegraph/sourcegraph/cmd/gitserver/internal) used
// for unit testing.
type MockService struct {
	// CloneRepoFunc is an instance of a mock function object controlling
	// the behavior of the method CloneRepo.
	CloneRepoFunc *ServiceCloneRepoFunc
	// CreateCommitFromPatchFunc is an instance of a mock function object
	// controlling the behavior of the method CreateCommitFromPatch.
	CreateCommitFromPatchFunc *ServiceCreateCommitFromPatchFunc
	// EnsureRevisionFunc is an instance of a mock function object
	// controlling the behavior of the method EnsureRevision.
	EnsureRevisionFunc *ServiceEnsureRevisionFunc
	// IsRepoCloneableFunc is an instance of a mock function object
	// controlling the behavior of the method IsRepoCloneable.
	IsRepoCloneableFunc *ServiceIsRepoCloneableFunc
	// LogIfCorruptFunc is an instance of a mock function object controlling
	// the behavior of the method LogIfCorrupt.
	LogIfCorruptFunc *ServiceLogIfCorruptFunc
	// MaybeStartCloneFunc is an instance of a mock function object
	// controlling the behavior of the method MaybeStartClone.
	MaybeStartCloneFunc *ServiceMaybeStartCloneFunc
	// RepoUpdateFunc is an instance of a mock function object controlling
	// the behavior of the method RepoUpdate.
	RepoUpdateFunc *ServiceRepoUpdateFunc
	// SearchWithObservabilityFunc is an instance of a mock function object
	// controlling the behavior of the method SearchWithObservability.
	SearchWithObservabilityFunc *ServiceSearchWithObservabilityFunc
}

// NewMockService creates a new mock of the service interface. All methods
// return zero values for all results, unless overwritten.
func NewMockService() *MockService {
	return &MockService{
		CloneRepoFunc: &ServiceCloneRepoFunc{
			defaultHook: func(context.Context, api.RepoName, CloneOptions) (r0 string, r1 error) {
				return
			},
		},
		CreateCommitFromPatchFunc: &ServiceCreateCommitFromPatchFunc{
			defaultHook: func(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) (r0 protocol.CreateCommitFromPatchResponse) {
				return
			},
		},
		EnsureRevisionFunc: &ServiceEnsureRevisionFunc{
			defaultHook: func(context.Context, api.RepoName, string) (r0 bool) {
				return
			},
		},
		IsRepoCloneableFunc: &ServiceIsRepoCloneableFunc{
			defaultHook: func(context.Context, api.RepoName) (r0 protocol.IsRepoCloneableResponse, r1 error) {
				return
			},
		},
		LogIfCorruptFunc: &ServiceLogIfCorruptFunc{
			defaultHook: func(context.Context, api.RepoName, error) {
				return
			},
		},
		MaybeStartCloneFunc: &ServiceMaybeStartCloneFunc{
			defaultHook: func(context.Context, api.RepoName) (r0 bool, r1 CloneStatus, r2 error) {
				return
			},
		},
		RepoUpdateFunc: &ServiceRepoUpdateFunc{
			defaultHook: func(context.Context, *protocol.RepoUpdateRequest) (r0 protocol.RepoUpdateResponse) {
				return
			},
		},
		SearchWithObservabilityFunc: &ServiceSearchWithObservabilityFunc{
			defaultHook: func(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (r0 bool, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockService creates a new mock of the service interface. All
// methods panic on invocation, unless overwritten.
func NewStrictMockService() *MockService {
	return &MockService{
		CloneRepoFunc: &ServiceCloneRepoFunc{
			defaultHook: func(context.Context, api.RepoName, CloneOptions) (string, error) {
				panic("unexpected invocation of MockService.CloneRepo")
			},
		},
		CreateCommitFromPatchFunc: &ServiceCreateCommitFromPatchFunc{
			defaultHook: func(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) protocol.CreateCommitFromPatchResponse {
				panic("unexpected invocation of MockService.CreateCommitFromPatch")
			},
		},
		EnsureRevisionFunc: &ServiceEnsureRevisionFunc{
			defaultHook: func(context.Context, api.RepoName, string) bool {
				panic("unexpected invocation of MockService.EnsureRevision")
			},
		},
		IsRepoCloneableFunc: &ServiceIsRepoCloneableFunc{
			defaultHook: func(context.Context, api.RepoName) (protocol.IsRepoCloneableResponse, error) {
				panic("unexpected invocation of MockService.IsRepoCloneable")
			},
		},
		LogIfCorruptFunc: &ServiceLogIfCorruptFunc{
			defaultHook: func(context.Context, api.RepoName, error) {
				panic("unexpected invocation of MockService.LogIfCorrupt")
			},
		},
		MaybeStartCloneFunc: &ServiceMaybeStartCloneFunc{
			defaultHook: func(context.Context, api.RepoName) (bool, CloneStatus, error) {
				panic("unexpected invocation of MockService.MaybeStartClone")
			},
		},
		RepoUpdateFunc: &ServiceRepoUpdateFunc{
			defaultHook: func(context.Context, *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse {
				panic("unexpected invocation of MockService.RepoUpdate")
			},
		},
		SearchWithObservabilityFunc: &ServiceSearchWithObservabilityFunc{
			defaultHook: func(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (bool, error) {
				panic("unexpected invocation of MockService.SearchWithObservability")
			},
		},
	}
}

// surrogateMockService is a copy of the service interface (from the package
// github.com/sourcegraph/sourcegraph/cmd/gitserver/internal). It is
// redefined here as it is unexported in the source package.
type surrogateMockService interface {
	CloneRepo(context.Context, api.RepoName, CloneOptions) (string, error)
	CreateCommitFromPatch(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) protocol.CreateCommitFromPatchResponse
	EnsureRevision(context.Context, api.RepoName, string) bool
	IsRepoCloneable(context.Context, api.RepoName) (protocol.IsRepoCloneableResponse, error)
	LogIfCorrupt(context.Context, api.RepoName, error)
	MaybeStartClone(context.Context, api.RepoName) (bool, CloneStatus, error)
	RepoUpdate(context.Context, *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse
	SearchWithObservability(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (bool, error)
}

// NewMockServiceFrom creates a new mock of the MockService interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockServiceFrom(i surrogateMockService) *MockService {
	return &MockService{
		CloneRepoFunc: &ServiceCloneRepoFunc{
			defaultHook: i.CloneRepo,
		},
		CreateCommitFromPatchFunc: &ServiceCreateCommitFromPatchFunc{
			defaultHook: i.CreateCommitFromPatch,
		},
		EnsureRevisionFunc: &ServiceEnsureRevisionFunc{
			defaultHook: i.EnsureRevision,
		},
		IsRepoCloneableFunc: &ServiceIsRepoCloneableFunc{
			defaultHook: i.IsRepoCloneable,
		},
		LogIfCorruptFunc: &ServiceLogIfCorruptFunc{
			defaultHook: i.LogIfCorrupt,
		},
		MaybeStartCloneFunc: &ServiceMaybeStartCloneFunc{
			defaultHook: i.MaybeStartClone,
		},
		RepoUpdateFunc: &ServiceRepoUpdateFunc{
			defaultHook: i.RepoUpdate,
		},
		SearchWithObservabilityFunc: &ServiceSearchWithObservabilityFunc{
			defaultHook: i.SearchWithObservability,
		},
	}
}

// ServiceCloneRepoFunc describes the behavior when the CloneRepo method of
// the parent MockService instance is invoked.
type ServiceCloneRepoFunc struct {
	defaultHook func(context.Context, api.RepoName, CloneOptions) (string, error)
	hooks       []func(context.Context, api.RepoName, CloneOptions) (string, error)
	history     []ServiceCloneRepoFuncCall
	mutex       sync.Mutex
}

// CloneRepo delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockService) CloneRepo(v0 context.Context, v1 api.RepoName, v2 CloneOptions) (string, error) {
	r0, r1 := m.CloneRepoFunc.nextHook()(v0, v1, v2)
	m.CloneRepoFunc.appendCall(ServiceCloneRepoFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the CloneRepo method of
// the parent MockService instance is invoked and the hook queue is empty.
func (f *ServiceCloneRepoFunc) SetDefaultHook(hook func(context.Context, api.RepoName, CloneOptions) (string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// CloneRepo method of the parent MockService instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ServiceCloneRepoFunc) PushHook(hook func(context.Context, api.RepoName, CloneOptions) (string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ServiceCloneRepoFunc) SetDefaultReturn(r0 string, r1 error) {
	f.SetDefaultHook(func(context.Context, api.RepoName, CloneOptions) (string, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ServiceCloneRepoFunc) PushReturn(r0 string, r1 error) {
	f.PushHook(func(context.Context, api.RepoName, CloneOptions) (string, error) {
		return r0, r1
	})
}

func (f *ServiceCloneRepoFunc) nextHook() func(context.Context, api.RepoName, CloneOptions) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ServiceCloneRepoFunc) appendCall(r0 ServiceCloneRepoFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ServiceCloneRepoFuncCall objects describing
// the invocations of this function.
func (f *ServiceCloneRepoFunc) History() []ServiceCloneRepoFuncCall {
	f.mutex.Lock()
	history := make([]ServiceCloneRepoFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ServiceCloneRepoFuncCall is an object that describes an invocation of
// method CloneRepo on an instance of MockService.
type ServiceCloneRepoFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.RepoName
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 CloneOptions
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ServiceCloneRepoFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ServiceCloneRepoFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ServiceCreateCommitFromPatchFunc describes the behavior when the
// CreateCommitFromPatch method of the parent MockService instance is
// invoked.
type ServiceCreateCommitFromPatchFunc struct {
	defaultHook func(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) protocol.CreateCommitFromPatchResponse
	hooks       []func(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) protocol.CreateCommitFromPatchResponse
	history     []ServiceCreateCommitFromPatchFuncCall
	mutex       sync.Mutex
}

// CreateCommitFromPatch delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockService) CreateCommitFromPatch(v0 context.Context, v1 protocol.CreateCommitFromPatchRequest, v2 io.Reader) protocol.CreateCommitFromPatchResponse {
	r0 := m.CreateCommitFromPatchFunc.nextHook()(v0, v1, v2)
	m.CreateCommitFromPatchFunc.appendCall(ServiceCreateCommitFromPatchFuncCall{v0, v1, v2, r0})
	return r0
}

// SetDefaultHook sets function that is called when the
// CreateCommitFromPatch method of the parent MockService instance is
// invoked and the hook queue is empty.
func (f *ServiceCreateCommitFromPatchFunc) SetDefaultHook(hook func(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) protocol.CreateCommitFromPatchResponse) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// CreateCommitFromPatch method of the parent MockService instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *ServiceCreateCommitFromPatchFunc) PushHook(hook func(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) protocol.CreateCommitFromPatchResponse) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ServiceCreateCommitFromPatchFunc) SetDefaultReturn(r0 protocol.CreateCommitFromPatchResponse) {
	f.SetDefaultHook(func(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) protocol.CreateCommitFromPatchResponse {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ServiceCreateCommitFromPatchFunc) PushReturn(r0 protocol.CreateCommitFromPatchResponse) {
	f.PushHook(func(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) protocol.CreateCommitFromPatchResponse {
		return r0
	})
}

func (f *ServiceCreateCommitFromPatchFunc) nextHook() func(context.Context, protocol.CreateCommitFromPatchRequest, io.Reader) protocol.CreateCommitFromPatchResponse {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ServiceCreateCommitFromPatchFunc) appendCall(r0 ServiceCreateCommitFromPatchFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ServiceCreateCommitFromPatchFuncCall
// objects describing the invocations of this function.
func (f *ServiceCreateCommitFromPatchFunc) History() []ServiceCreateCommitFromPatchFuncCall {
	f.mutex.Lock()
	history := make([]ServiceCreateCommitFromPatchFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ServiceCreateCommitFromPatchFuncCall is an object that describes an
// invocation of method CreateCommitFromPatch on an instance of MockService.
type ServiceCreateCommitFromPatchFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 protocol.CreateCommitFromPatchRequest
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 io.Reader
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 protocol.CreateCommitFromPatchResponse
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ServiceCreateCommitFromPatchFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ServiceCreateCommitFromPatchFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ServiceEnsureRevisionFunc describes the behavior when the EnsureRevision
// method of the parent MockService instance is invoked.
type ServiceEnsureRevisionFunc struct {
	defaultHook func(context.Context, api.RepoName, string) bool
	hooks       []func(context.Context, api.RepoName, string) bool
	history     []ServiceEnsureRevisionFuncCall
	mutex       sync.Mutex
}

// EnsureRevision delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockService) EnsureRevision(v0 context.Context, v1 api.RepoName, v2 string) bool {
	r0 := m.EnsureRevisionFunc.nextHook()(v0, v1, v2)
	m.EnsureRevisionFunc.appendCall(ServiceEnsureRevisionFuncCall{v0, v1, v2, r0})
	return r0
}

// SetDefaultHook sets function that is called when the EnsureRevision
// method of the parent MockService instance is invoked and the hook queue
// is empty.
func (f *ServiceEnsureRevisionFunc) SetDefaultHook(hook func(context.Context, api.RepoName, string) bool) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// EnsureRevision method of the parent MockService instance invokes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ServiceEnsureRevisionFunc) PushHook(hook func(context.Context, api.RepoName, string) bool) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ServiceEnsureRevisionFunc) SetDefaultReturn(r0 bool) {
	f.SetDefaultHook(func(context.Context, api.RepoName, string) bool {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ServiceEnsureRevisionFunc) PushReturn(r0 bool) {
	f.PushHook(func(context.Context, api.RepoName, string) bool {
		return r0
	})
}

func (f *ServiceEnsureRevisionFunc) nextHook() func(context.Context, api.RepoName, string) bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ServiceEnsureRevisionFunc) appendCall(r0 ServiceEnsureRevisionFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ServiceEnsureRevisionFuncCall objects
// describing the invocations of this function.
func (f *ServiceEnsureRevisionFunc) History() []ServiceEnsureRevisionFuncCall {
	f.mutex.Lock()
	history := make([]ServiceEnsureRevisionFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ServiceEnsureRevisionFuncCall is an object that describes an invocation
// of method EnsureRevision on an instance of MockService.
type ServiceEnsureRevisionFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.RepoName
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 bool
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ServiceEnsureRevisionFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ServiceEnsureRevisionFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ServiceIsRepoCloneableFunc describes the behavior when the
// IsRepoCloneable method of the parent MockService instance is invoked.
type ServiceIsRepoCloneableFunc struct {
	defaultHook func(context.Context, api.RepoName) (protocol.IsRepoCloneableResponse, error)
	hooks       []func(context.Context, api.RepoName) (protocol.IsRepoCloneableResponse, error)
	history     []ServiceIsRepoCloneableFuncCall
	mutex       sync.Mutex
}

// IsRepoCloneable delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockService) IsRepoCloneable(v0 context.Context, v1 api.RepoName) (protocol.IsRepoCloneableResponse, error) {
	r0, r1 := m.IsRepoCloneableFunc.nextHook()(v0, v1)
	m.IsRepoCloneableFunc.appendCall(ServiceIsRepoCloneableFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the IsRepoCloneable
// method of the parent MockService instance is invoked and the hook queue
// is empty.
func (f *ServiceIsRepoCloneableFunc) SetDefaultHook(hook func(context.Context, api.RepoName) (protocol.IsRepoCloneableResponse, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// IsRepoCloneable method of the parent MockService instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ServiceIsRepoCloneableFunc) PushHook(hook func(context.Context, api.RepoName) (protocol.IsRepoCloneableResponse, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ServiceIsRepoCloneableFunc) SetDefaultReturn(r0 protocol.IsRepoCloneableResponse, r1 error) {
	f.SetDefaultHook(func(context.Context, api.RepoName) (protocol.IsRepoCloneableResponse, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ServiceIsRepoCloneableFunc) PushReturn(r0 protocol.IsRepoCloneableResponse, r1 error) {
	f.PushHook(func(context.Context, api.RepoName) (protocol.IsRepoCloneableResponse, error) {
		return r0, r1
	})
}

func (f *ServiceIsRepoCloneableFunc) nextHook() func(context.Context, api.RepoName) (protocol.IsRepoCloneableResponse, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ServiceIsRepoCloneableFunc) appendCall(r0 ServiceIsRepoCloneableFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ServiceIsRepoCloneableFuncCall objects
// describing the invocations of this function.
func (f *ServiceIsRepoCloneableFunc) History() []ServiceIsRepoCloneableFuncCall {
	f.mutex.Lock()
	history := make([]ServiceIsRepoCloneableFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ServiceIsRepoCloneableFuncCall is an object that describes an invocation
// of method IsRepoCloneable on an instance of MockService.
type ServiceIsRepoCloneableFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.RepoName
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 protocol.IsRepoCloneableResponse
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ServiceIsRepoCloneableFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ServiceIsRepoCloneableFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ServiceLogIfCorruptFunc describes the behavior when the LogIfCorrupt
// method of the parent MockService instance is invoked.
type ServiceLogIfCorruptFunc struct {
	defaultHook func(context.Context, api.RepoName, error)
	hooks       []func(context.Context, api.RepoName, error)
	history     []ServiceLogIfCorruptFuncCall
	mutex       sync.Mutex
}

// LogIfCorrupt delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockService) LogIfCorrupt(v0 context.Context, v1 api.RepoName, v2 error) {
	m.LogIfCorruptFunc.nextHook()(v0, v1, v2)
	m.LogIfCorruptFunc.appendCall(ServiceLogIfCorruptFuncCall{v0, v1, v2})
	return
}

// SetDefaultHook sets function that is called when the LogIfCorrupt method
// of the parent MockService instance is invoked and the hook queue is
// empty.
func (f *ServiceLogIfCorruptFunc) SetDefaultHook(hook func(context.Context, api.RepoName, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// LogIfCorrupt method of the parent MockService instance invokes the hook
// at the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ServiceLogIfCorruptFunc) PushHook(hook func(context.Context, api.RepoName, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ServiceLogIfCorruptFunc) SetDefaultReturn() {
	f.SetDefaultHook(func(context.Context, api.RepoName, error) {
		return
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ServiceLogIfCorruptFunc) PushReturn() {
	f.PushHook(func(context.Context, api.RepoName, error) {
		return
	})
}

func (f *ServiceLogIfCorruptFunc) nextHook() func(context.Context, api.RepoName, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ServiceLogIfCorruptFunc) appendCall(r0 ServiceLogIfCorruptFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ServiceLogIfCorruptFuncCall objects
// describing the invocations of this function.
func (f *ServiceLogIfCorruptFunc) History() []ServiceLogIfCorruptFuncCall {
	f.mutex.Lock()
	history := make([]ServiceLogIfCorruptFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ServiceLogIfCorruptFuncCall is an object that describes an invocation of
// method LogIfCorrupt on an instance of MockService.
type ServiceLogIfCorruptFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.RepoName
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ServiceLogIfCorruptFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ServiceLogIfCorruptFuncCall) Results() []interface{} {
	return []interface{}{}
}

// ServiceMaybeStartCloneFunc describes the behavior when the
// MaybeStartClone method of the parent MockService instance is invoked.
type ServiceMaybeStartCloneFunc struct {
	defaultHook func(context.Context, api.RepoName) (bool, CloneStatus, error)
	hooks       []func(context.Context, api.RepoName) (bool, CloneStatus, error)
	history     []ServiceMaybeStartCloneFuncCall
	mutex       sync.Mutex
}

// MaybeStartClone delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockService) MaybeStartClone(v0 context.Context, v1 api.RepoName) (bool, CloneStatus, error) {
	r0, r1, r2 := m.MaybeStartCloneFunc.nextHook()(v0, v1)
	m.MaybeStartCloneFunc.appendCall(ServiceMaybeStartCloneFuncCall{v0, v1, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the MaybeStartClone
// method of the parent MockService instance is invoked and the hook queue
// is empty.
func (f *ServiceMaybeStartCloneFunc) SetDefaultHook(hook func(context.Context, api.RepoName) (bool, CloneStatus, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// MaybeStartClone method of the parent MockService instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ServiceMaybeStartCloneFunc) PushHook(hook func(context.Context, api.RepoName) (bool, CloneStatus, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ServiceMaybeStartCloneFunc) SetDefaultReturn(r0 bool, r1 CloneStatus, r2 error) {
	f.SetDefaultHook(func(context.Context, api.RepoName) (bool, CloneStatus, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ServiceMaybeStartCloneFunc) PushReturn(r0 bool, r1 CloneStatus, r2 error) {
	f.PushHook(func(context.Context, api.RepoName) (bool, CloneStatus, error) {
		return r0, r1, r2
	})
}

func (f *ServiceMaybeStartCloneFunc) nextHook() func(context.Context, api.RepoName) (bool, CloneStatus, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ServiceMaybeStartCloneFunc) appendCall(r0 ServiceMaybeStartCloneFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ServiceMaybeStartCloneFuncCall objects
// describing the invocations of this function.
func (f *ServiceMaybeStartCloneFunc) History() []ServiceMaybeStartCloneFuncCall {
	f.mutex.Lock()
	history := make([]ServiceMaybeStartCloneFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ServiceMaybeStartCloneFuncCall is an object that describes an invocation
// of method MaybeStartClone on an instance of MockService.
type ServiceMaybeStartCloneFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 api.RepoName
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 bool
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 CloneStatus
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ServiceMaybeStartCloneFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ServiceMaybeStartCloneFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}

// ServiceRepoUpdateFunc describes the behavior when the RepoUpdate method
// of the parent MockService instance is invoked.
type ServiceRepoUpdateFunc struct {
	defaultHook func(context.Context, *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse
	hooks       []func(context.Context, *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse
	history     []ServiceRepoUpdateFuncCall
	mutex       sync.Mutex
}

// RepoUpdate delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockService) RepoUpdate(v0 context.Context, v1 *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse {
	r0 := m.RepoUpdateFunc.nextHook()(v0, v1)
	m.RepoUpdateFunc.appendCall(ServiceRepoUpdateFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the RepoUpdate method of
// the parent MockService instance is invoked and the hook queue is empty.
func (f *ServiceRepoUpdateFunc) SetDefaultHook(hook func(context.Context, *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// RepoUpdate method of the parent MockService instance invokes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ServiceRepoUpdateFunc) PushHook(hook func(context.Context, *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ServiceRepoUpdateFunc) SetDefaultReturn(r0 protocol.RepoUpdateResponse) {
	f.SetDefaultHook(func(context.Context, *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ServiceRepoUpdateFunc) PushReturn(r0 protocol.RepoUpdateResponse) {
	f.PushHook(func(context.Context, *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse {
		return r0
	})
}

func (f *ServiceRepoUpdateFunc) nextHook() func(context.Context, *protocol.RepoUpdateRequest) protocol.RepoUpdateResponse {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ServiceRepoUpdateFunc) appendCall(r0 ServiceRepoUpdateFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ServiceRepoUpdateFuncCall objects
// describing the invocations of this function.
func (f *ServiceRepoUpdateFunc) History() []ServiceRepoUpdateFuncCall {
	f.mutex.Lock()
	history := make([]ServiceRepoUpdateFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ServiceRepoUpdateFuncCall is an object that describes an invocation of
// method RepoUpdate on an instance of MockService.
type ServiceRepoUpdateFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 *protocol.RepoUpdateRequest
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 protocol.RepoUpdateResponse
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ServiceRepoUpdateFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ServiceRepoUpdateFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// ServiceSearchWithObservabilityFunc describes the behavior when the
// SearchWithObservability method of the parent MockService instance is
// invoked.
type ServiceSearchWithObservabilityFunc struct {
	defaultHook func(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (bool, error)
	hooks       []func(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (bool, error)
	history     []ServiceSearchWithObservabilityFuncCall
	mutex       sync.Mutex
}

// SearchWithObservability delegates to the next hook function in the queue
// and stores the parameter and result values of this invocation.
func (m *MockService) SearchWithObservability(v0 context.Context, v1 trace.Trace, v2 *protocol.SearchRequest, v3 func(*protocol.CommitMatch) error) (bool, error) {
	r0, r1 := m.SearchWithObservabilityFunc.nextHook()(v0, v1, v2, v3)
	m.SearchWithObservabilityFunc.appendCall(ServiceSearchWithObservabilityFuncCall{v0, v1, v2, v3, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the
// SearchWithObservability method of the parent MockService instance is
// invoked and the hook queue is empty.
func (f *ServiceSearchWithObservabilityFunc) SetDefaultHook(hook func(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// SearchWithObservability method of the parent MockService instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *ServiceSearchWithObservabilityFunc) PushHook(hook func(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *ServiceSearchWithObservabilityFunc) SetDefaultReturn(r0 bool, r1 error) {
	f.SetDefaultHook(func(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (bool, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *ServiceSearchWithObservabilityFunc) PushReturn(r0 bool, r1 error) {
	f.PushHook(func(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (bool, error) {
		return r0, r1
	})
}

func (f *ServiceSearchWithObservabilityFunc) nextHook() func(context.Context, trace.Trace, *protocol.SearchRequest, func(*protocol.CommitMatch) error) (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ServiceSearchWithObservabilityFunc) appendCall(r0 ServiceSearchWithObservabilityFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ServiceSearchWithObservabilityFuncCall
// objects describing the invocations of this function.
func (f *ServiceSearchWithObservabilityFunc) History() []ServiceSearchWithObservabilityFuncCall {
	f.mutex.Lock()
	history := make([]ServiceSearchWithObservabilityFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ServiceSearchWithObservabilityFuncCall is an object that describes an
// invocation of method SearchWithObservability on an instance of
// MockService.
type ServiceSearchWithObservabilityFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 trace.Trace
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 *protocol.SearchRequest
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 func(*protocol.CommitMatch) error
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 bool
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ServiceSearchWithObservabilityFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ServiceSearchWithObservabilityFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
