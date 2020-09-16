package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

type provisioner struct {
	LexBot    lexmodelbuildingservice.PutBotInput        `locationName:"lexBot" type:"structure"`
	LexIntent lexmodelbuildingservice.PutIntentInput     `locationName:"lexIntent" type:"structure"`
	LexSlot   []lexmodelbuildingservice.PutSlotTypeInput `locationName:"lexISlot" type:"structure"`
}

type provisionerSpecification struct {
	LexBotProvisioner struct {
		LexBotName *string   `locationName:"lexBotName" type:"string"`
		LexBot     *string   `locationName:"lexBot" type:"string"`
		LexIntent  *string   `locationName:"lexIntent" type:"string"`
		LexSlot    []*string `locationName:"lexSlot" type:"string"`
	} `locationName:"lexBotProvisioner" type:"structure"`
}

func provision(svc *lexmodelbuildingservice.LexModelBuildingService, file string, poll int) {

	var myProvisionerSpecification provisionerSpecification
	readAndUnmarshal(file, &myProvisionerSpecification)

	if myProvisionerSpecification.LexBotProvisioner.LexBot == nil || myProvisionerSpecification.LexBotProvisioner.LexIntent == nil {
		log.Fatal("Yaml file is not as expected, please check your syntax.  You must specify at least a bot and intent")
	}

	var myProvisioner provisioner
	if *provisionCommandName != "" {
		myProvisioner.LexBot.Name = putIntentCommandName
	} else {
		myProvisioner.LexBot.Name = myProvisionerSpecification.LexBotProvisioner.LexBotName
	}

	//based on the mono bot yaml, load slots, intents and the bot and provision in the correct order

	//Slots first.  Slots are optional
	separator := string(os.PathSeparator)
	basePath := filepath.Dir(file) + separator + ".." + separator
	if len(myProvisionerSpecification.LexBotProvisioner.LexSlot) > 0 {
		fmt.Println("Adding Slots")
		//get the slot file.
		for _, v := range myProvisionerSpecification.LexBotProvisioner.LexSlot {
			putSlotType(svc, basePath+"slots"+separator+*v)
		}
	}

	//next intent
	fmt.Println("Adding the intent")
	putIntent(svc, basePath+"intents"+separator+*myProvisionerSpecification.LexBotProvisioner.LexIntent)

	//lastly the bot
	fmt.Println("Adding the bot and building it")
	putBot(svc, basePath+"bots"+separator+*myProvisionerSpecification.LexBotProvisioner.LexBot, poll)

}
