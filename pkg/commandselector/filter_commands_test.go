package commandselector

import (
	"reflect"
	"testing"
)

func TestFilterCommands(t *testing.T) {
	commands := []Command{
		{PackageName: "@internal/package1", CommandName: "build", CommandValue: "webpack"},
		{PackageName: "@internal/package1", CommandName: "test", CommandValue: "jest"},
		{PackageName: "package2", CommandName: "build", CommandValue: "webpack"},
		{PackageName: "package2", CommandName: "test", CommandValue: "jest"},
		{PackageName: "@internal/package3", CommandName: "foo", CommandValue: "some-command"},
		{PackageName: "@internal/package3", CommandName: "bar", CommandValue: "some-command"},
	}

	tests := []struct {
		query    string
		expected []Command
	}{
		{
			query:    "",
			expected: commands,
		},
		{
			query:    " ",
			expected: commands,
		},
		{
			query:    "build",
			expected: []Command{commands[0], commands[2]},
		},
		{
			query: "webpack",
			expected: []Command{
				commands[0],
				commands[2],
			},
		},
		{
			query:    "package2 test",
			expected: []Command{commands[3]},
		},
		{
			query:    "nonexistent",
			expected: []Command{},
		},
		{
			query:    "@inter",
			expected: []Command{commands[0], commands[1], commands[4], commands[5]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			result := filterCommands(commands, tt.query)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("\n"+"query: %q"+"\n"+"expected: %v"+"\n"+"result: %v", tt.query, tt.expected, result)
			}
		})
	}
}
