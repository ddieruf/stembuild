// Code generated by counterfeiter. DO NOT EDIT.
package stemcell_generatorfakes

import (
	io "io"
	sync "sync"

	stemcell_generator "github.com/cloudfoundry-incubator/stembuild/package_stemcell/stemcell_generator"
)

type FakeManifestGenerator struct {
	ManifestStub        func(io.Reader) (io.Reader, error)
	manifestMutex       sync.RWMutex
	manifestArgsForCall []struct {
		arg1 io.Reader
	}
	manifestReturns struct {
		result1 io.Reader
		result2 error
	}
	manifestReturnsOnCall map[int]struct {
		result1 io.Reader
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeManifestGenerator) Manifest(arg1 io.Reader) (io.Reader, error) {
	fake.manifestMutex.Lock()
	ret, specificReturn := fake.manifestReturnsOnCall[len(fake.manifestArgsForCall)]
	fake.manifestArgsForCall = append(fake.manifestArgsForCall, struct {
		arg1 io.Reader
	}{arg1})
	fake.recordInvocation("Manifest", []interface{}{arg1})
	fake.manifestMutex.Unlock()
	if fake.ManifestStub != nil {
		return fake.ManifestStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.manifestReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManifestGenerator) ManifestCallCount() int {
	fake.manifestMutex.RLock()
	defer fake.manifestMutex.RUnlock()
	return len(fake.manifestArgsForCall)
}

func (fake *FakeManifestGenerator) ManifestCalls(stub func(io.Reader) (io.Reader, error)) {
	fake.manifestMutex.Lock()
	defer fake.manifestMutex.Unlock()
	fake.ManifestStub = stub
}

func (fake *FakeManifestGenerator) ManifestArgsForCall(i int) io.Reader {
	fake.manifestMutex.RLock()
	defer fake.manifestMutex.RUnlock()
	argsForCall := fake.manifestArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManifestGenerator) ManifestReturns(result1 io.Reader, result2 error) {
	fake.manifestMutex.Lock()
	defer fake.manifestMutex.Unlock()
	fake.ManifestStub = nil
	fake.manifestReturns = struct {
		result1 io.Reader
		result2 error
	}{result1, result2}
}

func (fake *FakeManifestGenerator) ManifestReturnsOnCall(i int, result1 io.Reader, result2 error) {
	fake.manifestMutex.Lock()
	defer fake.manifestMutex.Unlock()
	fake.ManifestStub = nil
	if fake.manifestReturnsOnCall == nil {
		fake.manifestReturnsOnCall = make(map[int]struct {
			result1 io.Reader
			result2 error
		})
	}
	fake.manifestReturnsOnCall[i] = struct {
		result1 io.Reader
		result2 error
	}{result1, result2}
}

func (fake *FakeManifestGenerator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.manifestMutex.RLock()
	defer fake.manifestMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeManifestGenerator) recordInvocation(key string, args []interface{}) {
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

var _ stemcell_generator.ManifestGenerator = new(FakeManifestGenerator)
