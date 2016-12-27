package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeNumberString(t *testing.T) {
	type TestCase struct {
		Name string

		Text     string
		Expected string
	}

	for _, tc := range []TestCase{
		{"Null", "", ""},
		{"Non-null", "asd", "\"=\"\"asd\"\"\""},
	} {
		actual := escapeNumberString(tc.Text)
		assert.Equal(t, tc.Expected, actual, fmt.Sprintf("Test case %s", tc.Name))
	}
}

func TestWriteResults(t *testing.T) {
	var tests = []struct {
		name string

		component Component
		expected  string
	}{
		{"Name-only", Component{Name: "blah"},
			header + fmt.Sprintf(template, "blah", "", "", "", "", "", "", "", "", "", "")},
		{"Full", Component{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", true, "11", true, "12", "13", "14", "15", "16", "17", "18", "19", true},
			header + fmt.Sprintf(template, "1", "4", "5", "6", "7", "8", "9", "14", "15", "16", "17")},
	}
	for _, tc := range tests {
		var written int
		var actual string

		writer := bytes.NewBufferString("")
		ch := make(chan Component)
		done := make(chan bool)

		go func() {
			written = writeResults(ch, writer)
			done <- true
		}()

		ch <- tc.component
		close(ch) // close channel to signal writeResults to exit

		<-done // block on signal channel to avoid data race

		actual = writer.String()

		assert.Equal(t, tc.expected, actual, fmt.Sprintf("Test case: %s", tc.name))
		assert.Equal(t, 1, written, fmt.Sprintf("Test case: %s", tc.name))
	}
}
