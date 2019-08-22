// Name of file: chat_server_test.go
// Author:       Emanuel Aracena
// Date created: August 21, 2019
// Description:  Unit tests for chat_server.go

package main

import (
	"testing"
	"bufio"
	"strings"
)

func TestReceiveMessage(t *testing.T) {	
	b := []byte(`{"ID":"TestID","Body":"Test Message."}`)
	s := string(b) + "\n"
	result := ReceiveMessage(bufio.NewReader(strings.NewReader(s)), true)

	if result != string(b) + "\n"  {
		t.Errorf("ReceiveMessage(strings.NewReader(s), true) = %s; want %s",
			result, string(b) + "\n")
	}
}

func TestDecodeMessage(t *testing.T) {
	encm := []byte(`{"ID":"TestID","Body":"Test Message."}`)
	m := DecodeMessage(encm)
	if m.ID != "TestID" || m.Body != "Test Message." {
		t.Errorf("m = DecodeMessage(encm); m.ID = %s; m.Body = %s; want m.ID = %s;" +
			" m.Body = %s", m.ID, m.Body, "TestID", "Test Message.")
	}
}

func TestFormatMessage(t *testing.T) {
	m := Message {
		ID: "TestID",
		Body: "TestMessage",
	}

	fm, time := FormatMessage(m)
	want := "[+] (" + time + ") " + m.ID + " " + m.Body + "\n"
	if fm != want {
		t.Errorf("FormatMessage(m) = %s; want = %s", fm, want)
	}
}
