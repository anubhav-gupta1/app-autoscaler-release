// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"github.com/cloudfoundry/app-autoscaler-release/scalingengine/schedule"
	"sync"
)

type FakeActiveScheduleSychronizer struct {
	SyncStub         func()
	syncMutex        sync.RWMutex
	syncArgsForCall  []struct{}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeActiveScheduleSychronizer) Sync() {
	fake.syncMutex.Lock()
	fake.syncArgsForCall = append(fake.syncArgsForCall, struct{}{})
	fake.recordInvocation("Sync", []interface{}{})
	fake.syncMutex.Unlock()
	if fake.SyncStub != nil {
		fake.SyncStub()
	}
}

func (fake *FakeActiveScheduleSychronizer) SyncCallCount() int {
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	return len(fake.syncArgsForCall)
}

func (fake *FakeActiveScheduleSychronizer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.syncMutex.RLock()
	defer fake.syncMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeActiveScheduleSychronizer) recordInvocation(key string, args []interface{}) {
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

var _ schedule.ActiveScheduleSychronizer = new(FakeActiveScheduleSychronizer)
