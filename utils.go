package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/ghodss/yaml"
)

func checkError(err error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {

			case lexmodelbuildingservice.ErrCodeNotFoundException:
				break
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Fatal(err.Error())
		}

	}

}

func readAndUnmarshal(fileName string, destination interface{}) {
	thefile, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal("reading config file", err.Error())
	}

	switch filepath.Ext(fileName) {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(thefile, destination)

	default:
		err = json.Unmarshal(thefile, destination)

	}

	if err != nil {
		log.Fatal("parsing config file", err.Error())
	}

}

func marshalAndWrite(fileName string, toJson bool, source interface{}) {
	var file []byte
	extension := ".yaml"
	if toJson {
		file, _ = json.MarshalIndent(source, "", "  ")
		extension = ".json"
	} else {
		file, _ = yaml.Marshal(source)
	}
	err := ioutil.WriteFile(fileName+extension, file, 0644)
	checkError(err)
}
