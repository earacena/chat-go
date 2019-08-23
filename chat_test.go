// Name of file: chat_test.go
// Author:       Emanuel Aracena
// Date created: August 23, 2019
// Description:  Unit tests for chat.go


package main

import (
	"testing"
	"bufio"
	"strings"
	"bytes"
	"net"
)

func TestConnectToServer(t *testing.T) {
	l, _ := net.Listen("tcp", "localhost:9999")
	defer l.Close()
	
	c := Client { ipAndPort: "localhost:9999", }
	err := c.ConnectToServer()
	if err != nil {
		t.Errorf("ConnectToServer = %s; want nil", err.Error())
	}
}	

func TestReadUserInput(t *testing.T) {
	given1 := "ABC123\n"
	given2 := "23 123 er f #$  4\n"
	given3 := "         \n"

	r1 := bufio.NewReader(strings.NewReader(given1))
	r2 := bufio.NewReader(strings.NewReader(given2))
	r3 := bufio.NewReader(strings.NewReader(given3))

	result1 := ReadUserInput(r1, false)
	result2 := ReadUserInput(r2, false)
	result3 := ReadUserInput(r3, false)

	if given1 != result1 {
		t.Errorf("ReadUserInput(r1, false) = %s; want %s", result1, given1)
	}
	
	if given2 != result2 {
		t.Errorf("ReadUserInput(r2, false) = %s; want %s", result2, given2)
	}
	
	if given3 != result3 {
		t.Errorf("ReadUserInput(r3, false) = %s; want %s", result3, given3)
	}
}

func TestCheckHaltCommand(t *testing.T) {
	s1 := "!halt"
	s2 := "halt"
	s3 := "!hal"

	b1 := CheckHaltCommand(s1, true)
	b2 := CheckHaltCommand(s2, true)
	b3 := CheckHaltCommand(s3, true)

	if b1 != true {
		t.Errorf("CheckHaltCommand(s1) = false; want true")
	}

	if b2 != false {
		t.Errorf("CheckHaltCommand(s1) = true; want false")
	}
	
	if b3 != false {
		t.Errorf("CheckHaltCommand(s1) = true; want false")
	}
}

func TestEncodeMessage(t *testing.T) {
	m := Message {
		ID: "TestID",
		Body: "Test Message",
	}

	b := EncodeMessage(m)
	if bytes.Equal(b, []byte(`{"ID":"TestID","Body":"Test Message"}`)) != true {
		t.Errorf("EncodeMessage(m) = %s; want %s", b,
			[]byte(`{"ID":"TestID","Body":"Test Message"}`))
	}
}
			
func TestSend(t *testing.T) {
	var buf bytes.Buffer
	b := []byte(`"ID":"TestID","Body":"Test Message."`)
	Send(&buf, b)


	
	if buf.String() != string(b) + "\n" {
		t.Errorf("Send(buf, b); buf.String() = %s; string(b) = %s; " +
			"want buf.String() == %s", buf.Bytes(), b, b)
	}
}

func TestChooseID(t *testing.T) {
	given := "UserABC" + "\n"
	r := bufio.NewReader(strings.NewReader(given))
	ID := ChooseID(r, true)
	given = strings.Replace(given, "\n", "", -1)
	if "[" + given + "]" != ID {
		t.Errorf("ChooseID = %s; want %s", ID, "[" + given + "]")
	}
}


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
