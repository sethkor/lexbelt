package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

///Command line flags
var (
	app         = kingpin.New("lexbelt", "AWS Lex CLI utilities")
	pProfile    = app.Flag("profile", "AWS credentials/config file profile to use").String()
	pRegion     = app.Flag("region", "AWS region").String()
	//pVerbose    = app.Flag("verbose", "Verbose Logging").Default("false").Bool()

	pst  			= app.Command("put-slot-type", "Adds or updates a slot type")
	pstName   	= pst.Flag("name", "Name of Slot Type").Required().String()
	pstFile		= pst.Flag("cli-input-json","JSON file of Slot Type").Required().URL()

	pi  			= app.Command("put-intent", "Adds or updates an intent")
	piName   	= pi.Flag("name", "Name of Intent").Required().String()
	piFile		= pi.Flag("cli-input-json","JSON file of Intent").Required().URL()


	pb  			= app.Command("put-bot", "Adds or updates a bot")
	pbName   	= pi.Flag("name", "Name of Intent").Required().String()
	pbFile		= pi.Flag("cli-input-json","JSON file of Intent").Required().URL()
)

//version variable which can be overidden at compile time
var (
	version = "dev-local-version"
	commit  = "none"
	date    = "unknown"
)


func checkError (err error) {
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

func putSlotType (svc *lexmodelbuildingservice.LexModelBuildingService, decoder json.Decoder) {

	//check for existing slot type and get checksump

	var mySlotType lexmodelbuildingservice.PutSlotTypeInput

	if err := decoder.Decode(&mySlotType); err != nil {
		log.Fatal("parsing config file", err.Error())
	} else {
		fmt.Println(mySlotType)
	}

	getResult, err := svc.GetSlotType(&lexmodelbuildingservice.GetSlotTypeInput{
		Name: mySlotType.Name,
		Version: aws.String("$LATEST"),
	})

	checkError(err)

	mySlotType.Checksum = getResult.Checksum
	putResult, err := svc.PutSlotType(&mySlotType)

	checkError(err)
	fmt.Println(putResult)

}

func putIntent (svc *lexmodelbuildingservice.LexModelBuildingService, decoder json.Decoder){

	//check for existing slot type and get checksump

	var myIntent lexmodelbuildingservice.PutIntentInput

	if err := decoder.Decode(&myIntent); err != nil {
		log.Fatal("parsing config file", err.Error())
	} else {
		fmt.Println(myIntent)
	}

	getResult, err := svc.GetIntent(&lexmodelbuildingservice.GetIntentInput{
		Name: myIntent.Name,
		Version: aws.String("$LATEST"),
	})

	checkError(err)

	myIntent.Checksum = getResult.Checksum
	putResult, err := svc.PutIntent(&myIntent)

	checkError(err)
	fmt.Println(putResult)

}

func putBot (svc *lexmodelbuildingservice.LexModelBuildingService, decoder json.Decoder){

	//check for existing slot type and get checksump

	var myBot lexmodelbuildingservice.PutBotInput

	if err := decoder.Decode(&myBot); err != nil {
		log.Fatal("parsing config file", err.Error())
	} else {
		fmt.Println(myBot)
	}

	getResult, err := svc.GetIntent(&lexmodelbuildingservice.GetIntentInput{
		Name: myBot.Name,
		Version: aws.String("$LATEST"),
	})

	checkError(err)

	myBot.Checksum = getResult.Checksum
	putResult, err := svc.PutBot(&myBot)

	checkError(err)
	fmt.Println(putResult)

}


func main () {

	app.Version(version)
	kingpin.CommandLine.HelpFlag.Short('h')

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

	jsonFile, err := os.Open((*pstFile).Path)

	if err != nil {
		fmt.Println("Error opening json file", err.Error())
	}

	jsonParser := json.NewDecoder(jsonFile)

	switch command {
	case pst.FullCommand():
		putSlotType(svc, *jsonParser)
	case pi.FullCommand():
		putIntent(svc, *jsonParser)
	case pb.FullCommand():
		putBot(svc, *jsonParser)
	}

}