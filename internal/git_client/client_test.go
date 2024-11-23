package gitclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifyRepository(t *testing.T) {
	// Setup environment variables
	gitClient := NewGitClient()
	testKey := "dXNlcjFjb3Vyc2Ux"

	// Create manifests using test key
	err := gitClient.ModifyRepository(testKey)
	assert.NoError(t, err, "ModifyRepository should not return an error")

}