package main

import (
	gen "gepaplexx/day-x-generator/pkg/generator"
	sealedSecrets "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
	"log"
	"os"

	"io/ioutil"
)

func readYamlConfiguration(path string) ([]byte, error) {
	config, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func run() {
	args := os.Args[1:]
	utils.PrintDescription(args[0]) // TODO check arguments
	utils.WaitToContinue()

	config, err := readYamlConfiguration(args[0])
	if err != nil {
		log.Fatal("Cannot read config file: ", err)
		return
	}

	// TODO check prerequisites

	// TODO target aus config lesen => default = generated
	err = os.MkdirAll("generated", os.ModePerm)
	if err != nil {
		log.Fatal("cannot create dir: ", err)
		return
	}

	utils.PrintActionHeader("GENERATE CERTIFICATE FOR SEALED SECRETS")
	env, err := utils.FindValue(config, "env")
	if err != nil {
		log.Fatal("cannot find env paramter in config")
		return
	}
	_, _, err = sealedSecrets.GenerateCertificate(env.(string)) // TODO am schluss ausgeben
	if err != nil {
		log.Fatal(err)
		return
	}

	err = gen.Process(config, gen.GENERATORS)
	if err != nil {
		panic(err)
	}

	// TODO erst ganz am Ende augeben inkl. Anleitung was zu tun ist
	// fmt.Println(vaultCertificate)
	// fmt.Println(vaultPrivateKey)
}

func main() {
	run()
}
