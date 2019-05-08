package models

import "testing"

func TestSessionsCount(t *testing.T) {
	num := ReturnCountOfSessions()
	if num < 0 {
		t.Error("count of sessions is not ok")
	}
}
