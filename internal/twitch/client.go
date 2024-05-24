package twitch

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type TwitchClient struct {
	TwitchIRL    string
	Channel      string
	Conn         net.Conn
	Debug        bool
	WriterMain   func(string)
	WriterCmd    func(string)
	WriterConfig func(string, ...bool)
	Plugins      []string
}

type chatMsg struct {
	source     string
	command    string
	subcommand string
	message    string
}

// NewClient creates a Twitch Client and sets desired configurations
func NewClient(irl, channel string, debug bool, writerMain func(string), writerCmd func(string), writerConfig func(string, ...bool)) *TwitchClient {
	client := TwitchClient{
		TwitchIRL:    irl,
		Channel:      channel,
		Debug:        debug,
		WriterMain:   writerMain,
		WriterCmd:    writerCmd,
		WriterConfig: writerConfig,
		Plugins:      make([]string, 0),
	}
	return &client
}

// Parse chat message into source, command , subcommand and message
func parseMessage(line string) chatMsg {
	parsed := chatMsg{}
	if strings.HasPrefix(line, ":") {
		sline := strings.Split(line, " :")
		// msg
		if len(sline) == 2 {
			parsed.message = sline[1]
		}
		// cmd
		scmd := strings.SplitN(sline[0], " ", 3)
		parsed.source = scmd[0]
		parsed.command = scmd[1]
		parsed.subcommand = scmd[2]

	} else {
		scmd := strings.SplitN(line, " ", 2)
		parsed.command = scmd[0]
		parsed.subcommand = scmd[1]
	}
	return parsed
}

// Close Close connetion to IRC server
func (client *TwitchClient) Close() {
	client.Conn.Close()
}

func (client *TwitchClient) WriteCurrentConfigs() {
	client.WriterConfig("Configurations:\n", true)
	client.WriterConfig(fmt.Sprintf(" Channel: %s\n", client.Channel))
	client.WriterConfig(fmt.Sprintf(" Debug mode: %v\n", client.Debug))
	client.WriterConfig("\nPlugins:\n")
	for _, plugin := range client.Plugins {
		client.WriterConfig(fmt.Sprintf(" %s\n", plugin))
	}
}

func (client *TwitchClient) StartBot(uiStarted chan struct{}) error {
	// wait for UI to start
	<-uiStarted
	//time.Sleep(250 * time.Millisecond)
	// write configuratons
	client.WriteCurrentConfigs()
	err := client.ConnectChannel()
	if err != nil {
		return err
	}

	client.ReadChat()

	return nil
}

// ConnectChannel create connection to IRC server and join Channel
func (client *TwitchClient) ConnectChannel() error {
	var err error
	// Initiate connection
	client.Conn, err = net.Dial("tcp", client.TwitchIRL)
	if err != nil {
		return err
	}
	// join channel
	// client.WriterMain(fmt.Sprintf("Joinning Channel %s\n", client.Channel))
	client.WriterCmd(fmt.Sprintf("Joinning Channel %s\n", client.Channel))
	fmt.Fprintf(client.Conn, "PASS %s\r\nNICK %s\r\nJOIN #%s\r\n", "justinfan6493", "justinfan6493", client.Channel)

	chatBuffer := bufio.NewReader(client.Conn)
	for connecting := true; connecting; {
		bytes, _, err := chatBuffer.ReadLine()
		if err != nil {
			// close connection and buffer
			client.Close()
			return err
		}
		line := string(bytes)
		parsedMsg := parseMessage(line)
		if client.Debug {
			client.WriterMain(fmt.Sprintf("C:%s %s\n", parsedMsg.command, parsedMsg.message))
		}
		if parsedMsg.command == "366" {
			// End of /Names list
			connecting = false
			if client.Debug {
				client.WriterMain("\n")
			}
			client.WriterCmd("Channel Joined\n")
		}
	}
	return nil
}

// ReadChat
func (client *TwitchClient) ReadChat() {
	// create read buffer
	buffReader := bufio.NewReader(client.Conn)
	// read block
	for connected := true; connected; {
		// set deadline for reading
		client.Conn.SetReadDeadline(time.Now().Add(time.Second))
		bytes, _, err := buffReader.ReadLine()
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				if client.Debug {
					client.WriterMain("read timeout, continue\n")
				}
				continue
			}
			client.WriterMain(fmt.Sprintf("ERROR[%T] %s\n", err, err))
			connected = false
		}
		parsedMsg := parseMessage(string(bytes))
		// client.WriterMain(fmt.Sprintf("%+v\n", parsedMsg))
		switch parsedMsg.command {
		case "PING":
			// respond with PONG
			pong := fmt.Sprintf("PONG %s", parsedMsg.subcommand)
			client.Conn.Write([]byte(pong))
		case "PRIVMSG":
			// get user
			user := parsedMsg.source[1:strings.Index(parsedMsg.source, "!")]
			client.WriterMain(fmt.Sprintf("%s:> %s\n", user, parsedMsg.message))
		case "001":
			// Logged in (successfully authenticated).
			fallthrough
		case "002", "003", "004":
			fallthrough
		case "353":
			// Tells you who else is in the chat room you're joining.
			fallthrough
		case "366", "372", "375", "376":
			client.WriterMain(fmt.Sprintf("C:%s %s\n", parsedMsg.command, parsedMsg.message))
		}
	}
}
