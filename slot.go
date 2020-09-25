package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

func putSlotType(svc *lexmodelbuildingservice.LexModelBuildingService, file string, publish bool) (string, error) {
	var mySlotType lexmodelbuildingservice.PutSlotTypeInput

	readAndUnmarshal(file, &mySlotType)

	if *putSlotTypeCommandName != "" {
		mySlotType.Name = putSlotTypeCommandName
	}
	getResult, err := svc.GetSlotType(&lexmodelbuildingservice.GetSlotTypeInput{
		Name:    mySlotType.Name,
		Version: aws.String("$LATEST"),
	})

	checkError(err)
	mySlotType.Checksum = getResult.Checksum
	mySlotType.CreateVersion = aws.Bool(publish)
	result, err := svc.PutSlotType(&mySlotType)

	checkError(err)
	fmt.Printf("Slot %s was published as version %s\n", *mySlotType.Name, *result.Version)
	return *result.Version, err
}
