package rspec

import (
	"os"
	"strconv"
	"strings"
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
func NewResource(cpu float64, ram int)*Resources{
	return &Resources{
		CPUs: cpu,
		RAM: ram,
	}
}
type Resources struct {
	CPUs float64
	RAM  int // in MB
}

func (s *Service)ParseResourcesFromEnv() (*Resources, error) {
	resources := &Resources{}

	cpuStr := os.Getenv("CPUS")
	cpus, err := strconv.ParseFloat(cpuStr, 64)
	if err != nil {
		return nil, err
	}
	resources.CPUs = cpus

	ramStr := os.Getenv("RAM")
	ram, err := strconv.Atoi(strings.TrimSuffix(ramStr, "M"))
	if err != nil {
		return nil, err
	}
	resources.RAM = ram

	return resources, nil
}

func (s *Service)GetRspec() (*Resources, error) {
	return s.ParseResourcesFromEnv()
}