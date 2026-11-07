package data

import "testing"

func TestReadFile(t *testing.T) {
	content, err := ReadFile("quotes.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Fatal("embedded quote data is empty")
	}
	if _, err := ReadFile("missing.json"); err == nil {
		t.Fatal("ReadFile accepted missing file")
	}
}
