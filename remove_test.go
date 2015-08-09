package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveImport(t *testing.T) {
	err := removeImport(testRepoDir, "github.com/hectorj/gpm-test-repo", false)
	assert.Nil(t, err)

	err = removeImport(testRepoDir, "github.com/hectorj/gpm-test-repo", false)
	assert.NotNil(t, err)

	// @TODO : more thorough checks of the resulting repo state
}
