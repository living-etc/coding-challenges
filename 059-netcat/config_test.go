package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_ParseConfig(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected Config
	}{
		{"No args", []string{}, Config{}},
		{
			"Listen Mode (TCP)",
			[]string{"-l", "-p", "8888"},
			Config{Listen: &ListenConfig{Port: "8888"}, Scan: nil},
		},
		{
			"Listen Mode (UDP)",
			[]string{"-l", "-p", "8888", "-u"},
			Config{Listen: &ListenConfig{Port: "8888", Udp: true}, Scan: nil},
		},
		{
			"Scan Mode (Single Port) (TCP)",
			[]string{"-z", "localhost", "8888"},
			Config{
				Listen: nil,
				Scan: &ScanConfig{
					Host:  "localhost",
					Ports: []int{8888},
				},
			},
		},
		{
			"Scan Mode (Range of Port) (TCP)",
			[]string{"-z", "localhost", "8885-8888"},
			Config{
				Listen: nil,
				Scan: &ScanConfig{
					Host:  "localhost",
					Ports: []int{8885, 8886, 8887, 8888},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseConfig(tc.args)
			want := tc.expected

			if !cmp.Equal(got, want) {
				t.Fatalf("%v: want %+v, got %+v", tc.name, want, got)
			}
		})
	}
}
