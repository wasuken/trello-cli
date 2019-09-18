package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/adlio/trello"
	"os"
	"os/user"
)

type Config struct {
	API APIConfig
}

type APIConfig struct {
	Apikey string
	Token  string
}

func printBoards(member trello.Member) {
	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		panic("torello client error.")
	}
	for _, board := range boards {
		fmt.Println("name:" + board.Name + ", id:" + board.ID)
	}
}
func printLists(bid string, client trello.Client) {
	board, er := client.GetBoard(bid, trello.Defaults())
	fmt.Println(board.Name)
	if er != nil {
		panic("not get board")
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
func addCard(lid, name, desc string, client trello.Client) {
	list, err := client.GetList(lid, trello.Defaults())
	if err != nil {
		panic("failed get list")
	}
	er := list.AddCard(&trello.Card{Name: name, Desc: desc}, trello.Defaults())
	if er != nil {
		panic("failed create card")
	}
}
func removeCard(cid string, client trello.Client) {
	card, err := client.GetCard(cid, trello.Defaults())
	if err != nil {
		panic("failed get card")
	}
	er := client.Delete("cards/"+card.ID, trello.Defaults(), nil)
	if er != nil {
		fmt.Println(er)
		panic("failed remove card")
	}
}
func main() {
	var config Config
	usr, _ := user.Current()
	_, err := toml.DecodeFile(usr.HomeDir+"/.config/torello-cron/config.toml", &config)
	if err != nil {
		fmt.Println(err)
		panic("config file not found.")
	}
	client := trello.NewClient(config.API.Apikey, config.API.Token)
	member, err := client.GetMember("5d800d1e3ca8aa7f4be3e12b", trello.Defaults())
	if err != nil {
		panic("member error.")
	}
	if os.Args[1] == "boards" {
		printBoards(member)
	} else if os.Args[1] == "lists" {
		if len(os.Args) < 3 {
			return
		}
		bid := os.Args[2]
		printList(bid, client)
	} else if os.Args[1] == "addCard" {
		if len(os.Args) < 5 {
			return
		}
		lid := os.Args[2]
		name := os.Args[3]
		desc := os.Args[4]
		addCard(lid, name, desc)
	} else if os.Args[1] == "removeCard" {
		if len(os.Args) < 3 {
			return
		}
		cid := os.Args[2]
		removeCard(cid, client)
	}
}
