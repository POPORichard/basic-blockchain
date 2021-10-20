package server

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const yamlFilePath = "config.yaml"

type conf struct {
	TrustAddress string	`yaml:"TrustAddress"`
	LocalAddress string	`yaml:"LocalAddress"`
}

func (c *conf)getConf() *conf{
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil{
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil{
		panic(err)
	}
	return c
}

func readConfig()(TrustAddress, LocalAddress string){
	var c conf
	c.getConf()
	fmt.Println("读取配置",c)
	return c.TrustAddress, c.LocalAddress
}