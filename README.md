# Secondary Twitter Account Handler Bot

This implementation of a Twitter Bot using Golang allows the user to control a secondary account from a primary account by using direct messages. At this point, you can:

* Send a tweet as a DM from the primary account to be published from the secondary account.
* Have the secondary account look at its user stream for keywords and ask the primary account, via DM, for authorization to retweent.
* Check whether the bot is active.

Further functionality to follow / unfollow will be added next.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

This bot uses the go-twitter client library for the Twitter API, you will need to install it before you can use the Twitter bot.

```
go get github.com/dghubble/go-twitter/twitter
go get github.com/dghubble/oauth1
```
### Installing

To test it, fork a copy of my repository. Once that is done, you will need to create a couple of additional files.

keys.txt - This file should contain four lines with your keys and access tokens for your secondary Twitter account. In order:

* Consumer Key (API Key)
* Consumer Secret (API Secret)
* Access Token
* Access Token Secret

You can create these at https://apps.twitter.com/

Additionally, in your application settings, you will need to set the access level on both "Application Settings" as well as "Your Access Token" to "Read, write, and direct messages". This will allow the app to communicate via DM.

filters.txt - This one will contain all words or phrases that you want the bot to identify as retweet candidates. For example:

```
VR
virtual
falcon 9
artificial intelligence
rocket
developer
```
Both files need to be placed in the same directory as main.go.

Once you are ready, compile and run by typing:

```
go run main.go
```

## Contributing

This is still to be defined, as the project is a very early stage. That said, I welcome any suggestions or improvements to the code.

## Authors

* **Luis Conejo** [ldconejo](https://github.com/ldconejo)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* Dalton Hubble, this project would be non-existent without go-twitter.
