package model_test

import (
	"testing"
	"time"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
)

func TestCreateUser(t *testing.T) {
	now := time.Now().Unix()
	user := model.UserBuilder{}
	user.
		WithName(" Ramon ").
		WithLastName("Almeida de Carvalho ").
		WithEmail("ramon@email.com")
}
