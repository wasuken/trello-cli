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
	er := client.Delete("cards/"+card.ID, trello.Defaults(), nil)
	if er != nil {
		fmt.Println(er)
		panic("failed remove card")
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
func main() {
	if len(os.Args) < 2 {
		fmt.Println("command not found")
		return
	}
	var config Config
	usr, _ := user.Current()
	_, err := toml.DecodeFile(usr.HomeDir+"/.config/trello-cli/config.toml", &config)
	if err != nil {
		fmt.Println(err)
		panic("config file not found.")
	}
	client := trello.NewClient(config.API.Apikey, config.API.Token)
	member, err := client.GetMember(config.API.Member, trello.Defaults())
	if err != nil {
		panic("member error.")
	}
	if os.Args[1] == "boards" {
		printBoards(member)
	} else if os.Args[1] == "lists" {
		if len(os.Args) < 3 {
			fmt.Println("lists <board-id>")
			return
		}
		bid := os.Args[2]
		printLists(bid, client)
	} else if os.Args[1] == "addCard" {
		if len(os.Args) < 5 {
			fmt.Println("addCard <list-id> <card-name> <card-description>")
			return
		}
		lid := os.Args[2]
		name := os.Args[3]
		desc := os.Args[4]
		addCard(lid, name, desc, client)
	} else if os.Args[1] == "removeCard" {
		if len(os.Args) < 3 {
			fmt.Println("removeCard <card-id>")
			return
		}
		cid := os.Args[2]
		removeCard(cid, client)
	} else if os.Args[1] == "moveCard" {
		if len(os.Args) < 4 {
			fmt.Println("moveCard <card-id> <list-id>")
			return
		}
		cid := os.Args[2]
		lid := os.Args[3]
		moveCard(cid, lid, client)
	} else {
		fmt.Println("command not found")
		fmt.Println("command = [boards lists addCard removeCard moveCard]")
		return
	}
}
