/*
 * Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
 * Copyright (c) 2021, Red Hat.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1

const Version = "v1"

type TestMessageType int64

const (
	TestMessageTypeInfo TestMessageType = iota
	TestMessageTypeWarning
	TestMessageTypeError
	TestMessageTypeFlake
)

func (t TestMessageType) String() string {
	switch t {
	case TestMessageTypeInfo:
		return "_INFO"
	case TestMessageTypeWarning:
		return "_WARNING"
	case TestMessageTypeError:
		return "_ERROR"
	case TestMessageTypeFlake:
		return "_FLAKE"
	}
	return "unknown"
}

type MatricesSpec struct {
	Version     string                `json:"version"`
	Description string                `json:"description,omitempty"`
	TestHistory int                   `json:"test_history"`
	Matrices    map[string]MatrixSpec `json:"matrices,omitempty"`
}

type ToolboxStepResult struct {
	Name string

	Ok       int
	Failures int
	Ignored  int

	ExpectedFailure string

	FlakeFailure string
}

type TestResult struct {
	BuildId    string
	Passed     bool
	Result     string
	FinishDate string

	StepExecuted bool
	StepPassed   bool
	StepResult   string

	Messages map[TestMessageType]map[string]string

	/* *** */

	OperatorVersion    string
	OpenShiftVersion   string
	CiArtifactsVersion string

	// New field for pull request number
	PullNumber string `json:"pull_number,omitempty"`

	/* *** */
	TestSpec *TestSpec

	ToolboxSteps []string

	ToolboxStepsResults []ToolboxStepResult

	/* *** */

	Ok       int
	Failures int
	Ignored  int

	FlakeFailure bool
}

type TestSpec struct {
	TestName        string `json:"test_name,omitempty"`
	Branch          string `json:"branch,omitempty"`
	OperatorVersion string `json:"operator_version,omitempty"`
	Variant         string `json:"variant,omitempty"`
	ProwStep        string `json:"prow_step,omitempty"`

	ProwName     string `json:"prow_name,omitempty"`
	IsCiOperator *bool  `json:"is_ci_operator,omitempty"`

	// New field to indicate the type of Prow job.
	// If not specified, we default to "periodic".
	ProwType string `json:"prow_type,omitempty"`

	/* *** */

	Matrix *MatrixSpec

	TestGroup string

	OldTests []*TestResult
}

type MatrixSpec struct {
	Description    string                `json:"description,omitempty"`
	ViewerURL      string                `json:"viewer_url,omitempty"`
	ArtifactsURL   string                `json:"artifacts_url,omitempty"`
	ArtifactsCache string                `json:"artifacts_cache,omitempty"`
	ProwConfig     string                `json:"prow_config,omitempty"`
	ProwStep       string                `json:"prow_step,omitempty"`
	OperatorName   string                `json:"operator_name,omitempty"`
	RepositoryURL  string                `json:"repository_url,omitempty"`
	Tests          map[string][]TestSpec `json:"tests,omitempty"`

	/* *** */
	ProwType string `json:"prow_type,omitempty"`

	Name string
}
