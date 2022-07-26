package data

import (
	"github.com/mcolomerc/kafkasr/proto/model"

	"github.com/brianvoe/gofakeit/v6"
)

func GetPerson() model.Person {
	return model.Person{
		Name:    gofakeit.Name(),
		Email:   gofakeit.Email(),
		Phone:   gofakeit.Phone(),
		Company: gofakeit.Company(),
		Job:     gofakeit.JobTitle(),
		Gender:  gofakeit.Gender(),
		Age:     gofakeit.Int32(),
	}
}

func GetPersonKey() string {
	return gofakeit.UUID()
}
