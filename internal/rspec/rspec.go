package rspec

import (
	"os"
	"strconv"
	"strings"
)

type Resources struct {
	CPUs float64
	RAM  int // in MB
}

func ParseResourcesFromEnv() (*Resources, error) {
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

func GetRspec() (*Resources, error) {
	return ParseResourcesFromEnv()
}