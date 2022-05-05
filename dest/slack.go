package dest

import (
	"fmt"

	"github.com/slack-go/slack"
)

// channelは#をいれない。
func SlackSendContent(token, channel, content string) {
	fmt.Println(token)
	fmt.Println(content)
	c := slack.New(token)
	ch := "#" + channel

	_, _, err := c.PostMessage(ch, slack.MsgOptionText(content, true))
	if err != nil {
		panic(err)
	}
}
