package main

import (
	"fmt"
	"log"
	"io"
	"os"
	"os/signal"
	"syscall"

	// These two libraries had to be installed from the github repositories
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Returns a list of lines in a file
// For now, only used to get the Twitter keys for the account
func processKeyFile(keyFile string) [4]string {
	// Open file
	file, err := os.Open(keyFile)
	if err != nil {
		panic(err)
	}
}

func configure(){
	// Pass in your consumer key (API key) and your consumer secret (API secret)
	config := oauth1.NewConfig("placeholder", "placeholder")

	// Pass in your access token and your access token secret
	token := oauth1.NewToken("placeholder", "placeholder")

	// NoContext is the default for most cases
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Demux seems to be an event handler where 'Tweet' and 'DM' are events
	demux := twitter.NewSwitchDemux()

	// This is where the tweet gets printed
	demux.Tweet = func(tweet *twitter.Tweet){
		fmt.Println(tweet.Text)
		//client.DirectMessages.New(false, "eran_marno", tweet.Text)
	}

	// This one handles direct messages that are received
	demux.DM = func(dm *twitter.DirectMessage){
		fmt.Println("FINALLY: " + dm.Text)

	}

	fmt.Println("Starting stream...")

	// DM params
		//dmParams := &twitter.DirectMessageNewParams

	// Filter
		// StreamFilterParams is a struct type, note that filterParams is really a pointer
		filterParams := &twitter.StreamFilterParams{
			// Note that []string is simply the type for a string slice literal (dynamically sized portion
			// of an array)
			Track:		[]string{"cat"},
			StallWarnings:	twitter.Bool(true),
		}

		stream, err := client.Streams.Filter(filterParams)
		if err != nil {
			log.Fatal(err)
		}

	// Receive messages until stopped or stream quits
		go demux.HandleChan(stream.Messages)

		// Wait for SIGINT and SIGTERM (Hit CTRL-C)
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		log.Println(<-ch)

	// Send a Tweet
	// tweet, resp, err := client.Statuses.Update("just setting up my twttr", nil)

	// Send a DM
	// directMessage, resp, err := client.DirectMessages.New(directMessageParams)
	/**type DirectMessageNewParams struct {
		UserID     int64  `url:"user_id,omitempty"`
		ScreenName string `url:"screen_name,omitempty"`
		Text       string `url:"text"`
	}**/
}

func main() {
	fmt.Println("Go-Twitter Bot v0.01")

	// Launch the bot
	configure()
}

