package main

import (
	"fmt"
	"strings"
	"regexp"
	"os"

	"github.com/dghubble/go-twitter/twitter"
)

// Decodes direct messages to master account
// Returns a string and an array containing the type of message
// and any parameters
func decodeMasterMessage(masterMessage string) (bool, string, string) {
	var validCommand = regexp.MustCompile(`([A-Z]{3}) (.*)`)
	var result bool
	var commandArray []string

	if validCommand.MatchString(masterMessage){
		commandArray = validCommand.FindStringSubmatch(masterMessage)
		command := commandArray[1]
		commandParameters := commandArray[2]

		fmt.Println("Command:", command)
		fmt.Println("Parameters:", commandParameters)

		result = true
		return result, command, commandParameters
	}
	result = true
	return result, "empty", "empty"

}

// Act on master message
// Uses the output from decodeMasterMessage and takes action on it
func ActOnMasterMessage (client *twitter.Client, master string, servant string, retweetCandidate *twitter.Tweet, result bool, command string, commandParameters string )  {
	if  result == true {
		switch command {
		// Command is to tweet something
		case "TWT":
			// Tweets the command parameters and sends confirmation to master account
			client.Statuses.Update(commandParameters, nil)
			SendDirectMessage(client, master, "I have posted your tweet.")
			// Command is just to confirm that the bot is active
			// AYT means "Are You There?
		case "AYT":
			// Set parameters for a response via direct message
			SendDirectMessage(client, master, "Hey, boss. I'm active. What's up?")
			// RTW is confirmation or denial of retweeting request
		case "RTW":
			// If reply is yes, then retweet
			if commandParameters == "YES" {
				client.Statuses.Update(retweetCandidate.Text, nil)
			}
			// Request for a list of followers
		case "FLS":
			userParams := &twitter.FollowerListParams{
				ScreenName: servant,
			}
			listOfUsers, httpResponse, err := client.Followers.List(userParams)
			if err != nil {
				fmt.Println(listOfUsers, httpResponse, err)
			}

			var response string = ""

			var followers []string

			for _,user := range listOfUsers.Users{
				response = response + "\n" + "@" + user.ScreenName
				followers = append(followers, user.ScreenName)
			}
			// Replies with all followers
			switch commandParameters {
			case "ALL":
				fmt.Println(response)

				SendDirectMessage(client, master, response )
			case "NEW":
				// Check if followers.txt file exists
				if _, err := os.Stat("followers.txt"); err == nil {
					originalList := processKeyFile("followers.txt")
					delta := compareSlices(followers, originalList)
					response = ""
					for _, screenName := range delta {
						response = response + "\n@" + screenName
					}
				} else {
					response := "No previous followers file. Send FLS ALL to create a new one."
					fmt.Println(response)
				}
				if response == "" {
					response = "No new followers"
				}
				// Save list of followers
				writeTextFile("followers.txt", followers)
				SendDirectMessage(client, master, response )
			}
		}
		// RPT stands for report, depending on parameters, provides a report on new followers
	} else {
		fmt.Println( "Received message WAS NOT an order")
	}
}

// Sends a direct message to the master account
func SendDirectMessage(client *twitter.Client, screenName string, messageText string){
	dmParams := &twitter.DirectMessageNewParams{
		ScreenName: screenName,
		Text: messageText,
	}
	directMessage, httpResponse, err := client.DirectMessages.New(dmParams)
	if err != nil {
		fmt.Println(directMessage, httpResponse, err)
	}
}

// Examines a tweet and decides whether to ask the master for permission to tweet it
// Uses simple matching to select tweets that content words in the list
func ExamineTweet(tweetText string) bool {

	var filterWords []string
	filterWords = processKeyFile("filters.txt")

	filterRegex := ""

	// Assemble combined string for regex
	for _,word := range filterWords{
		word = strings.ToUpper(word)
		if filterRegex == "" {
			filterRegex = "(" + word + ")+"
		} else {
			filterRegex = filterRegex + "|(" + word + ")+"
		}
	}

	fmt.Println("INFO: Regex: " + filterRegex)

	// Capitalizes the string being tested
	tweetText = strings.ToUpper(tweetText)

	// Now, compile regex
	interestingTweet := regexp.MustCompile(filterRegex)

	// Handle tweets that are possible retweet targets
	if interestingTweet.MatchString(tweetText){
		commandArray := interestingTweet.FindStringSubmatch(tweetText)
		command := commandArray[0]
		fmt.Println("INFO: String match: " + command)
		return true
	}

	return false
}

