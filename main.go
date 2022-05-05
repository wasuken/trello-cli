package main

import (
	"fmt"
	"os"
	"os/user"
	"time"
	"trello-cli/dest"
	tr "trello-cli/trello"

	"github.com/BurntSushi/toml"
	"github.com/adlio/trello"
	"github.com/urfave/cli"
)

type Config struct {
	API   APIConfig
	Slack SlackConfig
	Gmail GmailConfig
	Ymail YmailConfig
}
type GmailConfig struct {
	Email      string
	Smtpserver string
	Password   string
	Port       int
}
type YmailConfig struct {
	Email      string
	Smtpserver string
	Password   string
	Port       int
}
type SlackConfig struct {
	Token string
}
type APIConfig struct {
	Apikey  string
	Token   string
	Member  string
	Boardid string
	List    []string
}

func main() {
	var config Config
	usr, _ := user.Current()
	_, err := toml.DecodeFile(usr.HomeDir+"/.config/trello-cli/config.toml", &config)
	if err != nil {
		fmt.Println(err)
		panic("config file not found.")
	}
	client := trello.NewClient(config.API.Apikey, config.API.Token)
	app := cli.NewApp()
	app.Name = "trello-cli"
	app.Usage = "trello cli tool"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "notify",
			Aliases: []string{},
			Usage:   "notify list info",
			Action: func(c *cli.Context) error {
				send := c.Args().Get(0)
				bid := config.API.Boardid

				contents, err := tr.WatchListsToContent(bid, client, config.API.List)
				if err != nil {
					panic(err)
				}
				// 関数だけ置換したいけどできるのだろうか
				switch send {
				case "slack":
					channel := c.Args().Get(1)
					for i, content := range contents {
						ind := time.Duration((i % 4) + 1)
						time.Sleep(time.Second * ind)
						dest.SlackSendContent(config.Slack.Token, channel, content)
					}
				case "mail":
					body := "from trello-cli"
					for _, content := range contents {
						body += "\n\n" + content
					}
					email := config.Gmail.Email
					pass := config.Gmail.Password
					srv := config.Gmail.Smtpserver
					port := config.Gmail.Port
					mail := dest.NewMail(email, pass, srv, port)
					mail.Send(email, "[trello] notify", body)
				case "ymail":
					body := "from trello-cli"
					for _, content := range contents {
						body += "\n\n" + content
					}
					email := config.Ymail.Email
					pass := config.Ymail.Password
					srv := config.Ymail.Smtpserver
					port := config.Ymail.Port
					mail := dest.NewMail(email, pass, srv, port)
					mail.Send(email, "[trello] notify", body)
				}
				fmt.Println("finished.")
				return nil
			},
		},
	}

	er := app.Run(os.Args)
	if er != nil {
		fmt.Println(er)
		panic("error app run.")
	}
}
