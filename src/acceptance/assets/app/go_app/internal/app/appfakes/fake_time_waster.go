// Code generated by counterfeiter. DO NOT EDIT.
package appfakes

import (
	"sync"
	"time"

	"code.cloudfoundry.org/app-autoscaler-release/src/acceptance/assets/app/go_app/internal/app"
)

type FakeTimeWaster struct {
	SleepStub        func(time.Duration)
	sleepMutex       sync.RWMutex
	sleepArgsForCall []struct {
		arg1 time.Duration
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTimeWaster) Sleep(arg1 time.Duration) {
	fake.sleepMutex.Lock()
	fake.sleepArgsForCall = append(fake.sleepArgsForCall, struct {
		arg1 time.Duration
	}{arg1})
	stub := fake.SleepStub
	fake.recordInvocation("Sleep", []interface{}{arg1})
	fake.sleepMutex.Unlock()
	if stub != nil {
		fake.SleepStub(arg1)
	}
}

func (fake *FakeTimeWaster) SleepCallCount() int {
	fake.sleepMutex.RLock()
	defer fake.sleepMutex.RUnlock()
	return len(fake.sleepArgsForCall)
}

func (fake *FakeTimeWaster) SleepCalls(stub func(time.Duration)) {
	fake.sleepMutex.Lock()
	defer fake.sleepMutex.Unlock()
	fake.SleepStub = stub
}

func (fake *FakeTimeWaster) SleepArgsForCall(i int) time.Duration {
	fake.sleepMutex.RLock()
	defer fake.sleepMutex.RUnlock()
	argsForCall := fake.sleepArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTimeWaster) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sleepMutex.RLock()
	defer fake.sleepMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTimeWaster) recordInvocation(key string, args []interface{}) {
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

var _ app.TimeWaster = new(FakeTimeWaster)
