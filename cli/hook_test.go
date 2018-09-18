// Copyright (c) 2018 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"context"
	"os"
	"testing"

	. "github.com/kata-containers/runtime/virtcontainers/pkg/mock"
	"github.com/kata-containers/runtime/virtcontainers/pkg/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/assert"
)

// Important to keep these values in sync with hook test binary
var testKeyHook = "test-key"
var testContainerIDHook = "test-container-id"
var testControllerIDHook = "test-controller-id"
var testBinHookPath = "/usr/bin/virtcontainers/bin/test/hook"
var testBundlePath = "/test/bundle"

func getMockHookBinPath() string {
	if DefaultMockHookBinPath == "" {
		return testBinHookPath
	}

	return DefaultMockHookBinPath
}

func createHook(timeout int) specs.Hook {
	to := &timeout
	if timeout == 0 {
		to = nil
	}

	return specs.Hook{
		Path:    getMockHookBinPath(),
		Args:    []string{testKeyHook, testContainerIDHook, testControllerIDHook},
		Env:     os.Environ(),
		Timeout: to,
	}
}

func createWrongHook() specs.Hook {
	return specs.Hook{
		Path: getMockHookBinPath(),
		Args: []string{"wrong-args"},
		Env:  os.Environ(),
	}
}

func TestRunHook(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip(testDisabledNeedNonRoot)
	}

	assert := assert.New(t)

	ctx := context.Background()

	// Run with timeout 0
	hook := createHook(0)
	err := runHook(ctx, hook, testSandboxID, testBundlePath)
	assert.NoError(err)

	// Run with timeout 1
	hook = createHook(1)
	err = runHook(ctx, hook, testSandboxID, testBundlePath)
	assert.NoError(err)

	// Run timeout failure
	hook = createHook(1)
	hook.Args = append(hook.Args, "2")
	err = runHook(ctx, hook, testSandboxID, testBundlePath)
	assert.Error(err)

	// Failure due to wrong hook
	hook = createWrongHook()
	err = runHook(ctx, hook, testSandboxID, testBundlePath)
	assert.Error(err)
}

func TestPreStartHooks(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip(testDisabledNeedNonRoot)
	}

	assert := assert.New(t)

	ctx := context.Background()

	// Hooks field is nil
	spec := oci.CompatOCISpec{}
	err := preStartHooks(ctx, spec, "", "")
	assert.NoError(err)

	// Hooks list is empty
	spec = oci.CompatOCISpec{
		Spec: specs.Spec{
			Hooks: &specs.Hooks{},
		},
	}
	err = preStartHooks(ctx, spec, "", "")
	assert.NoError(err)

	// Run with timeout 0
	hook := createHook(0)
	spec = oci.CompatOCISpec{
		Spec: specs.Spec{
			Hooks: &specs.Hooks{
				Prestart: []specs.Hook{hook},
			},
		},
	}
	err = preStartHooks(ctx, spec, testSandboxID, testBundlePath)
	assert.NoError(err)

	// Failure due to wrong hook
	hook = createWrongHook()
	spec = oci.CompatOCISpec{
		Spec: specs.Spec{
			Hooks: &specs.Hooks{
				Prestart: []specs.Hook{hook},
			},
		},
	}
	err = preStartHooks(ctx, spec, testSandboxID, testBundlePath)
	assert.Error(err)
}

func TestPostStartHooks(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip(testDisabledNeedNonRoot)
	}

	assert := assert.New(t)

	ctx := context.Background()

	// Hooks field is nil
	spec := oci.CompatOCISpec{}
	err := postStartHooks(ctx, spec, "", "")
	assert.NoError(err)

	// Hooks list is empty
	spec = oci.CompatOCISpec{
		Spec: specs.Spec{
			Hooks: &specs.Hooks{},
		},
	}
	err = postStartHooks(ctx, spec, "", "")
	assert.NoError(err)

	// Run with timeout 0
	hook := createHook(0)
	spec = oci.CompatOCISpec{
		Spec: specs.Spec{
			Hooks: &specs.Hooks{
				Poststart: []specs.Hook{hook},
			},
		},
	}
	err = postStartHooks(ctx, spec, testSandboxID, testBundlePath)
	assert.NoError(err)

	// Failure due to wrong hook
	hook = createWrongHook()
	spec = oci.CompatOCISpec{
		Spec: specs.Spec{
			Hooks: &specs.Hooks{
				Poststart: []specs.Hook{hook},
			},
		},
	}
	err = postStartHooks(ctx, spec, testSandboxID, testBundlePath)
	assert.Error(err)
}

func TestPostStopHooks(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip(testDisabledNeedNonRoot)
	}

	assert := assert.New(t)

	ctx := context.Background()

	// Hooks field is nil
	spec := oci.CompatOCISpec{}
	err := postStopHooks(ctx, spec, "", "")
	assert.NoError(err)

	// Hooks list is empty
	spec = oci.CompatOCISpec{
		Spec: specs.Spec{
			Hooks: &specs.Hooks{},
		},
	}
	err = postStopHooks(ctx, spec, "", "")
	assert.NoError(err)

	// Run with timeout 0
	hook := createHook(0)
	spec = oci.CompatOCISpec{
		Spec: specs.Spec{
			Hooks: &specs.Hooks{
				Poststop: []specs.Hook{hook},
			},
		},
	}
	err = postStopHooks(ctx, spec, testSandboxID, testBundlePath)
	assert.NoError(err)

	// Failure due to wrong hook
	hook = createWrongHook()
	spec = oci.CompatOCISpec{
		Spec: specs.Spec{
			Hooks: &specs.Hooks{
				Poststop: []specs.Hook{hook},
			},
		},
	}
	err = postStopHooks(ctx, spec, testSandboxID, testBundlePath)
	assert.Error(err)
}