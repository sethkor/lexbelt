package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

type intentYaml struct {
	LexIntent *lexmodelbuildingservice.PutIntentInput `locationName:"lexIntent" type:"structure"`
}

func publishIntent(svc *lexmodelbuildingservice.LexModelBuildingService, name string) string {

	getResult, err := svc.CreateIntentVersion(&lexmodelbuildingservice.CreateIntentVersionInput{
		Name: aws.String(name),
	})

	checkError(err)

	return *getResult.Version
}

func putIntent(svc *lexmodelbuildingservice.LexModelBuildingService, file string, publish bool) (string, error) {

	var myIntent intentYaml
	readAndUnmarshal(file, &myIntent)

	if myIntent.LexIntent == nil {
		log.Fatal("Yaml file is not as expected, please check your syntax.")
	}

	if publish {
		//For publishing, check the intent to see if it has any slots custom slots.  If we find a custom slot and the
		//version specified is latest, we must publish the slot and then update this intent put operation with the
		//version
		separator := string(os.PathSeparator)
		basePath := filepath.Dir(file) + separator + ".." + separator
		for _, v := range myIntent.LexIntent.Slots {
			if v.SlotTypeVersion != nil {
				//assume if no version is present, it's a built in type and we don't need to do anything
				if *v.SlotTypeVersion == latestVersion {
					slotVersion, err := putSlotType(svc, basePath+"slots"+separator+*v.SlotType+filepath.Ext(file), true)
					if err != nil {
						checkError(err)
					}
					*v.SlotTypeVersion = slotVersion
				}
			}
		}
	}

	if *putIntentCommandName != "" {
		myIntent.LexIntent.Name = putIntentCommandName
	}

	getResult, err := svc.GetIntent(&lexmodelbuildingservice.GetIntentInput{
		Name:    myIntent.LexIntent.Name,
		Version: aws.String(latestVersion),
	})

	checkError(err)
	myIntent.LexIntent.Checksum = getResult.Checksum

	myIntent.LexIntent.CreateVersion = aws.Bool(publish)
	result, err := svc.PutIntent(myIntent.LexIntent)
	fmt.Printf("Intent %s was published as version %s\n", *myIntent.LexIntent.Name, *result.Version)
	checkError(err)
	return *result.Version, err

}
