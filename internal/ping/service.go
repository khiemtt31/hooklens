package ping

type Request struct {
	Pong string `json:"pong"`
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Ping(value string) string {
	return value
}
