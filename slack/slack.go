package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

/*
   TODO: Change @BOT_NAME to the same thing you entered when creating your Slack application.
   NOTE: command_arg_1 and command_arg_2 represent optional parameteras that you define
   in the Slack API UI
*/
const helpMessage = "type in '@BOT_NAME <command_arg_1> <command_arg_2>'"

/*
   CreateSlackClient sets up the slack RTM (real-timemessaging) client library,
   initiating the socket connection and returning the client.
   DO NOT EDIT THIS FUNCTION. This is a fully complete implementation.
*/
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

/*
   RespondToEvents waits for messages on the Slack client's incomingEvents channel,
   and sends a response when it detects the bot has been tagged in a message with @<botTag>.

   EDIT THIS FUNCTION IN THE SPACE INDICATED ONLY!
*/
func RespondToEvents(slackClient *slack.RTM) {
	for msg := range slackClient.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)

			// if strings.Contains(ev.Msg.Text, botTagString) {
			// 	continue
			// }
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)

			// TODO: Make your bot do more than respond to a help command. See notes below.
			// Make changes below this line and add additional funcs to support your bot's functionality.
			// sendHelp is provided as a simple example. Your team may want to call a free external API
			// in a function called sendResponse that you'd create below the definition of sendHelp,
			// and call in this context to ensure execution when the bot receives an event.

			// START SLACKBOT CUSTOM CODE
			// ===============================================================
			sendResponse(slackClient, message, ev.Channel)
			sendHelp(slackClient, message, ev.Channel)
			// ===============================================================
			// END SLACKBOT CUSTOM CODE
		default:

		}
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

// sendResponse is NOT unimplemented --- write code in the function body to complete!

func sendResponse(slackClient *slack.RTM, message, slackChannel string) {
	args := strings.Split(message, " ")

	switch command := strings.ToLower(args[0]); command {
	case "say":
		slackClient.SendMessage(slackClient.NewOutgoingMessage(strings.Join(args[1:], " "), slackChannel))
	case "math":
		var a, b int64
		var err error

		a, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			slackClient.SendMessage(slackClient.NewOutgoingMessage("silly human you didn't enter valid math", slackChannel))
			break
		}

		b, err = strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			slackClient.SendMessage(slackClient.NewOutgoingMessage("silly human you didn't enter valid math", slackChannel))
			break
		}

		switch operator := args[2]; operator {
		case "+":
			slackClient.SendMessage(slackClient.NewOutgoingMessage(fmt.Sprintf("%s %s %s = %d", args[1], args[2], args[3], a+b), slackChannel))
		case "-":
			slackClient.SendMessage(slackClient.NewOutgoingMessage(fmt.Sprintf("%s %s %s = %d", args[1], args[2], args[3], a-b), slackChannel))
		case "/":
			if b == 0 {
				slackClient.SendMessage(slackClient.NewOutgoingMessage("nooooooo", slackChannel))
				break
			}
			slackClient.SendMessage(slackClient.NewOutgoingMessage(fmt.Sprintf("%s %s %s = %d", args[1], args[2], args[3], a/b), slackChannel))
		case "*":
			slackClient.SendMessage(slackClient.NewOutgoingMessage(fmt.Sprintf("%s %s %s = %d", args[1], args[2], args[3], a*b), slackChannel))
		default:
			slackClient.SendMessage(slackClient.NewOutgoingMessage("silly human you didn't enter valid math", slackChannel))
		}
	case "xkcd":
		slackClient.SendMessage(slackClient.NewOutgoingMessage(getImage(), slackChannel))
		return
	default:
		slackClient.SendMessage(slackClient.NewOutgoingMessage("Silence Human", slackChannel))

	}
	println("[RECEIVED] sendResponse:", args[0])

	// START SLACKBOT CUSTOM CODE
	// ===============================================================
	// TODO:
	//      1. Implement sendResponse for one or more of your custom Slackbot commands.
	//         You could call an external API here, or create your own string response. Anything goes!
	//      2. STRETCH: Write a goroutine that calls an external API based on the data received in this function.
	// ===============================================================
	// END SLACKBOT CUSTOM CODE
}

// get makes an api request to get weather for a location
func getImage() string {
	i := rand.Intn(1999)

	// make request
	res, _ := http.Get(fmt.Sprintf("https://xkcd.com/%d/info.0.json", i+1))

	// read all into body
	body, _ := ioutil.ReadAll(res.Body)

	// declare struct
	type response struct {
		Image string `json:"img"`
	}
	data := response{}

	// marshal into struct
	json.Unmarshal(body, &data)

	// return image url
	return data.Image

}
