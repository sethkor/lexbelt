package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

func putBot(svc *lexmodelbuildingservice.LexModelBuildingService, file string, name string, poll int, publish *string) {

	var myBot lexmodelbuildingservice.PutBotInput
	readAndUnmarshal(file, &myBot)

	if publish != nil {
		//publish each intent
		separator := string(os.PathSeparator)
		basePath := filepath.Dir(file) + separator + ".." + separator
		for _, v := range myBot.Intents {
			intentVersion, err := putIntent(svc, basePath+"intents"+separator+*v.IntentName+filepath.Ext(file), true)
			checkError(err)
			v.IntentVersion = aws.String(intentVersion)
		}
	}

	if name != "" {
		myBot.Name = &name
	}

	getResult, err := svc.GetBot(&lexmodelbuildingservice.GetBotInput{
		Name:           myBot.Name,
		VersionOrAlias: aws.String("$LATEST"),
	})

	checkError(err)

	myBot.Checksum = getResult.Checksum

	if publish != nil {
		myBot.CreateVersion = aws.Bool(true)
	}

	putResult, err := svc.PutBot(&myBot)

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
					BotName: myBot.Name,
				})

				checkError(err)
				_, err = svc.PutBotAlias(&lexmodelbuildingservice.PutBotAliasInput{
					Name:       publish,
					BotName:    myBot.Name,
					BotVersion: putResult.Version,
					Checksum:   getAliasResult.Checksum,
				})

				checkError(err)
				fmt.Printf("Bot %s was published as version %s and alias \"%s\".\n",
					*myBot.Name,
					*putResult.Version,
					*publish)

				break
			} else if currentStatus == "FAILED" {
				fmt.Printf("\n%s\n", *getResult.FailureReason)
				break
			}

			time.Sleep(time.Duration(poll) * time.Second)

			getResult, err = svc.GetBot(&lexmodelbuildingservice.GetBotInput{
				Name:           myBot.Name,
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
