package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
)

type RenderPlan struct {
	Resources []cue.Value
	err       error
}

var (
	r           []RenderPlan
	currentPlan *RenderPlan
)

func getCueFiles(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if !info.IsDir() && strings.HasSuffix(path, ".cue") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// Find kubernetes manifests
func parseResource(v cue.Value) bool {
	plan := RenderPlan{}
	resource := v.LookupPath(cue.ParsePath("apiVersion"))
	if resource.Exists() {
		compileError := v.Validate(
			cue.All(),
			cue.Attributes(true),
			cue.Definitions(true),
			cue.InlineImports(true),
			cue.Concrete(true),
			cue.Final(),
			cue.DisallowCycles(true),
			cue.Hidden(true),
			cue.Optional(true),
		)
		if compileError != nil {
			plan.err = compileError
			r = append(r, plan)
			return false
		}
		r = append(r, plan)
		currentPlan = &r[len(r)-1]
		currentPlan.Resources = append(currentPlan.Resources, v)
		return false
	}
	return true
}

func parseCue(files []string) error {
	ctx := cuecontext.New()
	instances := load.Instances(files, nil)
	values, err := ctx.BuildInstances(instances)
	if err != nil {
		panic(err)
	}

	// Parse the values
	for _, value := range values {
		_, err := yaml.Encode(value)
		if err != nil {
			return err
		}
		value.Walk(parseResource, nil)
	}
	return nil
}

func renderResultsDir(outputDir string) {
	if !strings.Contains(outputDir, "_rendered") {
		fmt.Printf("error: outputDir must contain '_rendered' in its path\n")
		return
	}
	err := os.RemoveAll(outputDir)
	if err != nil {
		fmt.Printf("error removing output directory: %v\n", err)
		return
	}
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("error creating output directory: %v\n", err)
		return
	}
	for _, plan := range r {
		if plan.err != nil {
			fmt.Printf("error: %v\n", plan.err)
			continue
		}

		for _, v := range plan.Resources {
			kind, err := v.LookupPath(cue.ParsePath("kind")).String()
			if err != nil {
				fmt.Printf("error getting resource kind: %v\n", err)
				continue
			}
			namespace, err := v.LookupPath(cue.ParsePath("metadata.namespace")).String()
			if err != nil {
				namespace = ""
			}
			name, err := v.LookupPath(cue.ParsePath("metadata.name")).String()
			if err != nil {
				fmt.Printf("error getting resource name: %v\n", err)
				continue
			}
			filename := ""
			switch namespace {
			case "":
				filename = filepath.Join(outputDir, strings.ToLower(kind+"-"+name+".yaml"))
			default:
				filename = filepath.Join(outputDir, strings.ToLower(namespace+"-"+kind+"-"+name+".yaml"))
			}

			c, err := yaml.Encode(v)
			if err != nil {
				fmt.Printf("error encoding resource to YAML: %v\n", err)
				continue
			}

			// Write YAML to file
			if err := os.WriteFile(filename, c, 0644); err != nil {
				fmt.Printf("error writing file %s: %v\n", filename, err)
				continue
			}

			fmt.Printf("Created: %s\n", filename)
		}
	}
}

func renderResultsStdout() {
	for _, plan := range r {
		if plan.err != nil {
			fmt.Printf("error: %v\n", plan.err)
			continue
		}

		for _, v := range plan.Resources {
			c, err := yaml.Encode(v)
			if err != nil {
				fmt.Printf("error encoding resource to YAML: %v\n", err)
				continue
			}

			// Print YAML to stdout
			fmt.Printf("---\n%s", c)
		}
	}
}

// Find cue files in path and render to multi-document yaml (stdout) or to _rendered/ directory.
// Filename will be _rendered/[namespace-]kind-name.yaml
func main() {
	out := flag.String("out", "stdout", "Output destination: 'stdout' or 'files'")
	flag.Parse()

	args := flag.Args()
	path := ""
	switch len(args) {
	case 0:
		path = "."
	case 1:
		path = args[0]
	default:
		panic("too many arguments")
	}

	files := getCueFiles(path)
	err := parseCue(files)
	if err != nil {
		fmt.Printf("error parsing CUE files: %v\n", err)
		os.Exit(1)
	}
	switch *out {
	case "files":
		*out = "_rendered"
		renderResultsDir(*out)
	case "stdout":
		renderResultsStdout()
	default:
		fmt.Printf("error: invalid output destination '%s'\nMust be one of: [stdout, files]\n", *out)
		os.Exit(1)
	}
}
