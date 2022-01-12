package main

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology/common/log"
	"io/ioutil"
	"os"
)

type Config struct {
	RpcAddr     string
	BlockHeight uint32
	PanicHeight []uint32
}

func NewSvrConfig(configFilePath string) (*Config, error) {
	fileContent, err := ReadFile(configFilePath)
	if err != nil {
		log.Errorf("NewSvrConfig: failed, err: %s", err)
		return nil, err
	}
	servConfig := &Config{}
	err = json.Unmarshal(fileContent, servConfig)
	if err != nil {
		log.Errorf("NewSvrConfig: failed, err: %s", err)
		return nil, err
	}
	return servConfig, nil
}

func ReadFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: open file %s error %s", fileName, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Errorf("ReadFile: File %s close error %s", fileName, err)
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: ioutil.ReadAll %s error %s", fileName, err)
	}
	return data, nil
}
