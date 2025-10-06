package model_test

import (
	"testing"
	"time"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	now := time.Now().Unix()
	user, err := model.NewUser("Ramon", "Almeida de Carvalho", "ramon@email.com", "123")
	assert.Nil(t, err)
	assert.Equal(t, "Ramon", user.Name)
	assert.Equal(t, "Almeida de Carvalho", user.LastName)
	assert.Equal(t, "ramon@email.com", user.Email)
	assert.True(t, now <= user.CreatedAt)
	assert.True(t, now <= user.UpdatedAt)
	assert.True(t, user.Active)
	user2, err := model.NewUser("Emilly", "Almeida de Carvalho", "ramon@email.com.", "123")
	assert.Nil(t, user2)
	assert.Error(t, err)
	lastUpdatedTime := user.UpdatedAt
	user.DeactivateUser()
	assert.True(t, lastUpdatedTime <= user.UpdatedAt, "user update time funcion working")
}
