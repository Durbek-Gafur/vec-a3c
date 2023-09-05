package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type config struct {
	queue string
	cpu   string
	ram   string
}

func main() {
	// Define lists of possible values
	// backendQueueSizes := []string{"5", "6", "7", "8", "9", "10"}
	// backendCpus := []string{"0.5", "1", "1.5", "2", "2.5", "3"}
	// backendRams := []string{"512Mi", "1Gi", "1.5Gi", "2Gi", "2.5Gi", "3Gi"}
	config := []config{
		{
			queue: "5",
			cpu:   "4",
			ram:   "8Gi",
		},
		{
			queue: "10",
			cpu:   "6",
			ram:   "16Gi",
		},
		{
			queue: "15",
			cpu:   "8",
			ram:   "32Gi",
		},
	}

	for i := 1; i <= 10; i++ { // assuming we want 10 VENs
		venName := fmt.Sprintf("ven%d", i)
		backendQueueSize := config[(i-1)%len(config)].queue
		backendCpu := config[(i-1)%len(config)].cpu
		backendRam := config[(i-1)%len(config)].ram

		content, err := ioutil.ReadFile("ven_template.cue") // your CUE template
		if err != nil {
			panic(err)
		}

		newContent := strings.Replace(string(content), `ven1`, venName, -1)
		newContent = strings.Replace(newContent, `yetti`, backendQueueSize, -1)
		newContent = strings.Replace(newContent, `yarim`, backendCpu, -1)
		newContent = strings.Replace(newContent, `sakkizyuz`, backendRam, -1)

		err = ioutil.WriteFile(venName+".cue", []byte(newContent), 0644)
		if err != nil {
			panic(err)
		}

		// Execute `cue export` command
		cmd := exec.Command("cue", "export", "--out", "yaml", venName+".cue")
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error executing cue export for %s: %v\n", venName, err)
			continue
		}

		// Remove lines before the first "remove" occurrence
		yamlContent := string(output)
		removeIndex := strings.Index(yamlContent, "remove")
		if removeIndex != -1 {
			yamlContent = yamlContent[removeIndex:]

			// Replace all occurrences of "remove" with "---\n#"
			yamlContent = strings.ReplaceAll(yamlContent, "remove", "---\n#")
		} else {
			fmt.Printf("No 'remove' line found in the YAML output for %s\n", venName)
		}

		err = ioutil.WriteFile(venName+".yaml", []byte(yamlContent), 0644)
		if err != nil {
			fmt.Printf("Error writing YAML file for %s: %v\n", venName, err)
		}

		// Run `kubectl apply -f <filename>.yaml` command
		kubectlCmd := exec.Command("kubectl", "apply", "-f", venName+".yaml")
		kubectlCmd.Stdout = os.Stdout
		kubectlCmd.Stderr = os.Stderr

		err = kubectlCmd.Run()
		if err != nil {
			fmt.Printf("Error running kubectl apply for %s: %v\n", venName, err)
		}
	}
}
