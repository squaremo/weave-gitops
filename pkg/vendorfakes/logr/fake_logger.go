package fakelogr

import (
	"sync"

	"github.com/go-logr/logr"
)

type FakeLogger struct {
	EnabledStub        func() bool
	enabledMutex       sync.RWMutex
	enabledArgsForCall []struct {
	}
	enabledReturns struct {
		result1 bool
	}
	enabledReturnsOnCall map[int]struct {
		result1 bool
	}
	ErrorStub        func(error, string, ...interface{})
	errorMutex       sync.RWMutex
	errorArgsForCall []struct {
		arg1 error
		arg2 string
		arg3 []interface{}
	}
	InfoStub        func(string, ...interface{})
	infoMutex       sync.RWMutex
	infoArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	VStub        func(int) logr.Logger
	vMutex       sync.RWMutex
	vArgsForCall []struct {
		arg1 int
	}
	vReturns struct {
		result1 logr.Logger
	}
	vReturnsOnCall map[int]struct {
		result1 logr.Logger
	}
	WithNameStub        func(string) logr.Logger
	withNameMutex       sync.RWMutex
	withNameArgsForCall []struct {
		arg1 string
	}
	withNameReturns struct {
		result1 logr.Logger
	}
	withNameReturnsOnCall map[int]struct {
		result1 logr.Logger
	}
	WithValuesStub        func(...interface{}) logr.Logger
	withValuesMutex       sync.RWMutex
	withValuesArgsForCall []struct {
		arg1 []interface{}
	}
	withValuesReturns struct {
		result1 logr.Logger
	}
	withValuesReturnsOnCall map[int]struct {
		result1 logr.Logger
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLogger) Enabled() bool {
	fake.enabledMutex.Lock()
	ret, specificReturn := fake.enabledReturnsOnCall[len(fake.enabledArgsForCall)]
	fake.enabledArgsForCall = append(fake.enabledArgsForCall, struct {
	}{})
	stub := fake.EnabledStub
	fakeReturns := fake.enabledReturns
	fake.recordInvocation("Enabled", []interface{}{})
	fake.enabledMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLogger) EnabledCallCount() int {
	fake.enabledMutex.RLock()
	defer fake.enabledMutex.RUnlock()
	return len(fake.enabledArgsForCall)
}

func (fake *FakeLogger) EnabledCalls(stub func() bool) {
	fake.enabledMutex.Lock()
	defer fake.enabledMutex.Unlock()
	fake.EnabledStub = stub
}

func (fake *FakeLogger) EnabledReturns(result1 bool) {
	fake.enabledMutex.Lock()
	defer fake.enabledMutex.Unlock()
	fake.EnabledStub = nil
	fake.enabledReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeLogger) EnabledReturnsOnCall(i int, result1 bool) {
	fake.enabledMutex.Lock()
	defer fake.enabledMutex.Unlock()
	fake.EnabledStub = nil
	if fake.enabledReturnsOnCall == nil {
		fake.enabledReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.enabledReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeLogger) Error(arg1 error, arg2 string, arg3 ...interface{}) {
	fake.errorMutex.Lock()
	fake.errorArgsForCall = append(fake.errorArgsForCall, struct {
		arg1 error
		arg2 string
		arg3 []interface{}
	}{arg1, arg2, arg3})
	stub := fake.ErrorStub
	fake.recordInvocation("Error", []interface{}{arg1, arg2, arg3})
	fake.errorMutex.Unlock()
	if stub != nil {
		fake.ErrorStub(arg1, arg2, arg3...)
	}
}

func (fake *FakeLogger) ErrorCallCount() int {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return len(fake.errorArgsForCall)
}

func (fake *FakeLogger) ErrorCalls(stub func(error, string, ...interface{})) {
	fake.errorMutex.Lock()
	defer fake.errorMutex.Unlock()
	fake.ErrorStub = stub
}

func (fake *FakeLogger) ErrorArgsForCall(i int) (error, string, []interface{}) {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	argsForCall := fake.errorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeLogger) Info(arg1 string, arg2 ...interface{}) {
	fake.infoMutex.Lock()
	fake.infoArgsForCall = append(fake.infoArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	stub := fake.InfoStub
	fake.recordInvocation("Info", []interface{}{arg1, arg2})
	fake.infoMutex.Unlock()
	if stub != nil {
		fake.InfoStub(arg1, arg2...)
	}
}

func (fake *FakeLogger) InfoCallCount() int {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	return len(fake.infoArgsForCall)
}

func (fake *FakeLogger) InfoCalls(stub func(string, ...interface{})) {
	fake.infoMutex.Lock()
	defer fake.infoMutex.Unlock()
	fake.InfoStub = stub
}

func (fake *FakeLogger) InfoArgsForCall(i int) (string, []interface{}) {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	argsForCall := fake.infoArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeLogger) V(arg1 int) logr.Logger {
	fake.vMutex.Lock()
	ret, specificReturn := fake.vReturnsOnCall[len(fake.vArgsForCall)]
	fake.vArgsForCall = append(fake.vArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.VStub
	fakeReturns := fake.vReturns
	fake.recordInvocation("V", []interface{}{arg1})
	fake.vMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLogger) VCallCount() int {
	fake.vMutex.RLock()
	defer fake.vMutex.RUnlock()
	return len(fake.vArgsForCall)
}

func (fake *FakeLogger) VCalls(stub func(int) logr.Logger) {
	fake.vMutex.Lock()
	defer fake.vMutex.Unlock()
	fake.VStub = stub
}

func (fake *FakeLogger) VArgsForCall(i int) int {
	fake.vMutex.RLock()
	defer fake.vMutex.RUnlock()
	argsForCall := fake.vArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeLogger) VReturns(result1 logr.Logger) {
	fake.vMutex.Lock()
	defer fake.vMutex.Unlock()
	fake.VStub = nil
	fake.vReturns = struct {
		result1 logr.Logger
	}{result1}
}

func (fake *FakeLogger) VReturnsOnCall(i int, result1 logr.Logger) {
	fake.vMutex.Lock()
	defer fake.vMutex.Unlock()
	fake.VStub = nil
	if fake.vReturnsOnCall == nil {
		fake.vReturnsOnCall = make(map[int]struct {
			result1 logr.Logger
		})
	}
	fake.vReturnsOnCall[i] = struct {
		result1 logr.Logger
	}{result1}
}

func (fake *FakeLogger) WithName(arg1 string) logr.Logger {
	fake.withNameMutex.Lock()
	ret, specificReturn := fake.withNameReturnsOnCall[len(fake.withNameArgsForCall)]
	fake.withNameArgsForCall = append(fake.withNameArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.WithNameStub
	fakeReturns := fake.withNameReturns
	fake.recordInvocation("WithName", []interface{}{arg1})
	fake.withNameMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLogger) WithNameCallCount() int {
	fake.withNameMutex.RLock()
	defer fake.withNameMutex.RUnlock()
	return len(fake.withNameArgsForCall)
}

func (fake *FakeLogger) WithNameCalls(stub func(string) logr.Logger) {
	fake.withNameMutex.Lock()
	defer fake.withNameMutex.Unlock()
	fake.WithNameStub = stub
}

func (fake *FakeLogger) WithNameArgsForCall(i int) string {
	fake.withNameMutex.RLock()
	defer fake.withNameMutex.RUnlock()
	argsForCall := fake.withNameArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeLogger) WithNameReturns(result1 logr.Logger) {
	fake.withNameMutex.Lock()
	defer fake.withNameMutex.Unlock()
	fake.WithNameStub = nil
	fake.withNameReturns = struct {
		result1 logr.Logger
	}{result1}
}

func (fake *FakeLogger) WithNameReturnsOnCall(i int, result1 logr.Logger) {
	fake.withNameMutex.Lock()
	defer fake.withNameMutex.Unlock()
	fake.WithNameStub = nil
	if fake.withNameReturnsOnCall == nil {
		fake.withNameReturnsOnCall = make(map[int]struct {
			result1 logr.Logger
		})
	}
	fake.withNameReturnsOnCall[i] = struct {
		result1 logr.Logger
	}{result1}
}

func (fake *FakeLogger) WithValues(arg1 ...interface{}) logr.Logger {
	fake.withValuesMutex.Lock()
	ret, specificReturn := fake.withValuesReturnsOnCall[len(fake.withValuesArgsForCall)]
	fake.withValuesArgsForCall = append(fake.withValuesArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	stub := fake.WithValuesStub
	fakeReturns := fake.withValuesReturns
	fake.recordInvocation("WithValues", []interface{}{arg1})
	fake.withValuesMutex.Unlock()
	if stub != nil {
		return stub(arg1...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLogger) WithValuesCallCount() int {
	fake.withValuesMutex.RLock()
	defer fake.withValuesMutex.RUnlock()
	return len(fake.withValuesArgsForCall)
}

func (fake *FakeLogger) WithValuesCalls(stub func(...interface{}) logr.Logger) {
	fake.withValuesMutex.Lock()
	defer fake.withValuesMutex.Unlock()
	fake.WithValuesStub = stub
}

func (fake *FakeLogger) WithValuesArgsForCall(i int) []interface{} {
	fake.withValuesMutex.RLock()
	defer fake.withValuesMutex.RUnlock()
	argsForCall := fake.withValuesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeLogger) WithValuesReturns(result1 logr.Logger) {
	fake.withValuesMutex.Lock()
	defer fake.withValuesMutex.Unlock()
	fake.WithValuesStub = nil
	fake.withValuesReturns = struct {
		result1 logr.Logger
	}{result1}
}

func (fake *FakeLogger) WithValuesReturnsOnCall(i int, result1 logr.Logger) {
	fake.withValuesMutex.Lock()
	defer fake.withValuesMutex.Unlock()
	fake.WithValuesStub = nil
	if fake.withValuesReturnsOnCall == nil {
		fake.withValuesReturnsOnCall = make(map[int]struct {
			result1 logr.Logger
		})
	}
	fake.withValuesReturnsOnCall[i] = struct {
		result1 logr.Logger
	}{result1}
}

func (fake *FakeLogger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.enabledMutex.RLock()
	defer fake.enabledMutex.RUnlock()
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	fake.vMutex.RLock()
	defer fake.vMutex.RUnlock()
	fake.withNameMutex.RLock()
	defer fake.withNameMutex.RUnlock()
	fake.withValuesMutex.RLock()
	defer fake.withValuesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLogger) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ logr.Logger = new(FakeLogger)
