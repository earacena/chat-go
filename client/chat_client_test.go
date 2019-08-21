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

