package twitch

import (
	"testing"
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
