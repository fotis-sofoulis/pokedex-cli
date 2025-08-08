package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "classic test 1",
			expected: []string{"classic", "test", "1"},
		},
		{
			input:    "   spaces before",
			expected: []string{"spaces", "before"},
		},
		{
			input:    "spaces after         ",
			expected: []string{"spaces", "after"},
		},
	}

	for _, c := range cases {
		// Check length of the actual and expected slice
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("len of actual not the same as expected. Actual: %d - vs - Expected: %d", len(actual), len(c.expected))
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("word of actual and expected do not match. Actual: %s - vs - Expected: %s", word, expectedWord)
			}
		}
	}

}
