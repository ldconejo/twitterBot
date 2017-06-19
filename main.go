package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"os/signal"
	"syscall"
	// flag is imported to support command-line flags
	"flag"
	"regexp"

	// These two libraries had to be installed from the github repositories
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"reflect"
)

////////////////////////////////////
// Implementation functions
///////////////////////////////////

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

// Parses arguments and returns a map
func get_commandline_args() map[string]string {
	// Declare command line parameters and their default values

	// The master account has control over what the bot does
	masterNamePtr := flag.String("master", "luisdconejo", "the name of the master account")
	// myname is the screen name of the servant account
	mynamePtr := flag.String("servant", "eran_marno", "screen name of the servant account")

	flag.Parse()

	// Create an empty map (similar to a Python dictionary)
	cmdLineArgs := map[string]string{}
	cmdLineArgs["masterName"] = *masterNamePtr
	cmdLineArgs["servantName"] = *mynamePtr
	return cmdLineArgs
}

// Decodes direct messages to master account
// Returns a string and an array containing the type of message
// and any parameters
func decodeMasterMessage(masterMessage string) (string, string, string) {
	var validCommand = regexp.MustCompile(`([A-Z]{3}) (.*)`)
	var result string
	var commandArray []string

	if validCommand.MatchString(masterMessage){
		commandArray = validCommand.FindStringSubmatch(masterMessage)
		command := commandArray[1]
		commandParameters := commandArray[2]

		fmt.Println("Command:", command)
		fmt.Println("Parameters:", commandParameters)

		result = "true"
		return result, command, commandParameters
	}
	result = "false"
	return result, "empty", "empty"

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

// Checks if a command is valid
// Valid commands are: TWT (tweet), FLW (follow)

// Launches the bot
func configure(){
	// Get commandline arguments
	cmdLineArgs := get_commandline_args()
	master := cmdLineArgs["masterName"]
	servant := cmdLineArgs["servantName"]

	var twitterKeys []string
	twitterKeys = processKeyFile("keys.txt")

	// Pass in your consumer key (API key) and your consumer secret (API secret)
	config := oauth1.NewConfig(twitterKeys[0], twitterKeys[1])

	// Pass in your access token and your access token secret
	token := oauth1.NewToken(twitterKeys[2], twitterKeys[3])

	// NoContext is the default for most cases
	httpClient := config.Client(oauth1.NoContext, token)
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
			result, command, commandParameters := decodeMasterMessage(dm.Text)

			if  result == "true" {
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
				}
				// RPT stands for report, depending on parameters, provides a report on new followers
			} else {
				fmt.Println( "Received message WAS NOT an order")
				fmt.Println(reflect.TypeOf(client))
			}
		}
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

/* Unused code, kept here as backup
	/*
	demux.All = func(message interface{}){
		fmt.Println(message)
	}
	*/

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
