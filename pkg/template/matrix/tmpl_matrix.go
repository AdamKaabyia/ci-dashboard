package matrix

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	v1 "github.com/openshift-psap/ci-dashboard/api/matrix/v1"
)

type TemplateBase struct {
	Spec        *v1.MatricesSpec
	Description string
	Date        string
}

func Generate(matrixTemplate string, matrices *v1.MatricesSpec, date string) ([]byte, error) {
	matrix_template, err := ioutil.ReadFile(matrixTemplate)
	if err != nil {
		return []byte{}, fmt.Errorf("Matrix template file %s cannot be read: %v", matrixTemplate, err)
	}

	tmpl_data := TemplateBase{
		Spec: matrices,
		Date: date,
	}

	fmap := template.FuncMap{
		"md_section": func(s string) string {
			return strings.Repeat("=", utf8.RuneCountInString(s))
		},
		"md_subsection": func(s string) string {
			return strings.Repeat("-", utf8.RuneCountInString(s))
		},
		"unescape_html": func(s string) template.HTML {
			return template.HTML(s)
		},
		"nb_last_test": func() string {
			return fmt.Sprintf("%d", matrices.TestHistory)
		},
		"no_test_history": func(test v1.TestSpec) []int {
			arr := []int{}
			for i := len(test.OldTests); i < matrices.TestHistory; i++ {
				arr = append(arr, i)
			}
			return arr
		},
		"group_name": func(txt string) string {
			pipe_pos := strings.Index(txt, "|")
			if pipe_pos == -1 {
				return txt
			} else {
				return txt[pipe_pos+1:]
			}
		},
		"artifacts_url": func(matrix v1.MatrixSpec, test v1.TestResult) string {
			if test.TestSpec == nil {
				return "INVALID"
			}
			var prow_step = matrix.ProwStep
			if test.TestSpec.ProwStep != "" {
				// override test_matrix.ProwStep if ProwStep is test_spec.ProwStep is specified
				prow_step = test.TestSpec.ProwStep
			}

			var base string
			if matrix.ProwType == "presubmit" {
				if test.PullNumber == "" {
					print("Missing pull number for %s", test.TestSpec.TestName)
				}
				base = fmt.Sprintf("%s/pull/%s/%s/%s/",
					matrix.ArtifactsURL,
					test.PullNumber,
					test.TestSpec.ProwName,
					test.BuildId)
			} else {
				base = fmt.Sprintf("%s/%s/%s/artifacts/%s/%s",
					matrix.ArtifactsURL,
					test.TestSpec.ProwName,
					test.BuildId,
					test.TestSpec.TestName,
					prow_step)
			}

			if test.TestSpec.IsCiOperator == nil || *test.TestSpec.IsCiOperator {
				return base + "/artifacts"
			}
			return base
		},
		"spyglass_url": func(matrix v1.MatrixSpec, prowName string, test v1.TestResult) string {
			return fmt.Sprintf("%s/%s/%s", matrix.ViewerURL, prowName, test.BuildId)
		},
		"repository_url": func(matrix v1.MatrixSpec, test v1.TestResult) string {
			base := matrix.RepositoryURL
			if base == "" {
				base = "https://github.com/openshift-psap/ci-artifacts"
			}
			return fmt.Sprintf("%s/commit/%s", base, test.CiArtifactsVersion)
		},
		"test_status_descr": func(test v1.TestResult, status string) string {
			if status == "success" {
				return "Test passed"
			} else if status == "known_flake" {
				msg := "Test failed because of a known flake: "
				for _, flake := range test.Messages[v1.TestMessageTypeFlake] {
					msg += "\n- " + flake
				}
				return msg
			} else if status == "step_success" {
				return "Test failed but the operator step passed"
			} else if status == "step_failed" {
				return "Test failed because the operator step failed"
			} else if status == "step_missing" {
				return "Test failed but operator step wasn't executed"
			} else {
				return fmt.Sprintf("Test: %t, Step: %t (status: %s)",
					test.Passed, test.StepPassed, status)
			}
		},
		"test_status": func(test v1.TestResult) string {
			if test.Passed {
				return "success"
			} else if len(test.Messages[v1.TestMessageTypeFlake]) != 0 {
				return "known_flake"
			} else if !test.StepExecuted {
				return "step_missing"
			} else if test.StepPassed {
				return "step_success"
			} else if !test.StepPassed {
				return "step_failed"
			} else {
				return "parsing_error"
			}
		},
		"test_messages": func(message_type string, test v1.TestResult) map[string]string {
			if message_type == "flake" {
				return test.Messages[v1.TestMessageTypeFlake]
			} else if message_type == "info" {
				return test.Messages[v1.TestMessageTypeInfo]
			} else if message_type == "warning" {
				return test.Messages[v1.TestMessageTypeWarning]
			} else if message_type == "error" {
				return test.Messages[v1.TestMessageTypeError]
			}
			return nil
		},
		"test_message_types": func() []string {
			return []string{"flake", "info", "warning", "error"}
		},
	}

	tmpl := template.Must(template.New("runtime").Funcs(fmap).Parse(string(matrix_template)))

	var buff bytes.Buffer
	if err = tmpl.Execute(&buff, tmpl_data); err != nil {
		return []byte{}, fmt.Errorf("Matrix template file %s could not applied: %v", matrixTemplate, err)
	}

	generated_html := buff.Bytes()

	return generated_html, nil
}
