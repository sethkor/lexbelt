package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

type slotYaml struct {
	LexSlot *lexmodelbuildingservice.PutSlotTypeInput `locationName:"lexSlot" type:"structure"`
}

func putSlotType(svc *lexmodelbuildingservice.LexModelBuildingService, file string) {
	var mySlotType slotYaml

	readAndUnmarshal(file, &mySlotType)

	if mySlotType.LexSlot == nil {
		log.Fatal("Yaml file is not as expected, please check your syntax.")
	}

	if *putSlotTypeCommandName != "" {
		mySlotType.LexSlot.Name = putSlotTypeCommandName
	}
	getResult, err := svc.GetSlotType(&lexmodelbuildingservice.GetSlotTypeInput{
		Name:    mySlotType.LexSlot.Name,
		Version: aws.String("$LATEST"),
	})

	checkError(err)

	mySlotType.LexSlot.Checksum = getResult.Checksum

	_, err = svc.PutSlotType(mySlotType.LexSlot)

	checkError(err)

}
