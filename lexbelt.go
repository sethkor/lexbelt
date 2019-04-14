package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

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
	pVerbose = app.Flag("verbose", "Verbose Logging").Default("false").Bool()

	putSlotTypeCommand     = app.Command("put-slot-type", "Adds or updates a slot type")
	putSlotTypeCommandName = putSlotTypeCommand.Flag("name", "Name of Slot Type").Required().String()
	putSlotTypeCommandFile = putSlotTypeCommand.Flag("cli-input-json", "JSON file of Slot Type").Required().URL()

	putIntentCommand     = app.Command("put-intent", "Adds or updates an intent")
	putIntentCommandName = putIntentCommand.Flag("name", "Name of Intent").Required().String()
	putIntentCommandFile = putIntentCommand.Flag("cli-input-json", "JSON file of Intent").Required().URL()

	putBotCommand     = app.Command("put-bot", "Adds or updates a bot")
	putBotCommandName = putBotCommand.Flag("name", "Name of Intent").Required().String()
	putBotCommandFile = putBotCommand.Flag("cli-input-json", "JSON file of Intent").Required().URL()
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

func genericConversion(lexStuct interface{}) interface{} {

	file := (*putIntentCommandFile).Path

	thefile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal("reading config file", err.Error())
	}

	resStruct := reflect.New(reflect.TypeOf(lexStuct))

	switch filepath.Ext(file) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(thefile, &resStruct)

	default:
		err = json.Unmarshal(thefile, &resStruct)

	}

	if err != nil {
		log.Fatal("parsing config file", err.Error())
	} else {
		fmt.Println(resStruct)
	}

	return resStruct
}

func putSlotType(svc *lexmodelbuildingservice.LexModelBuildingService) {

	var mySlotType lexmodelbuildingservice.PutSlotTypeInput

	file := (*putIntentCommandFile).Path

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
	} else {
		fmt.Println(mySlotType)
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

	putResult, err := svc.PutSlotType(&mySlotType)

	checkError(err)
	fmt.Println(putResult)

}

func putIntent(svc *lexmodelbuildingservice.LexModelBuildingService) {

	var myIntent lexmodelbuildingservice.PutIntentInput

	myIntentInterface := genericConversion(myIntent)
	myIntent = myIntentInterface.(lexmodelbuildingservice.PutIntentInput)

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
	} else {
		fmt.Println(myIntent)
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

	putResult, err := svc.PutIntent(&myIntent)

	checkError(err)
	fmt.Println(putResult)

}

func putBot(svc *lexmodelbuildingservice.LexModelBuildingService) {

	var myBot lexmodelbuildingservice.PutBotInput

	file := (*putIntentCommandFile).Path

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
	} else {
		fmt.Println(myBot)
	}

	if *putBotCommandName != "" {
		myBot.Name = putBotCommandName
	}

	getResult, err := svc.GetIntent(&lexmodelbuildingservice.GetIntentInput{
		Name:    myBot.Name,
		Version: aws.String("$LATEST"),
	})

	checkError(err)

	myBot.Checksum = getResult.Checksum

	putResult, err := svc.PutBot(&myBot)

	checkError(err)
	fmt.Println(putResult)

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
