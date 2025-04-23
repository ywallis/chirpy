package main

import "testing"

func Test_profanityFilter(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input  string
		filter map[string]bool
		want   string
	}{
		{
			name:   "test",
			input:  "poop",
			filter: map[string]bool{"poop": true},
			want:   "****",
		},
		{
			name:   "test",
			input:  "hello",
			filter: map[string]bool{"poop": true},
			want:   "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := profanityFilter(tt.input, tt.filter)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("profanityFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
