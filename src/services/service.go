package services

type Service interface {
	Check(string) <-chan string
}
