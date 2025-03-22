package twitch

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseMessage(t *testing.T) {
	var tests = []struct {
		name       string
		twitchMsg  string
		source     string
		command    string
		subCommand string
		message    string
	}{
		{
			"Logging in",
			":tmi.twitch.tv 001 <user> :Welcome, GLHF!",
			":tmi.twitch.tv",
			"001",
			"<user>",
			"Welcome, GLHF!",
		},
		{
			"Keep alive",
			"PING :tmi.twitch.tv",
			"",
			"PING",
			":tmi.twitch.tv",
			"",
		},
		{
			"Chat message",
			":foo!foo@foo.tmi.twitch.tv PRIVMSG #bar :bleedPurple",
			":foo!foo@foo.tmi.twitch.tv",
			"PRIVMSG",
			"#bar",
			"bleedPurple",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsed := parseMessage(test.twitchMsg)
			if parsed.source != test.source {
				t.Fail()
			}
			if parsed.command != test.command {
				t.Fail()
			}
			if parsed.subcommand != test.subCommand {
				t.Fail()
			}
			if parsed.message != test.message {
				t.Fail()
			}
		})
	}
}

func TestConnectClient(t *testing.T) {
	// server mock
	go func(address, port string) {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
		if err != nil {
			panic(err)
		}
		defer listener.Close()
		// listen for single connection
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		// work
		buffReader := bufio.NewReader(conn)
		runFlag := true
		passReceived := false
		userReceived := false
		welcomeSend := false
		user := ""
		defer conn.Close()
		//read cycle
		for runFlag {
			conn.SetReadDeadline(time.Now().Add(time.Second * 5))
			bytes, err := buffReader.ReadBytes('\n')
			if err != nil {
				switch {
				case errors.Is(err, os.ErrDeadlineExceeded):
					t.Fatal("read deadline reached")
				case errors.Is(err, io.EOF):
					runFlag = false
				default:
					t.Fatal(err)
				}
			}
			msg := strings.Split(string(bytes), " ")
			if msg[0] == "PASS" {
				passReceived = true
			} else if msg[0] == "NICK" {
				userReceived = true
				user = strings.TrimSpace(msg[1])
			} else if msg[0] == "JOIN" {
				channel := strings.TrimSpace(msg[1])
				if welcomeSend {
					conn.Write([]byte(fmt.Sprintf(
						":%[1]s!%[1]s@%[1]s.tmi.twitch.tv JOIN #%[2]s\r\n:%[1]s.tmi.twitch.tv 353 %[1]s = #%[2]s :%[1]s\r\n:%[1]s.tmi.twitch.tv 366 %[1]s #%[2]s :End of /NAMES list\r\n",
						user,
						channel,
					)))
					// end loop/connection
					runFlag = false
				}
			} else {
				t.Fatalf("unespected command %s", msg[0])
			}
			if passReceived && userReceived && !welcomeSend {
				conn.Write([]byte(fmt.Sprintf(
					":tmi.twitch.tv 001 %[1]s :Welcome, GLHF!\r\n:tmi.twitch.tv 002 %[1]s :Your host is tmi.twitch.tv\r\n:tmi.twitch.tv 003 %[1]s :This server is rather new\r\n:tmi.twitch.tv 004 %[1]s :-\r\n:tmi.twitch.tv 375 %[1]s :-\r\n:tmi.twitch.tv 372 %[1]s :You are in a maze of twisty passages.\r\n:tmi.twitch.tv 376 %[1]s :>\r\n",
					user,
				)))
				welcomeSend = true
			}
		}
	}("localhost", "8888")

	// create client
	client := TwitchClient{
		TwitchIRL:  "localhost:8888",
		Channel:    "testchannel",
		WriterMain: func(string) {},
		WriterCmd:  func(string) {},
	}
	err := client.ConnectChannel()
	// assert
	if err != nil {
		t.Fatal(err)
	}
}

type uiBuffer struct {
	data []string
}

func (buf *uiBuffer) Write(line string) {
	buf.data = append(buf.data, line)
}

func TestMessageClient(t *testing.T) {
	// server mock
	go func(address, port string) {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
		if err != nil {
			panic(err)
		}
		defer listener.Close()
		// listen for single connection
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		// work
		buffReader := bufio.NewReader(conn)
		// buffWriter := bufio.NewWriter(conn)
		defer conn.Close()
		// send single message
		fmt.Fprint(
			conn,
			":user!user@user.tmi.twitch.tv PRIVMSG #channel :this is a message\r\n",
		)
		// time.Sleep(time.Millisecond * 100)
		// send Ping heartbeat
		fmt.Fprint(conn, "PING :tmi.twitch.tv\r\n")
		time.Sleep(time.Second)
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		bytes, err := buffReader.ReadBytes('\n')
		// bytes, _, err := buffReader.ReadLine()
		if err != nil {
			switch {
			//case timeout
			case errors.Is(err, os.ErrDeadlineExceeded):
				t.Fatal(err)
			case errors.Is(err, io.EOF):
				t.Fatal(err)
			default:
				t.Fatal(err)
			}
		}
		if strings.TrimSpace(string(bytes)) != "PONG :tmi.twitch.tv" {
			// fail
			t.Fatalf("PONG message not received: %s", string(bytes))
		}
		conn.Close()
	}("localhost", "8889")

	// create client
	// mainBuffer := make([]string, 0)
	mainBuffer := uiBuffer{data: make([]string, 0)}
	cmdBuffer := uiBuffer{data: make([]string, 0)}
	client := TwitchClient{
		TwitchIRL:  "localhost:8889",
		Channel:    "testchannel",
		WriterMain: mainBuffer.Write,
		WriterCmd:  cmdBuffer.Write,
	}
	// connect to mock server
	var err error
	client.Conn, err = net.Dial("tcp", client.TwitchIRL)
	if err != nil {
		panic(err)
	}
	client.ReadChat()
	// assert
	if len(mainBuffer.data) != 1 {
		t.Fail()
	}
	if mainBuffer.data[0] != "user:> this is a message\n" {
		t.Fail()
	}
	if len(cmdBuffer.data) != 0 {
		t.Fail()
	}
}
