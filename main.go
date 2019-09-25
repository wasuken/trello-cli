package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/adlio/trello"
	"github.com/urfave/cli"
	"os"
	"os/user"
	"regexp"
)

type Config struct {
	API APIConfig
}

type APIConfig struct {
	Apikey string
	Token  string
	Member string
}

func printBoards(member *trello.Member) {
	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		panic("torello client error.")
	}
	for _, board := range boards {
		fmt.Println("name:" + board.Name + ", id:" + board.ID)
	}
}
func printLists(bid string, client *trello.Client) {
	board, er := client.GetBoard(bid, trello.Defaults())
	fmt.Println(board.Name)
	if er != nil {
		panic("please input base command")
	}
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		panic("not get lists")
	}
	for _, list := range lists {
		fmt.Println("name:" + list.Name + ", id:" + list.ID)
		cards, err := list.GetCards(trello.Defaults())
		if err != nil {
			panic("not get cards")
		}
		for _, card := range cards {
			fmt.Println("    name:" + card.Name + ", id:" + card.ID)
		}
	}
}
func addCard(lid, name, desc string, client *trello.Client) {
	list, err := client.GetList(lid, trello.Defaults())
	if err != nil {
		panic("failed get list")
	}
	er := list.AddCard(&trello.Card{Name: name, Desc: desc}, trello.Defaults())
	if er != nil {
		panic("failed create card")
	}
}
func removeCard(cid string, client *trello.Client) {
	card, err := client.GetCard(cid, trello.Defaults())
	if err != nil {
		panic("failed get card")
	}
	er := client.Delete("cards/"+card.ID, trello.Defaults(), card)
	if er != nil {
		fmt.Println(er)
		panic("failed remove card")
	}
}

func archiveList(lid string, client *trello.Client) {
	list, err := client.GetList(lid, trello.Defaults())
	if err != nil {
		panic("failed get list")
	}
	er := client.Put("lists/"+lid, trello.Arguments{"closed": "true"}, list)
	if er != nil {
		fmt.Println(er)
		panic("failed archive list")
	}
}

func moveCard(cid, after_lid string, client *trello.Client) {
	card, err := client.GetCard(cid, trello.Defaults())
	if err != nil {
		panic("card not found")
	}
	er := card.MoveToList(after_lid, trello.Defaults())
	if er != nil {
		panic("failed move card")
	}
}

func addList(bid, name string, client *trello.Client) {
	board, err := client.GetBoard(bid, trello.Defaults())
	if err != nil {
		panic("failed get board")
	}
	_, er := board.CreateList(name, trello.Defaults())
	if er != nil {
		panic("failed create list")
	}
}

func searchList(bid, query string, client *trello.Client) map[string]string {
	r := regexp.MustCompile(query)
	board, err := client.GetBoard(bid, trello.Defaults())
	if err != nil {
		panic("failed get board")
	}
	lists, er := board.GetLists(trello.Defaults())
	if er != nil {
		panic("failed get lists")
	}
	result := map[string]string{}
	for _, list := range lists {
		if r.MatchString(list.Name) {
			result[list.ID] = list.Name
		}
	}
	return result
}

func moveAllCards(fromListId, toBoardId, toListId string, client *trello.Client) {
	list, err := client.GetList(fromListId, trello.Defaults())
	if err != nil {
		panic("failed get list")
	}
	lists := []*trello.List{list}
	er := client.Post("lists/"+fromListId+"/moveAllCards",
		trello.Arguments{"idBoard": toBoardId, "idList": toListId}, lists)
	if er != nil {
		panic("failed post lists/moveAllCards")
	}
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
			Name:    "boards",
			Aliases: []string{},
			Usage:   "print boards",
			Action: func(c *cli.Context) error {
				member, err := client.GetMember(config.API.Member, trello.Defaults())
				if err != nil {
					fmt.Println(err)
					panic("config file not found.")
				}
				printBoards(member)
				return nil
			},
		},
		{
			Name:    "lists",
			Aliases: []string{},
			Usage:   "print lists and cards",
			Action: func(c *cli.Context) error {
				bid := c.Args().First()
				printLists(bid, client)
				return nil
			},
		},
		{
			Name:    "removeCard",
			Aliases: []string{},
			Usage:   "remove card",
			Action: func(c *cli.Context) error {
				cid := c.Args().First()
				removeCard(cid, client)
				return nil
			},
		},
		{
			Name:    "addCard",
			Aliases: []string{},
			Usage:   "add card",
			Action: func(c *cli.Context) error {
				lid := c.Args().First()
				name := c.Args().Get(1)
				desc := c.Args().Get(2)
				addCard(lid, name, desc, client)
				return nil
			},
		},
		{
			Name:    "moveCard",
			Aliases: []string{},
			Usage:   "move card",
			Action: func(c *cli.Context) error {
				cid := c.Args().First()
				lid := c.Args().Get(1)
				moveCard(cid, lid, client)
				return nil
			},
		},
		{
			Name:    "archiveList",
			Aliases: []string{},
			Usage:   "archive list",
			Action: func(c *cli.Context) error {
				lid := c.Args().First()
				archiveList(lid, client)
				return nil
			},
		},
		{
			Name:    "addList",
			Aliases: []string{},
			Usage:   "add list",
			Action: func(c *cli.Context) error {
				bid := c.Args().First()
				name := c.Args().Get(1)
				addList(bid, name, client)
				return nil
			},
		},
		{
			Name:    "searchList",
			Aliases: []string{},
			Usage:   "search list",
			Action: func(c *cli.Context) error {
				bid := c.Args().First()
				query := c.Args().Get(1)
				for key, value := range searchList(bid, query, client) {
					fmt.Println(key + " " + value)
				}
				return nil
			},
		},
		{
			Name:    "moveAllCards",
			Aliases: []string{},
			Usage:   "move All Cards",
			Action: func(c *cli.Context) error {
				fromListId := c.Args().First()
				toBoardId := c.Args().Get(1)
				toListId := c.Args().Get(2)
				moveAllCards(fromListId, toBoardId, toListId, client)
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
