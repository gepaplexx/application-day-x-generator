package main

import (
	gen "gepaplexx/day-x-generator/pkg/generator"
	utils "gepaplexx/day-x-generator/pkg/util"
	"log"
	"os"
)

func readYamlConfiguration(path string) ([]byte, error) {
	config, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func run(configFile string) {
	utils.PrintDescription(configFile)
	utils.WaitToContinue()

	config, err := readYamlConfiguration(configFile)
	if err != nil {
		log.Fatal("Cannot read config file: ", err)
		return
	}

	err = os.MkdirAll(utils.TARGET_DIR, os.ModePerm)
	if err != nil {
		log.Fatal("cannot create dir: ", err)
		return
	}

	err = gen.Process(config, gen.GENERATORS)
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("No config file was specified...")
	}

	run(args[0])
}
