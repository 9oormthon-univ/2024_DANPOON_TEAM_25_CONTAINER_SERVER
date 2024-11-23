package gitclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifyRepository(t *testing.T) {
	// Setup environment variables
	gitClient := NewGitClient()
	// Create manifests using test key
	err := gitClient.ModifyRepository("1", "1")
	assert.NoError(t, err, "ModifyRepository should not return an error")
}
