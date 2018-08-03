package ghokin

import (
	"os/exec"
	"testing"

	"github.com/cucumber/gherkin-go"
	"github.com/stretchr/testify/assert"
)

func TestFormatStepOrExampleLine(t *testing.T) {
	token := &gherkin.Token{Keyword: "Then ", Text: "match some JSON properties"}
	expected := "    Then match some JSON properties\n"

	assert.Equal(t, expected, formatStepOrExampleLine(token))
}

func TestFormatFeatureOrBackgroundLine(t *testing.T) {
	token := &gherkin.Token{Keyword: "Feature", Text: "Set api"}
	expected := "Feature: Set api\n"

	assert.Equal(t, expected, formatFeatureOrBackgroundLine(token))
}

func TestFormatDocStringOrRuleLine(t *testing.T) {
	token := &gherkin.Token{Keyword: `"""`}
	expected := `      """` + "\n"

	assert.Equal(t, expected, formatDocStringOrRuleLine(token))
}

func TestFormatTable(t *testing.T) {
	type scenario struct {
		tokens []*gherkin.Token
		test   func(string)
	}

	scenarios := []scenario{
		{
			[]*gherkin.Token{
				&gherkin.Token{
					Items: []*gherkin.LineSpan{
						&gherkin.LineSpan{Text: "whatever"},
						&gherkin.LineSpan{Text: "whatever whatever"},
					},
				},
				&gherkin.Token{
					Items: []*gherkin.LineSpan{
						&gherkin.LineSpan{Text: "test"},
						&gherkin.LineSpan{Text: "test"},
					},
				},
				&gherkin.Token{
					Items: []*gherkin.LineSpan{
						&gherkin.LineSpan{Text: "t"},
						&gherkin.LineSpan{Text: "t"},
					},
				},
			},
			func(output string) {
				expected := `      | whatever | whatever whatever |
      | test     | test              |
      | t        | t                 |
`
				assert.Equal(t, expected, output)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.test(formatTableRows(scenario.tokens))
	}
}

func TestExtractCommand(t *testing.T) {
	type scenario struct {
		token *gherkin.Token
		test  func(*exec.Cmd)
	}

	commandMatcher = map[string]string{
		"cat": "cat",
		"jq":  "jq",
	}

	scenarios := []scenario{
		{
			&gherkin.Token{
				Text: "",
			},
			func(cmd *exec.Cmd) {
				assert.Nil(t, cmd)
			},
		},
		{
			&gherkin.Token{
				Text: "# A comment",
			},
			func(cmd *exec.Cmd) {
				assert.Nil(t, cmd)
			},
		},
		{
			&gherkin.Token{
				Text: "# @jq",
			},
			func(cmd *exec.Cmd) {
				expected := exec.Command("sh", "-c", "jq")

				assert.Equal(t, expected, cmd)
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.test(extractCommand(scenario.token))
	}
}
