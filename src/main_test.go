package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	defer func() { os.Stdout = original }()
	if err := run([]string{"--version"}); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if _, err := io.Copy(&output, reader); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "quotes "+Version) || !strings.Contains(output.String(), OfficialSite) {
		t.Fatalf("version output missing build information: %s", output.String())
	}
}

func TestRunStatus(t *testing.T) {
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	defer func() { os.Stdout = original }()
	if err := run([]string{"--status"}); err != nil {
		t.Fatal(err)
	}
	_ = writer.Close()
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "OK\n" {
		t.Fatalf("status output = %q", output)
	}
}

func TestValidateListenAddress(t *testing.T) {
	for _, tc := range []struct {
		address string
		port    string
		wantErr bool
	}{
		{address: "127.0.0.1", port: "8080"},
		{address: "::1", port: "443"},
		{address: "", port: "8080", wantErr: true},
		{address: "localhost", port: "0", wantErr: true},
		{address: "localhost", port: "65536", wantErr: true},
		{address: "localhost", port: "https", wantErr: true},
	} {
		err := validateListenAddress(tc.address, tc.port)
		if (err != nil) != tc.wantErr {
			t.Errorf("validateListenAddress(%q, %q) error = %v", tc.address, tc.port, err)
		}
	}
}
