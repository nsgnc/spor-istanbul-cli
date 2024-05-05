package service

import (
	"github.com/umtdemr/spor-istanbul-cli/internal/client"
	"github.com/umtdemr/spor-istanbul-cli/internal/parser"
	"strings"
)

type Service struct {
	client *client.Client
	parser *parser.Parser
}

func NewService() *Service {
	return &Service{
		client: client.NewClient(),
		parser: parser.NewParser(),
	}
}

func (s *Service) Login(id string, password string) bool {
	body := s.client.Login(id, password)
	title, ok := s.parser.GetTitle(body)

	if !ok {
		return false
	}

	return !strings.Contains(title, "Giriş Yap")

}

func (s *Service) GetSubscriptions() {
	body := s.client.GetSubscriptionsPage()
	s.parser.GetSubscriptions(body)
}
