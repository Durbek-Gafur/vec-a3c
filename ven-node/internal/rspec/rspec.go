package rspec

import (
	"os"
)

//go:generate mockgen -destination=mocks/rspec_mock.go -package=rspec_mock vec-node/internal/rspec Rspec

// Rspec handles operations on queue
type Rspec interface {
	GetRspec() (*Resources, error)
}

type Service struct {
}

func NewService() *Service {
	return &Service{
	}
}
func NewResource(cpu , ram, queue_size string)*Resources{
	return &Resources{
		CPUs: cpu,
		RAM: ram,
		MAX_QUEUE: queue_size,
	}
}
type Resources struct {
	CPUs string
	RAM  string 
	MAX_QUEUE string
}

func (s *Service)ParseResourcesFromEnv() (*Resources, error) {
	resources := &Resources{}

	resources.CPUs = os.Getenv("CPUS")
	resources.RAM = os.Getenv("RAM")
	resources.RAM = os.Getenv("QUEUE_SIZE")

	return resources, nil
}

func (s *Service)GetRspec() (*Resources, error) {
	return s.ParseResourcesFromEnv()
}