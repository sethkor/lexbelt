package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

func putIntent(svc *lexmodelbuildingservice.LexModelBuildingService, file string, publish bool) (string, error) {

	var myIntent lexmodelbuildingservice.PutIntentInput
	readAndUnmarshal(file, &myIntent)

	if publish {
		//For publishing, check the intent to see if it has any slots custom slots.  If we find a custom slot and the
		//version specified is latest, we must publish the slot and then update this intent put operation with the
		//version
		separator := string(os.PathSeparator)
		basePath := filepath.Dir(file) + separator + ".." + separator
		for _, v := range myIntent.Slots {
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
		myIntent.Name = putIntentCommandName
	}

	getResult, err := svc.GetIntent(&lexmodelbuildingservice.GetIntentInput{
		Name:    myIntent.Name,
		Version: aws.String(latestVersion),
	})

	checkError(err)
	myIntent.Checksum = getResult.Checksum

	myIntent.CreateVersion = aws.Bool(publish)
	result, err := svc.PutIntent(&myIntent)
	fmt.Printf("Intent %s was published as version %s\n", *myIntent.Name, *result.Version)
	checkError(err)
	return *result.Version, err

}
