package security_test

import (
	"testing"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/security"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "123"
	hashedPassword, err := security.Hash(password)
	assert.NoError(t, err, "should not return error to generate hashed password")
	assert.NotEmpty(t, hashedPassword, "hashed password should not be empty")

	ok := security.CheckPassword(password, hashedPassword)
	assert.True(t, ok, "passwords should be equal")
}

func TestCheckPassword_FailsWithWrongPassword(t *testing.T) {
	password := "123"
	wrongPassword := "234"

	hashedPassword, err := security.Hash(password)
	assert.NoError(t, err, "should not return error from hash password")

	ok := security.CheckPassword(wrongPassword, hashedPassword)
	assert.False(t, ok, "ok should be false to wrong password")
}

func TestHash_ShouldGenerateDifferentHashsForSamePassowrd(t *testing.T) {
	password := "123"
	hashedPassword1, _ := security.Hash(password)
	hashedPassword2, _ := security.Hash(password)

	assert.NotEqual(t, hashedPassword1, hashedPassword2, "should generate different hashes for same password")
}
