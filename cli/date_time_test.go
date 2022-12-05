package cli

import (
	"testing"
	"time"
)

func TestParseDateTime(t *testing.T) {
	// formated := "2022-12-05T09:09:25.940827+08:00"
	// formated := "2022-12-05T00:23:43.936007"
	formated := "2022-12-05T09:37:19.355398Z"
	dt, err := time.Parse(time.RFC3339, formated)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dt)
	t.Log(dt.Local())
}
