package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
)

var wolframClient *wolfram.Client

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {

		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}
func main() {
	godotenv.Load(".env")
	bot := slacker.NewClient(os.Getenv("oauth_bottokens"), os.Getenv("socket_token"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("wolfram_id")}
	client := witai.NewClient(os.Getenv("wit.ai_token"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("i have a doubt - <message>", &slacker.CommandDefinition{
		Description: "Send any question :",
		Examples:    []string{"who is the chief minister of tamilnadu"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")
			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			data, _ := json.MarshalIndent(msg, "", "    ")
			softcopy := string(data[:])
			val := gjson.Get(softcopy, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			ans := val.String()
			res, err := wolframClient.GetSpokentAnswerQuery(ans, wolfram.Metric, 1000)
			if err != nil {
				fmt.Println("Here Is An Error")
			}
			fmt.Println(val)
			response.Reply(res)

		},
	})

	strt, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := bot.Listen(strt)
	if err != nil {
		log.Fatal(err)
	}
}
