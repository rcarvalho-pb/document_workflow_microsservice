package main

import (
	"log"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("erro ao criar user:", r)
		}
	}()
	user, _ := (&model.UserBuilder{}).
		WithName("Ramon").
		WithLastName("Carvalho").
		WithPassword("123").
		WithEmail("ramon@email.com").
		WithRole(2).
		Build()
	log.Printf("%T - %v\n", user, user)
}
