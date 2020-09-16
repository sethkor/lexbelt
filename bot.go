package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

type botYaml struct {
	LexBot *lexmodelbuildingservice.PutBotInput `locationName:"lexBot" type:"structure"`
}

func putBot(svc *lexmodelbuildingservice.LexModelBuildingService) {

	var myBot botYaml
	readAndUnmarshal((*putBotCommandFile).Name(), &myBot)

	if myBot.LexBot == nil {
		log.Fatal("Yaml file is not as expected, please check your syntax.")
	}

	if *putBotCommandName != "" {
		myBot.LexBot.Name = putBotCommandName
	}

	getResult, err := svc.GetBot(&lexmodelbuildingservice.GetBotInput{
		Name:           myBot.LexBot.Name,
		VersionOrAlias: aws.String("$LATEST"),
	})

	checkError(err)

	myBot.LexBot.Checksum = getResult.Checksum

	putResult, err := svc.PutBot(myBot.LexBot)

	checkError(err)

	//loop and poll the status
	if !*dontWait {
		currentStatus := *putResult.Status
		fmt.Print(currentStatus)
		for {

			if currentStatus == "READY" {
				fmt.Println()
				break
			} else if currentStatus == "FAILED" {
				fmt.Printf("\n%s\n", *getResult.FailureReason)
				break
			}

			time.Sleep((time.Duration(*poll) * time.Second))

			getResult, err = svc.GetBot(&lexmodelbuildingservice.GetBotInput{
				Name:           myBot.LexBot.Name,
				VersionOrAlias: aws.String("$LATEST"),
			})

			if currentStatus != *getResult.Status {
				currentStatus = *getResult.Status
				fmt.Printf("\n" + currentStatus)
			} else {
				fmt.Print(".")
			}
		}
	}
}
