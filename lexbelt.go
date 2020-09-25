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

	putSlotTypeCommand        = app.Command("put-slot-type", "Adds or updates a slot type")
	putSlotTypeCommandName    = putSlotTypeCommand.Flag("name", "Name of Slot Type").String()
	putSlotTypeCommandPublish = putSlotTypeCommand.Flag("publish", "Publish a new version of the slot").Default("false").Bool()
	putSlotTypeCommandFile    = putSlotTypeCommand.Arg("file", "The input specification in json or yaml").Required().File()

	putIntentCommand        = app.Command("put-intent", "Adds or updates an intent")
	putIntentCommandName    = putIntentCommand.Flag("name", "Name of Intent").String()
	putIntentCommandPublish = putIntentCommand.Flag("publish", "Publish a new version of the intent").Default("false").Bool()
	putIntentCommandFile    = putIntentCommand.Arg("file", "The input specification in json or yaml").Required().File()

	putBotCommand        = app.Command("put-bot", "Adds or updates a bot")
	putBotCommandName    = putBotCommand.Flag("name", "Name of Bot").String()
	putBotCommandFile    = putBotCommand.Arg("file", "The input specification in json or yaml").Required().File()
	putBotCommandPublish = putBotCommand.Flag("publish", "Publish a new version the bot with the provided alias").String()
	putBotCommandPoll    = putBotCommand.Flag("poll", "Poll time").Default("3").Int()
	dontWait             = putBotCommand.Flag("dont-wait", "Don't wait for the build to completed before exiting").Default("false").Bool()

	exportCommand     = app.Command("export", "Export existing AWS Lex configs and write to workspace")
	exportCommandJson = exportCommand.Flag("json", "export as json").Default("false").Bool()
)

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
		putSlotType(svc, (*putSlotTypeCommandFile).Name(), *putSlotTypeCommandPublish)
	case putIntentCommand.FullCommand():
		putIntent(svc, (*putIntentCommandFile).Name(), *putIntentCommandPublish)
	case putBotCommand.FullCommand():
		putBot(svc, (*putBotCommandFile).Name(), *putBotCommandName, *putBotCommandPoll, putBotCommandPublish)
	case exportCommand.FullCommand():
		exportLex(svc, *exportCommandJson)
	}

}
