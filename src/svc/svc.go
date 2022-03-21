package svc

type Result struct {
	URI     string
	// service Service
}

type Service interface {
	Check(string) <-chan Result
}
