package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

type intentYaml struct {
	LexIntent *lexmodelbuildingservice.PutIntentInput `locationName:"lexIntent" type:"structure"`
}

func putIntent(svc *lexmodelbuildingservice.LexModelBuildingService, file string) {

	var myIntent intentYaml
	readAndUnmarshal(file, &myIntent)

	if myIntent.LexIntent == nil {
		log.Fatal("Yaml file is not as expected, please check your syntax.")
	}
	if *putIntentCommandName != "" {
		myIntent.LexIntent.Name = putIntentCommandName
	}

	getResult, err := svc.GetIntent(&lexmodelbuildingservice.GetIntentInput{
		Name:    myIntent.LexIntent.Name,
		Version: aws.String("$LATEST"),
	})

	checkError(err)

	myIntent.LexIntent.Checksum = getResult.Checksum

	_, err = svc.PutIntent(myIntent.LexIntent)

	checkError(err)

}
