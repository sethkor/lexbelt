package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"gopkg.in/alecthomas/kingpin.v2"
)

///Command line flags
var (
	app      = kingpin.New("lexbelt", "AWS Lex CLI utilities")
	pProfile = app.Flag("profile", "AWS credentials/config file profile to use").String()
	pRegion  = app.Flag("region", "AWS region").String()

	putSlotTypeCommand     = app.Command("put-slot-type", "Adds or updates a slot type")
	putSlotTypeCommandName = putSlotTypeCommand.Flag("name", "Name of Slot Type").Required().String()
	putSlotTypeCommandFile = putSlotTypeCommand.Arg("file", "The input specification in json or yaml").Required().File()

	putIntentCommand     = app.Command("put-intent", "Adds or updates an intent")
	putIntentCommandName = putIntentCommand.Flag("name", "Name of Intent").Required().String()
	putIntentCommandFile = putIntentCommand.Arg("file", "The input specification in json or yaml").Required().File()

	putBotCommand     = app.Command("put-bot", "Adds or updates a bot")
	putBotCommandName = putBotCommand.Flag("name", "Name of Bot").Required().String()
	putBotCommandFile = putBotCommand.Arg("file", "The input specification in json or yaml").Required().File()
	putBotCommandPoll = putBotCommand.Flag("poll", "Poll time").Default("3").Int()
	dontWait          = putBotCommand.Flag("dont-wait", "Don't wait for the build to completed before exiting").Default("false").Bool()

	provisionCommand     = app.Command("provision", "Provisions and builds an entire Lex bot including slots, intents and the actual bot")
	provisionCommandName = provisionCommand.Flag("name", "Name of Lex Bot to Provision").String()
	provisionCommandPoll = provisionCommand.Flag("poll", "Poll time").Default("3").Int()
	provisionCommandFile = provisionCommand.Arg("file", "The input specification in json or yaml").Required().File()
)

//version variable which can be overidden at computIntentCommandle time
var (
	version = "dev-local-version"
	commit  = "none"
	date    = "unknown"
)

func getAwsSession() *session.Session {
	var sess *session.Session
	if *pProfile != "" {

		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           *pProfile,
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
				MaxRetries:                    aws.Int(30),
			},
		}))

	} else {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
				MaxRetries:                    aws.Int(30),
			},
		}))
	} //else

	if *pRegion != "" {
		sess.Config.Region = aws.String(*pRegion)
	}
	return sess
}

func main() {

	app.Version(version + " " + date + " " + commit)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	sess := getAwsSession()
	svc := lexmodelbuildingservice.New(sess)

	switch command {
	case putSlotTypeCommand.FullCommand():
		putSlotType(svc, (*putSlotTypeCommandFile).Name())
	case putIntentCommand.FullCommand():
		putIntent(svc, (*putIntentCommandFile).Name())
	case putBotCommand.FullCommand():
		putBot(svc, (*putBotCommandFile).Name(), *putBotCommandPoll)
	case provisionCommand.FullCommand():
		provision(svc, (*provisionCommandFile).Name(), *provisionCommandPoll)
	}

}
