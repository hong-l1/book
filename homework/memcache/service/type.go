package service

type Service interface {
	SendSMS(template string, numbers []string, number string) error
}
