package main

import (
	"log"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/dto"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("erro ao criar user:", r)
		}
	}()
	user, _ := model.NewUser("Ramon", "Carvalho", "ramon@email.com", "123", model.EMPLOYEE)
	log.Printf("%T - %v\n", user, user)
	userDTO := dto.UserDTO{
		Name:     "Ramon",
		LastName: "Carvalho",
		Email:    "ramon@email.com",
		Password: "123",
	}
	user2 := userDTO.ToUserModel()
	log.Println(user2)
}
