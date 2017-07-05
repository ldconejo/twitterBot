package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	// These libraries had to be installed from the github repositories
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"twitterBot/pkg"

)

////////////////////////////////////
// Implementation functions
///////////////////////////////////
// Checks if a command is valid
// Valid commands are: TWT (tweet), FLW (follow)

// Launches the bot
func configure(){
	// Get commandline arguments
	cmdLineArgs := pkg.Get_commandline_args()
	master := cmdLineArgs["masterName"]
	servant := cmdLineArgs["servantName"]

	// This variable holds retweet candidates
	var retweetCandidate *twitter.Tweet

	var twitterKeys []string
	twitterKeys = pkg.ProcessKeyFile("keys.txt")

	// Pass in your consumer key (API key) and your consumer secret (API secret)
	config := oauth1.NewConfig(twitterKeys[0], twitterKeys[1])

	// Pass in your access token and your access token secret
	token := oauth1.NewToken(twitterKeys[2], twitterKeys[3])

	// NoContext is the default for most cases
	httpClient := config.Client(oauth1.NoContext, token)

	// Creates the Twitter Client, which wil have services allow you to handle the account
	client := twitter.NewClient(httpClient)

	// Demux seems to be an event handler where 'Tweet' and 'DM' are events
	demux := twitter.NewSwitchDemux()

	// This one handles direct messages that are received
	// Part of SwitchDemux
	demux.DM = func(dm *twitter.DirectMessage){
		// Check if the message comes from the master account
		if dm.SenderScreenName != master && dm.SenderScreenName != servant{
			fmt.Println("Whoever you are, you're not my master")
			fmt.Println(dm.SenderScreenName)
			fmt.Println(master)
		} else if dm.SenderScreenName == master{
			// Decode instruction from master
			result, command, commandParameters := pkg.DecodeMasterMessage(dm.Text)

			// Take action on decoded master message
			pkg.ActOnMasterMessage(client, master, servant, retweetCandidate, result, command, commandParameters)

		}
	}

	// This one handles tweets on the user stream
	demux.Tweet = func(tweet *twitter.Tweet){
		if pkg.ExamineTweet(tweet.Text){
			fmt.Println("This tweet is interesting:" + tweet.Text + "\n")

			// Now, ask the master account for permission to retweet
			pkg.SendDirectMessage(client, master, "RTW CANDIDATE: " + tweet.Text)

			// Saves the candidate
			retweetCandidate = tweet
		}
		fmt.Println(tweet.Text)
	}

	// This one handles notifications of new followers (covered under "Event")
	demux.Event = func(event *twitter.Event) {
		fmt.Println("INFO: New event - " + event.Event)
		fmt.Println("Created at: " + event.CreatedAt)
		fmt.Println("Target: " + event.Target.ScreenName)
		fmt.Println("Source: " + event.Source.ScreenName)
	}

	fmt.Println("Starting stream...")

	// User stream
	userParams := &twitter.StreamUserParams{
		With: "followings",
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.User(userParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (Hit CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()

}

// Main flow, this is where the program launches
func main() {
	fmt.Println("Go-Twitter Master-Servant Bot v0.01")

	// Launch the bot
	configure()
}