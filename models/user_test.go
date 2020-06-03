package models

import "testing"

// TestCountUsers checks the number of returned values
func TestCountUsers(t *testing.T) {
	got := Users(1)
	want := 30
	if len(got) != want{
		t.Errorf("Users function returned incorrect amount users")
	}
}