package main

import "testing"

func Test_makeMap(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w    int
		h    int
		want [][]int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeMap(tt.w, tt.h)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("makeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

