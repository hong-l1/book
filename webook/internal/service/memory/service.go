package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
func (s *Service) SendSMS(ctx context.Context, template string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
