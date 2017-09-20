package pkg

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

// DecodeMasterMessage decodes direct messages to master account
// Returns a string and an array containing the type of message
// and any parameters
func DecodeMasterMessage(masterMessage string) (bool, string, string) {
	var validCommand = regexp.MustCompile(`([A-Z]{3}) (.*)`)
	var result bool
	var commandArray []string

	if validCommand.MatchString(masterMessage) {
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

// ActOnMasterMessage
// Uses the output from decodeMasterMessage and takes action on it
func ActOnMasterMessage(client *twitter.Client, master string, servant string, retweetCandidateMap map[int]*twitter.Tweet, result bool, command string, commandParameters string, pauseRetweet *bool) {
	if result == true {
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
			// Convert to integer and retweet if it exists
			if tweetIndex, err := strconv.Atoi(commandParameters); err == nil {
				fmt.Println("Retweet order received for index " + commandParameters)
				//fmt.Println(retweetCandidateMap[tweetIndex].Text)
				// Check that the index exists
				if retweetCandidate, ok := retweetCandidateMap[tweetIndex]; ok {
					client.Statuses.Update(retweetCandidate.Text, nil)
				} else {
					fmt.Println("ERROR: Incorrect index for retweet")
					SendDirectMessage(client, master, "ERROR: Incorrect index for retweet")
				}
			} else {
				fmt.Println(err)
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

			for _, user := range listOfUsers.Users {
				response = response + "\n" + "@" + user.ScreenName
				followers = append(followers, user.ScreenName)
			}
			// Replies with all followers
			switch commandParameters {
			case "ALL":
				fmt.Println(response)

				SendDirectMessage(client, master, response)
			case "NEW":
				// Check if followers.txt file exists
				if _, err := os.Stat("followers.txt"); err == nil {
					originalList := ProcessKeyFile("followers.txt")
					delta := CompareSlices(followers, originalList)
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
				WriteTextFile("followers.txt", followers)
				SendDirectMessage(client, master, response)
			}
		// Command is to return list of users that are being followed
		case "LST":
			userParams := &twitter.FriendListParams{
				ScreenName: servant,
			}
			listOfUsers, httpResponse, err := client.Friends.List(userParams)
			if err != nil {
				fmt.Println(listOfUsers, httpResponse, err)
			}

			var response string = ""

			var friends []string

			for _, user := range listOfUsers.Users {
				response = response + "\n" + "@" + user.ScreenName
				friends = append(friends, user.ScreenName)
			}
			SendDirectMessage(client, master, response)
		// Command is to follow a user
		case "FLW":
			newFriendParams := &twitter.FriendshipCreateParams{
				ScreenName: commandParameters,
			}
			client.Friendships.Create(newFriendParams)
			response := "Following " + commandParameters
			SendDirectMessage(client, master, response)

		// Command is to unfollow a user
		case "UFL":
			unfriendParams := &twitter.FriendshipDestroyParams{
				ScreenName: commandParameters,
			}
			client.Friendships.Destroy(unfriendParams)
			response := "I stopped following " + commandParameters
			SendDirectMessage(client, master, response)

		// Command to pause retweet requests
		case "PRT":
			switch commandParameters {
			case "YES":
				*pauseRetweet = true
				response := "Pausing all retweet requests"
				SendDirectMessage(client, master, response)
			case "NO":
				*pauseRetweet = false
				response := "Resuming retweet requests"
				SendDirectMessage(client, master, response)
			default:
				if *pauseRetweet == true {
					response := "Retweet requests are currently paused"
					SendDirectMessage(client, master, response)
				} else {
					response := "Retweet requests are currently active"
					SendDirectMessage(client, master, response)
				}
			}
		}

	} else {
		fmt.Println("Received message WAS NOT an order")
	}
}

// Sends a direct message to the master account
func SendDirectMessage(client *twitter.Client, screenName string, messageText string) {
	dmParams := &twitter.DirectMessageNewParams{
		ScreenName: screenName,
		Text:       messageText,
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
	filterWords = ProcessKeyFile("filters.txt")

	filterRegex := ""

	// Assemble combined string for regex
	for _, word := range filterWords {
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
	if interestingTweet.MatchString(tweetText) {
		commandArray := interestingTweet.FindStringSubmatch(tweetText)
		command := commandArray[0]
		fmt.Println("INFO: String match: " + command)
		return true
	}

	return false
}
