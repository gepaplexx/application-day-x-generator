package main

import (
	gen "gepaplexx/day-x-generator/pkg/generator"
	"gepaplexx/day-x-generator/pkg/sealedSecrets"
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

func prerequisitesMet() bool {
	utils.PrintAction("kubeseal")
	installed := utils.IsCommandAvailable("kubeseal")
	if !installed {
		utils.PrintFailure()
		return false
	}
	utils.PrintSuccess()

	return true
}

func run(configFile string) {
	utils.PrintDescription(configFile)
	utils.WaitToContinue()

	config, err := readYamlConfiguration(configFile)
	if err != nil {
		log.Fatal("Cannot read config file: ", err)
		return
	}

	utils.PrintActionHeader("CHECK PREREQUISITES")
	met := prerequisitesMet()
	if !met {
		log.Fatal("Prerequisites not met")
		return
	}
	//TODO: check if required parameters are set in config or fail with error message!

	err = os.MkdirAll(utils.TARGET_DIR, os.ModePerm)
	if err != nil {
		log.Fatal("cannot create dir: ", err)
		return
	}

	err = os.MkdirAll(utils.DEBUG_DIR, os.ModePerm)
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
	// utils.WaitToContinue()
	// fmt.Println(vaultCertificate)
	// fmt.Println(vaultPrivateKey)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("No config file was specified...")
	}

	run(args[0])
}
