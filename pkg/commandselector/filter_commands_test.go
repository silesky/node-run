package commandselector

import (
	"reflect"
	"testing"
)

func TestFilterCommands(t *testing.T) {
	commands := []Command{
		{PackageName: "package1", CommandName: "build", CommandValue: "npm run build"},
		{PackageName: "package2", CommandName: "test", CommandValue: "npm run test"},
		{PackageName: "package3", CommandName: "start", CommandValue: "npm start"},
	}

	tests := []struct {
		query    string
		expected []Command
	}{
		{
			query:    "build",
			expected: []Command{commands[0]},
		},
		{
			query: "npm run",
			expected: []Command{
				commands[0],
				commands[1],
			},
		},
		{
			query:    "package3 start",
			expected: []Command{commands[2]},
		},
		{
			query:    "nonexistent",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			result := filterCommands(commands, tt.query)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("filterCommands(%q) = %v, want %v", tt.query, result, tt.expected)
			}
		})
	}
}
