package trello

import (
	"fmt"
	"regexp"

	"github.com/adlio/trello"
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
		fmt.Println("name=>" + board.Name + ", id=>" + board.ID)
	}
}

func find(haystack []string, needle string) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}

// 監視対象Listから送付コンテンツを生成
// 長さ制限対策としてリストごとにコンテンツを分割
// @return string[]
func WatchListsToContent(bid string, client *trello.Client, watch_list []string) ([]string, error) {
	board, er := client.GetBoard(bid, trello.Defaults())
	fmt.Println(board.Name)
	if er != nil {
		return nil, er
	}
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return nil, err
	}
	var rst []string
	for _, list := range lists {
		if !find(watch_list, list.Name) {
			continue
		}
		content := fmt.Sprintf("# %s\n", list.Name)
		cards, err := list.GetCards(trello.Defaults())
		if err != nil {
			return nil, err
		}
		for _, card := range cards {
			line := fmt.Sprintf("\t%s(%s)\n", card.Name, card.ShortURL)

			content += line
		}
		if len(cards) > 0 && len(content) > 0 {
			rst = append(rst, content)
		}
	}
	return rst, nil
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
		fmt.Println("name=>" + list.Name + ", id=>" + list.ID)
		cards, err := list.GetCards(trello.Defaults())
		if err != nil {
			panic("not get cards")
		}
		for _, card := range cards {
			fmt.Println("    name=>" + card.Name + ", id=>" + card.ID)
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
func printCardInfo(cardId string, client *trello.Client) {
	card, err := client.GetCard(cardId, trello.Defaults())
	if err != nil {
		panic("failed get card")
	}
	const layout = "2006-01-02 15:04:05"
	fmt.Println("id=>" + card.ID + "\tname=>" + card.Name + "\tdate=>" + card.DateLastActivity.Format(layout))
}
func printListInfo(listId string, client *trello.Client) {
	list, err := client.GetList(listId, trello.Defaults())
	if err != nil {
		panic("failed get list")
	}
	cards, e := list.GetCards(trello.Defaults())
	if e != nil {
		panic("failed get cards")
	}
	for _, card := range cards {
		printCardInfo(card.ID, client)
	}
}
