package services

import "go-ns/src/models"

type Service interface {
	Check(string) <-chan models.Result
}
