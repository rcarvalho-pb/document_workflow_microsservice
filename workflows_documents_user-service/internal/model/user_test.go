package model_test

import (
	"testing"
	"time"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	now := time.Now().Unix()
	ub := model.UserBuilder{}
	user, err := ub.
		WithName(" Ramon ").
		WithLastName("Almeida de Carvalho ").
		WithEmail("ramon@email.com").
		WithPassword("123").Build()
	assert.Nil(t, err)
	assert.Equal(t, "Ramon", user.Name)
	assert.Equal(t, "Almeida de Carvalho", user.LastName)
	assert.Equal(t, "ramon@email.com", user.Email)
	assert.True(t, now <= user.CreatedAt)
	assert.True(t, now <= user.UpdatedAt)
	assert.True(t, user.Active)
	ub = model.UserBuilder{}
	user2, err := ub.
		WithName(" Emilly ").
		WithLastName("Almeida de Carvalho ").
		WithEmail("ramon@email.com.br.").
		WithPassword("123").Build()
	assert.Error(t, err)
	assert.Nil(t, user2)
	lastUpdatedTime := user.UpdatedAt
	user.UpdateUserTime()
	assert.True(t, lastUpdatedTime <= user.UpdatedAt, "user update time funcion working")
}
