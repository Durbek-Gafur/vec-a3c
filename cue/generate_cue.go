package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func main() {
	// Define lists of possible values
	backendQueueSizes := []string{"5", "6", "7", "8", "9", "10"}
	backendCpus := []string{"0.5", "1", "1.5", "2", "2.5", "3"}
	backendRams := []string{"512Mi", "1Gi", "1.5Gi", "2Gi", "2.5Gi", "3Gi"}

	for i := 1; i <= 10; i++ { // assuming we want 10 VENs
		venName := fmt.Sprintf("ven%d", i)
		backendQueueSize := backendQueueSizes[(i-1)%len(backendQueueSizes)]
		backendCpu := backendCpus[(i-1)%len(backendCpus)]
		backendRam := backendRams[(i-1)%len(backendRams)]

		content, err := ioutil.ReadFile("ven_template.cue") // your CUE template
		if err != nil {
			panic(err)
		}

		newContent := strings.Replace(string(content), `ven1`,venName, -1)
		newContent = strings.Replace(newContent, `7`, backendQueueSize, -1)
		newContent = strings.Replace(newContent, `0.5`, backendCpu, -1)
		newContent = strings.Replace(newContent, `800Mi`, backendRam, -1)

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
	}
}
