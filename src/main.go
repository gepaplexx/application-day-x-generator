package main

import (
	"flag"
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

	if utils.GetConfig().GetDebugSealedSecrets() {
		err = os.MkdirAll("generated/debug/", os.ModePerm)
	} else {
		err = os.MkdirAll("generated", os.ModePerm)
	}

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
	debugSealedSecrets := flag.Bool("debugSealedSecrets", false, "Write generated secrets before sealing")
	clusterConfig := flag.String("clusterConfig", "", "Cluster Configuration")
	flag.Parse()

	if *clusterConfig == "" {
		panic("No config file was specified...")
	}

	config := utils.GetConfig()
	config.SetDebugSealedSecrets(*debugSealedSecrets)

	run(*clusterConfig)
}
