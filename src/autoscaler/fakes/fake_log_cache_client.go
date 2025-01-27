// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"sync"
	"time"

	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/eventgenerator/client"
	clienta "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/v9/rpc/loggregator_v2"
)

type FakeLogCacheClientReader struct {
	ReadStub        func(context.Context, string, time.Time, ...clienta.ReadOption) ([]*loggregator_v2.Envelope, error)
	readMutex       sync.RWMutex
	readArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 time.Time
		arg4 []clienta.ReadOption
	}
	readReturns struct {
		result1 []*loggregator_v2.Envelope
		result2 error
	}
	readReturnsOnCall map[int]struct {
		result1 []*loggregator_v2.Envelope
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLogCacheClientReader) Read(arg1 context.Context, arg2 string, arg3 time.Time, arg4 ...clienta.ReadOption) ([]*loggregator_v2.Envelope, error) {
	fake.readMutex.Lock()
	ret, specificReturn := fake.readReturnsOnCall[len(fake.readArgsForCall)]
	fake.readArgsForCall = append(fake.readArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 time.Time
		arg4 []clienta.ReadOption
	}{arg1, arg2, arg3, arg4})
	stub := fake.ReadStub
	fakeReturns := fake.readReturns
	fake.recordInvocation("Read", []interface{}{arg1, arg2, arg3, arg4})
	fake.readMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeLogCacheClientReader) ReadCallCount() int {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return len(fake.readArgsForCall)
}

func (fake *FakeLogCacheClientReader) ReadCalls(stub func(context.Context, string, time.Time, ...clienta.ReadOption) ([]*loggregator_v2.Envelope, error)) {
	fake.readMutex.Lock()
	defer fake.readMutex.Unlock()
	fake.ReadStub = stub
}

func (fake *FakeLogCacheClientReader) ReadArgsForCall(i int) (context.Context, string, time.Time, []clienta.ReadOption) {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	argsForCall := fake.readArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeLogCacheClientReader) ReadReturns(result1 []*loggregator_v2.Envelope, result2 error) {
	fake.readMutex.Lock()
	defer fake.readMutex.Unlock()
	fake.ReadStub = nil
	fake.readReturns = struct {
		result1 []*loggregator_v2.Envelope
		result2 error
	}{result1, result2}
}

func (fake *FakeLogCacheClientReader) ReadReturnsOnCall(i int, result1 []*loggregator_v2.Envelope, result2 error) {
	fake.readMutex.Lock()
	defer fake.readMutex.Unlock()
	fake.ReadStub = nil
	if fake.readReturnsOnCall == nil {
		fake.readReturnsOnCall = make(map[int]struct {
			result1 []*loggregator_v2.Envelope
			result2 error
		})
	}
	fake.readReturnsOnCall[i] = struct {
		result1 []*loggregator_v2.Envelope
		result2 error
	}{result1, result2}
}

func (fake *FakeLogCacheClientReader) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLogCacheClientReader) recordInvocation(key string, args []interface{}) {
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

var _ client.LogCacheClientReader = new(FakeLogCacheClientReader)
