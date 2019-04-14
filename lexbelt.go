package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/ghodss/yaml"
	"gopkg.in/alecthomas/kingpin.v2"
)

///Command line flags
var (
	app      = kingpin.New("lexbelt", "AWS Lex CLI utilities")
	pProfile = app.Flag("profile", "AWS credentials/config file profile to use").String()
	pRegion  = app.Flag("region", "AWS region").String()
	pVerbose = app.Flag("verbose", "Verbose Logging - not implemented yet").Default("false").Bool()

	putSlotTypeCommand     = app.Command("put-slot-type", "Adds or updates a slot type")
	putSlotTypeCommandName = putSlotTypeCommand.Flag("name", "Name of Slot Type").Required().String()
	putSlotTypeCommandFile = putSlotTypeCommand.Flag("cli-input-json", "JSON file of Slot Type").Required().URL()

	putIntentCommand     = app.Command("put-intent", "Adds or updates an intent")
	putIntentCommandName = putIntentCommand.Flag("name", "Name of Intent").Required().String()
	putIntentCommandFile = putIntentCommand.Flag("cli-input-json", "JSON file of Intent").Required().URL()

	putBotCommand     = app.Command("put-bot", "Adds or updates a bot")
	putBotCommandName = putBotCommand.Flag("name", "Name of Intent").Required().String()
	putBotCommandFile = putBotCommand.Flag("cli-input-json", "JSON file of Intent").Required().URL()
	poll              = putBotCommand.Flag("poll", "Poll time").Default("3").Int()
	dontWait          = putBotCommand.Flag("dont-wait", "Don't wait for the build to completed before exiting").Default("false").Bool()
)

//version variable which can be overidden at computIntentCommandle time
var (
	version = "dev-local-version"
	commit  = "none"
	date    = "unknown"
)

func checkError(err error) {
	if err != nil {

		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {

			case lexmodelbuildingservice.ErrCodeNotFoundException:
				break
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Fatal(err.Error())
		}

	}

}

func putSlotType(svc *lexmodelbuildingservice.LexModelBuildingService) {

	var mySlotType lexmodelbuildingservice.PutSlotTypeInput

	file := (*putSlotTypeCommandFile).Path

	thefile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal("reading config file", err.Error())
	}

	switch filepath.Ext(file) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(thefile, &mySlotType)

	default:
		err = json.Unmarshal(thefile, &mySlotType)

	}

	if err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	if *putSlotTypeCommandName != "" {
		mySlotType.Name = putSlotTypeCommandName
	}

	getResult, err := svc.GetSlotType(&lexmodelbuildingservice.GetSlotTypeInput{
		Name:    mySlotType.Name,
		Version: aws.String("$LATEST"),
	})

	checkError(err)

	mySlotType.Checksum = getResult.Checksum

	_, err = svc.PutSlotType(&mySlotType)

	checkError(err)

}

func putIntent(svc *lexmodelbuildingservice.LexModelBuildingService) {

	var myIntent lexmodelbuildingservice.PutIntentInput

	file := (*putIntentCommandFile).Path

	thefile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal("reading config file", err.Error())
	}

	switch filepath.Ext(file) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(thefile, &myIntent)

	default:
		err = json.Unmarshal(thefile, &myIntent)

	}

	if err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	if *putIntentCommandName != "" {
		myIntent.Name = putIntentCommandName
	}

	getResult, err := svc.GetIntent(&lexmodelbuildingservice.GetIntentInput{
		Name:    myIntent.Name,
		Version: aws.String("$LATEST"),
	})

	checkError(err)

	myIntent.Checksum = getResult.Checksum

	_, err = svc.PutIntent(&myIntent)

	checkError(err)

}

func putBot(svc *lexmodelbuildingservice.LexModelBuildingService) {

	var myBot lexmodelbuildingservice.PutBotInput

	file := (*putBotCommandFile).Path

	thefile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("reading config file", err.Error())
	}

	switch filepath.Ext(file) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(thefile, &myBot)

	default:
		err = json.Unmarshal(thefile, &myBot)

	}

	if err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	if *putBotCommandName != "" {
		myBot.Name = putBotCommandName
	}

	getResult, err := svc.GetBot(&lexmodelbuildingservice.GetBotInput{
		Name:           myBot.Name,
		VersionOrAlias: aws.String("$LATEST"),
	})

	checkError(err)

	myBot.Checksum = getResult.Checksum

	putResult, err := svc.PutBot(&myBot)

	checkError(err)

	//loop and poll the status

	if !*dontWait {
		getResult.Status = putResult.Status
		for {

			fmt.Println(*getResult.Status)

			if *getResult.Status == "READY" {
				break
			} else if *getResult.Status == "FAILED" {
				fmt.Println(*getResult.FailureReason)
				break
			}

			time.Sleep((time.Duration(*poll) * time.Second))

			getResult, err = svc.GetBot(&lexmodelbuildingservice.GetBotInput{
				Name:           myBot.Name,
				VersionOrAlias: aws.String("$LATEST"),
			})

		}
	}

}

func main() {

	app.Version(version)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	var sess *session.Session
	if *pProfile != "" {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           *pProfile,
			SharedConfigState: session.SharedConfigEnable,
		}))

	} else {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	} //else
	if *pRegion != "" {
		sess.Config.Region = aws.String(*pRegion)
	}

	svc := lexmodelbuildingservice.New(sess)

	switch command {
	case putSlotTypeCommand.FullCommand():
		putSlotType(svc)
	case putIntentCommand.FullCommand():
		putIntent(svc)
	case putBotCommand.FullCommand():
		putBot(svc)
	}

}
