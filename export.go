package main

import (
	"os"

	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

func exportLex(svc *lexmodelbuildingservice.LexModelBuildingService, toJson bool) {

	//ensure the needed child dirs are present

	os.Mkdir("slots", 0744)
	os.Mkdir("intents", 0744)
	os.Mkdir("bots", 0744)

	//get slots
	err := svc.GetSlotTypesPages(&lexmodelbuildingservice.GetSlotTypesInput{}, func(page *lexmodelbuildingservice.GetSlotTypesOutput, lastPage bool) bool {

		for _, v := range page.SlotTypes {
			result, err := svc.GetSlotType(&lexmodelbuildingservice.GetSlotTypeInput{
				Name:    v.Name,
				Version: v.Version,
			})
			checkError(err)
			marshalAndWrite("./slots/"+*result.Name, toJson, result)

		}
		return true
	})
	checkError(err)

	//get Intents
	err = svc.GetIntentsPages(&lexmodelbuildingservice.GetIntentsInput{}, func(page *lexmodelbuildingservice.GetIntentsOutput, lastPage bool) bool {

		for _, v := range page.Intents {
			result, err := svc.GetIntent(&lexmodelbuildingservice.GetIntentInput{
				Name:    v.Name,
				Version: v.Version,
			})
			checkError(err)
			marshalAndWrite("./intents/"+*result.Name, toJson, result)

		}
		return true
	})
	checkError(err)

	//get bots
	err = svc.GetBotsPages(&lexmodelbuildingservice.GetBotsInput{}, func(page *lexmodelbuildingservice.GetBotsOutput, lastPage bool) bool {

		for _, v := range page.Bots {
			result, err := svc.GetBot(&lexmodelbuildingservice.GetBotInput{
				Name:           v.Name,
				VersionOrAlias: v.Version,
			})
			checkError(err)
			marshalAndWrite("./bots/"+*result.Name, toJson, result)

		}
		return true
	})
	checkError(err)
}
