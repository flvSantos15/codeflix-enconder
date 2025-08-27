package utils_test

import (
	"enconder/framework/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `{
			"id": "525b5fd9-700d-4feb-89c0-415a1e6e148c",
			"file_path": "test.mp4",
			"status": "ending"
	}`

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = `wes`
	err = utils.IsJson(json)
	require.Error(t, err)
}
