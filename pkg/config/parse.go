package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	v1 "github.com/openshift-psap/ci-dashboard/api/matrix/v1"
	"sigs.k8s.io/yaml"
)

func ParseMatricesConfigFile(configFile string) (*v1.MatricesSpec, error) {
	var err error
	var configYaml []byte

	log.Printf("Reading configuration from: %s", configFile)

	// If the configFile is "-" use stdin for input
	if configFile == "-" {
		log.Println("Reading YAML from stdin...")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("Read line: %s", line) // Log each line read from stdin
			configYaml = append(configYaml, scanner.Bytes()...)
			configYaml = append(configYaml, '\n')
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from stdin: %v", err)
			return nil, fmt.Errorf("stdin read error: %v", err)
		}
	} else {
		log.Printf("Reading YAML file: %s...", configFile)
		configYaml, err = ioutil.ReadFile(configFile)
		if err != nil {
			log.Printf("Error reading file %s: %v", configFile, err)
			return nil, fmt.Errorf("read error: %v", err)
		}
		log.Printf("Successfully read %d bytes from file %s", len(configYaml), configFile)
	}

	// Print the raw YAML content before parsing
	log.Println("Raw YAML content:")
	log.Println(string(configYaml))

	log.Println("Parsing YAML...")
	var spec v1.MatricesSpec
	err = yaml.Unmarshal(configYaml, &spec)
	if err != nil {
		log.Printf("Error unmarshaling YAML: %v", err)
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	log.Println("YAML successfully parsed.")

	log.Println("Setting default values for ProwType if missing...")
	for _, matrix := range spec.Matrices {
		for testGroup, tests := range matrix.Tests {
			for i, test := range tests {
				if test.ProwType == "" {
					// Log the TestName
					log.Printf("Test '%s' in group %s has no ProwType, setting to 'periodic'", test.TestName, testGroup)
					test.ProwType = "periodic"
				}
				tests[i] = test
			}
			matrix.Tests[testGroup] = tests
		}
	}

	log.Println("Final parsed YAML:")
	str, err := yaml.Marshal(spec)
	if err != nil {
		log.Printf("Error marshaling YAML: %v", err)
	} else {
		log.Println("Parsed YAML output:\n", string(str))
	}

	log.Println("Returning parsed configuration...")
	return &spec, nil
}
