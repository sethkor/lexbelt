package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

type slotYaml struct {
	LexSlot *lexmodelbuildingservice.PutSlotTypeInput `locationName:"lexSlot" type:"structure"`
}

//func publishSlot(svc *lexmodelbuildingservice.LexModelBuildingService, name string) (string, error) {
//
//	getResult, err := svc.CreateSlotTypeVersion(&lexmodelbuildingservice.CreateSlotTypeVersionInput{
//		Name: aws.String(name),
//	})
//
//	var slotVersion string
//	if err == nil {
//		slotVersion = *getResult.Version
//	}
//
//	return slotVersion, err
//}

func putSlotType(svc *lexmodelbuildingservice.LexModelBuildingService, file string, publish bool) (string, error) {
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
	mySlotType.LexSlot.CreateVersion = aws.Bool(publish)
	result, err := svc.PutSlotType(mySlotType.LexSlot)

	checkError(err)
	fmt.Printf("Slot %s was published as version %s\n", *mySlotType.LexSlot.Name, *result.Version)
	return *result.Version, err
}
