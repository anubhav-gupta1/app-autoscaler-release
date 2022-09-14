// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/cloudfoundry/app-autoscaler-release/db"
	"github.com/cloudfoundry/app-autoscaler-release/models"
)

type FakeStoredProcedureDB struct {
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	closeReturns struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	CreateCredentialsStub        func(models.CredentialsOptions) (*models.Credential, error)
	createCredentialsMutex       sync.RWMutex
	createCredentialsArgsForCall []struct {
		arg1 models.CredentialsOptions
	}
	createCredentialsReturns struct {
		result1 *models.Credential
		result2 error
	}
	createCredentialsReturnsOnCall map[int]struct {
		result1 *models.Credential
		result2 error
	}
	DeleteAllInstanceCredentialsStub        func(string) error
	deleteAllInstanceCredentialsMutex       sync.RWMutex
	deleteAllInstanceCredentialsArgsForCall []struct {
		arg1 string
	}
	deleteAllInstanceCredentialsReturns struct {
		result1 error
	}
	deleteAllInstanceCredentialsReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteCredentialsStub        func(models.CredentialsOptions) error
	deleteCredentialsMutex       sync.RWMutex
	deleteCredentialsArgsForCall []struct {
		arg1 models.CredentialsOptions
	}
	deleteCredentialsReturns struct {
		result1 error
	}
	deleteCredentialsReturnsOnCall map[int]struct {
		result1 error
	}
	PingStub        func() error
	pingMutex       sync.RWMutex
	pingArgsForCall []struct {
	}
	pingReturns struct {
		result1 error
	}
	pingReturnsOnCall map[int]struct {
		result1 error
	}
	ValidateCredentialsStub        func(models.Credential) (*models.CredentialsOptions, error)
	validateCredentialsMutex       sync.RWMutex
	validateCredentialsArgsForCall []struct {
		arg1 models.Credential
	}
	validateCredentialsReturns struct {
		result1 *models.CredentialsOptions
		result2 error
	}
	validateCredentialsReturnsOnCall map[int]struct {
		result1 *models.CredentialsOptions
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStoredProcedureDB) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	stub := fake.CloseStub
	fakeReturns := fake.closeReturns
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStoredProcedureDB) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeStoredProcedureDB) CloseCalls(stub func() error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeStoredProcedureDB) CloseReturns(result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoredProcedureDB) CloseReturnsOnCall(i int, result1 error) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoredProcedureDB) CreateCredentials(arg1 models.CredentialsOptions) (*models.Credential, error) {
	fake.createCredentialsMutex.Lock()
	ret, specificReturn := fake.createCredentialsReturnsOnCall[len(fake.createCredentialsArgsForCall)]
	fake.createCredentialsArgsForCall = append(fake.createCredentialsArgsForCall, struct {
		arg1 models.CredentialsOptions
	}{arg1})
	stub := fake.CreateCredentialsStub
	fakeReturns := fake.createCredentialsReturns
	fake.recordInvocation("CreateCredentials", []interface{}{arg1})
	fake.createCredentialsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStoredProcedureDB) CreateCredentialsCallCount() int {
	fake.createCredentialsMutex.RLock()
	defer fake.createCredentialsMutex.RUnlock()
	return len(fake.createCredentialsArgsForCall)
}

func (fake *FakeStoredProcedureDB) CreateCredentialsCalls(stub func(models.CredentialsOptions) (*models.Credential, error)) {
	fake.createCredentialsMutex.Lock()
	defer fake.createCredentialsMutex.Unlock()
	fake.CreateCredentialsStub = stub
}

func (fake *FakeStoredProcedureDB) CreateCredentialsArgsForCall(i int) models.CredentialsOptions {
	fake.createCredentialsMutex.RLock()
	defer fake.createCredentialsMutex.RUnlock()
	argsForCall := fake.createCredentialsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStoredProcedureDB) CreateCredentialsReturns(result1 *models.Credential, result2 error) {
	fake.createCredentialsMutex.Lock()
	defer fake.createCredentialsMutex.Unlock()
	fake.CreateCredentialsStub = nil
	fake.createCredentialsReturns = struct {
		result1 *models.Credential
		result2 error
	}{result1, result2}
}

func (fake *FakeStoredProcedureDB) CreateCredentialsReturnsOnCall(i int, result1 *models.Credential, result2 error) {
	fake.createCredentialsMutex.Lock()
	defer fake.createCredentialsMutex.Unlock()
	fake.CreateCredentialsStub = nil
	if fake.createCredentialsReturnsOnCall == nil {
		fake.createCredentialsReturnsOnCall = make(map[int]struct {
			result1 *models.Credential
			result2 error
		})
	}
	fake.createCredentialsReturnsOnCall[i] = struct {
		result1 *models.Credential
		result2 error
	}{result1, result2}
}

func (fake *FakeStoredProcedureDB) DeleteAllInstanceCredentials(arg1 string) error {
	fake.deleteAllInstanceCredentialsMutex.Lock()
	ret, specificReturn := fake.deleteAllInstanceCredentialsReturnsOnCall[len(fake.deleteAllInstanceCredentialsArgsForCall)]
	fake.deleteAllInstanceCredentialsArgsForCall = append(fake.deleteAllInstanceCredentialsArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.DeleteAllInstanceCredentialsStub
	fakeReturns := fake.deleteAllInstanceCredentialsReturns
	fake.recordInvocation("DeleteAllInstanceCredentials", []interface{}{arg1})
	fake.deleteAllInstanceCredentialsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStoredProcedureDB) DeleteAllInstanceCredentialsCallCount() int {
	fake.deleteAllInstanceCredentialsMutex.RLock()
	defer fake.deleteAllInstanceCredentialsMutex.RUnlock()
	return len(fake.deleteAllInstanceCredentialsArgsForCall)
}

func (fake *FakeStoredProcedureDB) DeleteAllInstanceCredentialsCalls(stub func(string) error) {
	fake.deleteAllInstanceCredentialsMutex.Lock()
	defer fake.deleteAllInstanceCredentialsMutex.Unlock()
	fake.DeleteAllInstanceCredentialsStub = stub
}

func (fake *FakeStoredProcedureDB) DeleteAllInstanceCredentialsArgsForCall(i int) string {
	fake.deleteAllInstanceCredentialsMutex.RLock()
	defer fake.deleteAllInstanceCredentialsMutex.RUnlock()
	argsForCall := fake.deleteAllInstanceCredentialsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStoredProcedureDB) DeleteAllInstanceCredentialsReturns(result1 error) {
	fake.deleteAllInstanceCredentialsMutex.Lock()
	defer fake.deleteAllInstanceCredentialsMutex.Unlock()
	fake.DeleteAllInstanceCredentialsStub = nil
	fake.deleteAllInstanceCredentialsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoredProcedureDB) DeleteAllInstanceCredentialsReturnsOnCall(i int, result1 error) {
	fake.deleteAllInstanceCredentialsMutex.Lock()
	defer fake.deleteAllInstanceCredentialsMutex.Unlock()
	fake.DeleteAllInstanceCredentialsStub = nil
	if fake.deleteAllInstanceCredentialsReturnsOnCall == nil {
		fake.deleteAllInstanceCredentialsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteAllInstanceCredentialsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoredProcedureDB) DeleteCredentials(arg1 models.CredentialsOptions) error {
	fake.deleteCredentialsMutex.Lock()
	ret, specificReturn := fake.deleteCredentialsReturnsOnCall[len(fake.deleteCredentialsArgsForCall)]
	fake.deleteCredentialsArgsForCall = append(fake.deleteCredentialsArgsForCall, struct {
		arg1 models.CredentialsOptions
	}{arg1})
	stub := fake.DeleteCredentialsStub
	fakeReturns := fake.deleteCredentialsReturns
	fake.recordInvocation("DeleteCredentials", []interface{}{arg1})
	fake.deleteCredentialsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStoredProcedureDB) DeleteCredentialsCallCount() int {
	fake.deleteCredentialsMutex.RLock()
	defer fake.deleteCredentialsMutex.RUnlock()
	return len(fake.deleteCredentialsArgsForCall)
}

func (fake *FakeStoredProcedureDB) DeleteCredentialsCalls(stub func(models.CredentialsOptions) error) {
	fake.deleteCredentialsMutex.Lock()
	defer fake.deleteCredentialsMutex.Unlock()
	fake.DeleteCredentialsStub = stub
}

func (fake *FakeStoredProcedureDB) DeleteCredentialsArgsForCall(i int) models.CredentialsOptions {
	fake.deleteCredentialsMutex.RLock()
	defer fake.deleteCredentialsMutex.RUnlock()
	argsForCall := fake.deleteCredentialsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStoredProcedureDB) DeleteCredentialsReturns(result1 error) {
	fake.deleteCredentialsMutex.Lock()
	defer fake.deleteCredentialsMutex.Unlock()
	fake.DeleteCredentialsStub = nil
	fake.deleteCredentialsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoredProcedureDB) DeleteCredentialsReturnsOnCall(i int, result1 error) {
	fake.deleteCredentialsMutex.Lock()
	defer fake.deleteCredentialsMutex.Unlock()
	fake.DeleteCredentialsStub = nil
	if fake.deleteCredentialsReturnsOnCall == nil {
		fake.deleteCredentialsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteCredentialsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoredProcedureDB) Ping() error {
	fake.pingMutex.Lock()
	ret, specificReturn := fake.pingReturnsOnCall[len(fake.pingArgsForCall)]
	fake.pingArgsForCall = append(fake.pingArgsForCall, struct {
	}{})
	stub := fake.PingStub
	fakeReturns := fake.pingReturns
	fake.recordInvocation("Ping", []interface{}{})
	fake.pingMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStoredProcedureDB) PingCallCount() int {
	fake.pingMutex.RLock()
	defer fake.pingMutex.RUnlock()
	return len(fake.pingArgsForCall)
}

func (fake *FakeStoredProcedureDB) PingCalls(stub func() error) {
	fake.pingMutex.Lock()
	defer fake.pingMutex.Unlock()
	fake.PingStub = stub
}

func (fake *FakeStoredProcedureDB) PingReturns(result1 error) {
	fake.pingMutex.Lock()
	defer fake.pingMutex.Unlock()
	fake.PingStub = nil
	fake.pingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoredProcedureDB) PingReturnsOnCall(i int, result1 error) {
	fake.pingMutex.Lock()
	defer fake.pingMutex.Unlock()
	fake.PingStub = nil
	if fake.pingReturnsOnCall == nil {
		fake.pingReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.pingReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStoredProcedureDB) ValidateCredentials(arg1 models.Credential) (*models.CredentialsOptions, error) {
	fake.validateCredentialsMutex.Lock()
	ret, specificReturn := fake.validateCredentialsReturnsOnCall[len(fake.validateCredentialsArgsForCall)]
	fake.validateCredentialsArgsForCall = append(fake.validateCredentialsArgsForCall, struct {
		arg1 models.Credential
	}{arg1})
	stub := fake.ValidateCredentialsStub
	fakeReturns := fake.validateCredentialsReturns
	fake.recordInvocation("ValidateCredentials", []interface{}{arg1})
	fake.validateCredentialsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStoredProcedureDB) ValidateCredentialsCallCount() int {
	fake.validateCredentialsMutex.RLock()
	defer fake.validateCredentialsMutex.RUnlock()
	return len(fake.validateCredentialsArgsForCall)
}

func (fake *FakeStoredProcedureDB) ValidateCredentialsCalls(stub func(models.Credential) (*models.CredentialsOptions, error)) {
	fake.validateCredentialsMutex.Lock()
	defer fake.validateCredentialsMutex.Unlock()
	fake.ValidateCredentialsStub = stub
}

func (fake *FakeStoredProcedureDB) ValidateCredentialsArgsForCall(i int) models.Credential {
	fake.validateCredentialsMutex.RLock()
	defer fake.validateCredentialsMutex.RUnlock()
	argsForCall := fake.validateCredentialsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStoredProcedureDB) ValidateCredentialsReturns(result1 *models.CredentialsOptions, result2 error) {
	fake.validateCredentialsMutex.Lock()
	defer fake.validateCredentialsMutex.Unlock()
	fake.ValidateCredentialsStub = nil
	fake.validateCredentialsReturns = struct {
		result1 *models.CredentialsOptions
		result2 error
	}{result1, result2}
}

func (fake *FakeStoredProcedureDB) ValidateCredentialsReturnsOnCall(i int, result1 *models.CredentialsOptions, result2 error) {
	fake.validateCredentialsMutex.Lock()
	defer fake.validateCredentialsMutex.Unlock()
	fake.ValidateCredentialsStub = nil
	if fake.validateCredentialsReturnsOnCall == nil {
		fake.validateCredentialsReturnsOnCall = make(map[int]struct {
			result1 *models.CredentialsOptions
			result2 error
		})
	}
	fake.validateCredentialsReturnsOnCall[i] = struct {
		result1 *models.CredentialsOptions
		result2 error
	}{result1, result2}
}

func (fake *FakeStoredProcedureDB) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.createCredentialsMutex.RLock()
	defer fake.createCredentialsMutex.RUnlock()
	fake.deleteAllInstanceCredentialsMutex.RLock()
	defer fake.deleteAllInstanceCredentialsMutex.RUnlock()
	fake.deleteCredentialsMutex.RLock()
	defer fake.deleteCredentialsMutex.RUnlock()
	fake.pingMutex.RLock()
	defer fake.pingMutex.RUnlock()
	fake.validateCredentialsMutex.RLock()
	defer fake.validateCredentialsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeStoredProcedureDB) recordInvocation(key string, args []interface{}) {
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

var _ db.StoredProcedureDB = new(FakeStoredProcedureDB)
