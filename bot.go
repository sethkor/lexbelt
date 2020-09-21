package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

type botYaml struct {
	LexBot *lexmodelbuildingservice.PutBotInput `locationName:"lexBot" type:"structure"`
}

func publishBot(svc *lexmodelbuildingservice.LexModelBuildingService, name string) string {

	getResult, err := svc.CreateBotVersion(&lexmodelbuildingservice.CreateBotVersionInput{
		Name: aws.String(name),
	})

	checkError(err)

	return *getResult.Version
}

func putBot(svc *lexmodelbuildingservice.LexModelBuildingService, file string, name string, poll int, publish *string) {

	var myBot botYaml
	readAndUnmarshal(file, &myBot)

	if myBot.LexBot == nil {
		log.Fatal("Yaml file is not as expected, please check your syntax.")
	}

	if publish != nil {
		//publish each intent
		separator := string(os.PathSeparator)
		basePath := filepath.Dir(file) + separator + ".." + separator
		for _, v := range myBot.LexBot.Intents {
			intentVersion, err := putIntent(svc, basePath+"intents"+separator+*v.IntentName+filepath.Ext(file), true)
			checkError(err)
			v.IntentVersion = aws.String(intentVersion)
		}
	}

	if name != "" {
		myBot.LexBot.Name = &name
	}

	getResult, err := svc.GetBot(&lexmodelbuildingservice.GetBotInput{
		Name:           myBot.LexBot.Name,
		VersionOrAlias: aws.String("$LATEST"),
	})

	checkError(err)

	myBot.LexBot.Checksum = getResult.Checksum

	if publish != nil {
		myBot.LexBot.CreateVersion = aws.Bool(true)
	}

	putResult, err := svc.PutBot(myBot.LexBot)

	checkError(err)

	//loop and poll the status
	if !*dontWait || publish != nil {
		currentStatus := *putResult.Status
		fmt.Print(currentStatus)
		for {

			if currentStatus == "READY" {
				fmt.Println()

				getAliasResult, err := svc.GetBotAlias(&lexmodelbuildingservice.GetBotAliasInput{
					Name:    publish,
					BotName: myBot.LexBot.Name,
				})

				checkError(err)
				_, err = svc.PutBotAlias(&lexmodelbuildingservice.PutBotAliasInput{
					Name:       publish,
					BotName:    myBot.LexBot.Name,
					BotVersion: putResult.Version,
					Checksum:   getAliasResult.Checksum,
				})

				checkError(err)
				fmt.Printf("Bot %s was published as version %s and alias \"%s\".\n",
					*myBot.LexBot.Name,
					*putResult.Version,
					*publish)

				break
			} else if currentStatus == "FAILED" {
				fmt.Printf("\n%s\n", *getResult.FailureReason)
				break
			}

			time.Sleep(time.Duration(poll) * time.Second)

			getResult, err = svc.GetBot(&lexmodelbuildingservice.GetBotInput{
				Name:           myBot.LexBot.Name,
				VersionOrAlias: aws.String("$LATEST"),
			})

			checkError(err)

			if currentStatus != *getResult.Status {
				currentStatus = *getResult.Status
				fmt.Printf("\n" + currentStatus)
			} else {
				fmt.Print(".")
			}
		}
	}
}
