package rspec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestGetRspec(t *testing.T) {
	// Set up the environment
	os.Setenv("CPUS", "0.75")
	os.Setenv("RAM", "512M")

	rspecService := NewService()
	// Run the function
	resources, err := rspecService.GetRspec()
	assert.NoError(t, err, "GetRspec failed")

	// Check the results
	assert.Equal(t, 0.75, resources.CPUs, "Unexpected CPUs value")
	assert.Equal(t, 512, resources.RAM, "Unexpected RAM value")

	// Clean up the environment
	os.Unsetenv("CPUS")
	os.Unsetenv("RAM")
}

func TestParseResourcesFromEnv(t *testing.T) {
	// Set up the environment
	os.Setenv("CPUS", "0.5")
	os.Setenv("RAM", "256M")
	rspecService := NewService()
	// Run the function
	resources, err := rspecService.ParseResourcesFromEnv()
	assert.NoError(t, err, "ParseResourcesFromEnv failed")

	// Check the results
	assert.Equal(t, 0.5, resources.CPUs, "Unexpected CPUs value")
	assert.Equal(t, 256, resources.RAM, "Unexpected RAM value")

	// Clean up the environment
	os.Unsetenv("CPUS")
	os.Unsetenv("RAM")
}


