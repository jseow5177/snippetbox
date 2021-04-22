package main

import (
	"testing"
	"time"
)


func TestFormatDate(t *testing.T) {
	// Create a slice of annonymous structs containing the test case name,
	// input to the formatDate() function and the expected output
	tests := []struct{
		name string
		tm time.Time
		want string
	} {
		{
			name: "UTC",
			tm: time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			tm: time.Time{}, // Empty time
			want: "",
		},
		{
			name: "CET Time Zone",
			// CET is one hour ahead of UTC
			tm: time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			// Should expect UTC time (one hour before)
			want: "17 Dec 2020 at 09:00",
		},
	}
	
	// Loop over table of test cases
	for _, tt := range tests {
		// Use t.Run() to run a sub-test for each test case. The first parameter
		// to this is the name of the test (to identify the sub-test in any log output)
		// and the second parameter is an annonymous function containing the actual test for 
		// each case.
		t.Run(tt.name, func(t *testing.T) {
			fd := formatDate(tt.tm)

			if fd != tt.want {
				t.Errorf("want %s; got %s", tt.want, fd)
			}
		})
	}

}