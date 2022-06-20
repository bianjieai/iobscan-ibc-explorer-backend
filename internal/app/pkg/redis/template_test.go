package redis

import (
	"reflect"
	"testing"
)

func TestFormatAny(t *testing.T) {
	type Student struct {
		Name string `json:"name"`
	}
	student := Student{Name: "zhang"}

	var tests = []struct {
		input interface{}
		want  string
	}{
		{"hello", "hello"},
		{1, "1"},
		{3.23, "3.23"},
		{true, "true"},
		{false, "false"},
		{student, "{\"name\":\"zhang\"}"},
		{-33, "-33"},
		{0, "0"},
	}

	for _, test := range tests {
		got := formatAny(reflect.ValueOf(test.input))
		if got != test.want {
			t.Errorf("input: %q, want: %q", got, test.want)
		}
	}
}
