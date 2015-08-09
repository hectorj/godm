package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVendorImport(t *testing.T) {
	err, ok := vendorImport(testRepoDir, "github.com/hectorj/gpm")
	assert.Nil(t, err)
	assert.True(t, ok)

	err, ok = vendorImport(testRepoDir, "github.com/hectorj/gpm")
	assert.Nil(t, err)
	assert.False(t, ok)

	// @TODO : more thorough checks of the resulting repo state
}
