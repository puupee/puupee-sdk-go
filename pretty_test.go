package puupeesdk

import "testing"

type testPrintObject struct {
	Username string
	Password string
}

func TestPrintObject(t *testing.T) {
	PrintObject(&testPrintObject{
		Username: "test_user_name",
		Password: "hahahaah",
	})
}

func TestPrintObjectArray(t *testing.T) {
	PrintArray([]*testPrintObject{
		{
			Username: "test_user_name",
			Password: "hahahaah",
		},
		{
			Username: "aaaaaaa",
			Password: "bbbbbbbb",
		},
	})
}
