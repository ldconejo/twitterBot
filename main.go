package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"os/signal"
	"syscall"

	// These two libraries had to be installed from the github repositories
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Returns a list of lines in a file
// For now, only used to get the Twitter keys for the account
func processKeyFile(keyFile string) []string {
	// Open file
	file, err := os.Open(keyFile)
	if err != nil {
		panic(err)
	}
	// Defer occurs only after the function ends
	// which makes sense, considering it closes the file
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Gets each line from the file, using the scanner, and appends it to the array
	for scanner.Scan(){
		lines = append(lines, scanner.Text())
	}
	return lines
}

// Launches the bot
func configure(){
	var twitterKeys []string
	twitterKeys = processKeyFile("keys.txt")

	fmt.Println(twitterKeys[0])
	fmt.Println(twitterKeys[1])
	fmt.Println(twitterKeys[2])
	fmt.Println(twitterKeys[3])

	// Pass in your consumer key (API key) and your consumer secret (API secret)
	config := oauth1.NewConfig(twitterKeys[0], twitterKeys[1])

	// Pass in your access token and your access token secret
	token := oauth1.NewToken(twitterKeys[2], twitterKeys[3])

	// NoContext is the default for most cases
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Demux seems to be an event handler where 'Tweet' and 'DM' are events
	demux := twitter.NewSwitchDemux()

	// This is where the tweet gets printed
	/*
	demux.Tweet = func(tweet *twitter.Tweet){
		// Direct message (DM) params
		// This is used to send messages to the master account
		dmParams := &twitter.DirectMessageNewParams{
			ScreenName: "eran_marno",
			Text: tweet.Text,
		}
		fmt.Println(tweet.Text)
		client.DirectMessages.New(dmParams)
	}
	*/

	// This one handles direct messages that are received
	// Part of SwitchDemux
	demux.DM = func(dm *twitter.DirectMessage){
		fmt.Println("FINALLY: " + dm.Text)
		fmt.Println(dm.SenderID)

	}

	/*
	demux.All = func(message interface{}){
		fmt.Println(message)
	}
	*/


	fmt.Println("Starting stream...")

	// Filter
	// StreamFilterParams is a struct type, note that filterParams is really a pointer

	/*
	filterParams := &twitter.StreamFilterParams{
		// Note that []string is simply the type for a string slice literal (dynamically sized portion
		// of an array)
		Track:		[]string{"LUISCONEJO"},
		StallWarnings:	twitter.Bool(true),
	}
	*/

	// Sample stream (instead of filtered)
	/*
	params := &twitter.StreamSampleParams{
		StallWarnings: twitter.Bool(true),
	}
	*/

	// User stream
	userParams := &twitter.StreamUserParams{
		With: "followings",
		StallWarnings: twitter.Bool(true),
	}

	//stream, err := client.Streams.Filter(filterParams)
	//stream, err := client.Streams.Sample(params)
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

func main() {
	fmt.Println("Go-Twitter Bot v0.01")

	// Launch the bot
	configure()
}

