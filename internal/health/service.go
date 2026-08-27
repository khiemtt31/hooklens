package health

import "time"

type Report struct {
	Status        string    `json:"status"`
	Service       string    `json:"service"`
	Version       string    `json:"version"`
	UptimeSeconds int64     `json:"uptimeSeconds"`
	Timestamp     time.Time `json:"timestamp"`
}

type Service struct {
	serviceName string
	version     string
	startedAt   time.Time
}

func NewService(serviceName string, version string) *Service {
	return &Service{
		serviceName: serviceName,
		version:     version,
		startedAt:   time.Now(),
	}
}

func (s *Service) Check() Report {
	now := time.Now().UTC()

	return Report{
		Status:        "ok",
		Service:       s.serviceName,
		Version:       s.version,
		UptimeSeconds: int64(time.Since(s.startedAt).Seconds()),
		Timestamp:     now,
	}
}
