// Command fixspec fixes OpenAPI 3.1 type arrays for ogen compatibility.
// OpenAPI 3.1 allows type: [string, null] but ogen expects type: string.
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	content, err := os.ReadFile("openapi-spec.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read spec: %v\n", err)
		os.Exit(1)
	}

	var spec map[string]any
	if err := yaml.Unmarshal(content, &spec); err != nil {
		fmt.Fprintf(os.Stderr, "parse yaml: %v\n", err)
		os.Exit(1)
	}

	fixTypes(spec)

	out, err := yaml.Marshal(spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal yaml: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("openapi-spec.yaml", out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write spec: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Fixed OpenAPI spec for ogen compatibility")
}

// fixTypes recursively walks the spec and converts type arrays to single values.
// OpenAPI 3.1: type: [string, null] -> OpenAPI 3.0: type: string
func fixTypes(v any) {
	switch val := v.(type) {
	case map[string]any:
		if typeVal, ok := val["type"]; ok {
			if types, ok := typeVal.([]any); ok {
				for _, t := range types {
					if s, ok := t.(string); ok && s != "null" {
						val["type"] = s
						break
					}
				}
			}
		}
		for _, child := range val {
			fixTypes(child)
		}
	case []any:
		for _, item := range val {
			fixTypes(item)
		}
	}
}
