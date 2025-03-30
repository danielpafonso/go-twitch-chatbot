# Go Twicth Chatbot

Twitch Chat Reader bot that enables to filter and reaction to commands

|             | WebSocket clients               | IRC clients                   |
| ----------- | ------------------------------- | ----------------------------- |
| **SSL**     | wss://irc-ws.chat.twitch.tv:443 | irc://irc.chat.twitch.tv:6697 |
| **Non SSL** | ws://irc-ws.chat.twitch.tv:80   | irc://irc.chat.twitch.tv:6667 |

# Rate limits

The Twitch IRC server enforces the following limits. It is up to your bot to keep track of its usage and not exceed the limits. Rate limit counters begin when the server processes the first message and resets at the end of the window. For example, if the limit is 20 messages per 30 seconds, the window starts when the server processes the first message and lasts for 30 seconds. At the end of the window, the counter resets and a new window begins with the next message.
Command and message rate limits

The following tables show the rate limits for the number of messages that your bot may send. If you exceed these limits, Twitch ignores the bots messages for the next 30 minutes.

| Limit                       | Description                                                                                                         |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| 20 messages per 30 seconds  | If the user isn’t the channel’s broadcaster or moderator, the bot may send a maximum of 20 messages per 30 seconds. |
| 100 messages per 30 seconds | If the user is the channel’s broadcaster or moderator, the bot may send a maximum of 100 messages per 30 seconds.   |

laksdjlaksdj
