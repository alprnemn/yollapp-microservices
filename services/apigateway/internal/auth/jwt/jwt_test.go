package jwt

import (
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTokenValidator_Validate(t *testing.T) {

	validator := NewTokenValidator(
		55*time.Second,
		"hf5y4f8zhyfd4zxxv3d",
		"yollapp.com",
	)

	tkn, err := validator.Validate("")
	require.NoError(t, err)

	id, err := tkn.Claims.GetSubject()
	assert.Equal(t, id, "c8d00537-4ee3-4ee2-9472-a4a36ca44f05")
}
