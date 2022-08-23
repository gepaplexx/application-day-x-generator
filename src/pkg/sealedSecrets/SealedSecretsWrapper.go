package sealedSecrets

import (
	"fmt"
	utils "gepaplexx/day-x-generator/pkg/util"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	vault "github.com/sosedoff/ansible-vault-go"

	"bytes"
)

type Value = utils.Value

func SealValues(secret []byte, env utils.Value, keys ...string) (map[string]Value, error) {
	encryptedValues, err := seal(secret, env, keys...)
	if err != nil {
		return nil, err
	}

	return encryptedValues, nil
}

// data 	=> YAML formatted String
// env		=> Target Cluster, eg. play, steppe, ...
// ks... 	=> YAML keys wich should be returned
func seal(data []byte, env utils.Value, keys ...string) (map[string]Value, error) {

	cmd := exec.Command("kubeseal", "--cert", "generated/"+env.String()+".crt", "-o", "yaml")
	cmd.Stdin = bytes.NewReader(data)

	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err, ": ", stderr.String())
		return nil, err
	}

	prefixed := prefix("spec:encryptedData:", keys...)
	res, err := utils.FindValues(out.Bytes(), prefixed...)
	if err != nil {
		return nil, err
	}

	if utils.GetConfig().GetDebugSealedSecrets() {
		writeSecretAndValuesToFile(data, keys...)
	}

	return res, nil
}

func encryptForVault(payload string, vaultPassword string) (string, error) {
	// https://pkg.go.dev/github.com/sosedoff/ansible-vault-go#Encrypt
	// Encrypt secret data
	// vaultEncrypt([DATA TO ENCRYPT], [VAULT PASSWORD])
	str, err := vault.Encrypt(payload, vaultPassword)
	if err != nil {
		return "", err
	}

	return str, nil
}

func prefix(prefix string, keys ...string) []string {
	res := []string{}
	for _, key := range keys {
		res = append(res, prefix+key)
	}
	return res
}

func writeSecretAndValuesToFile(secret []byte, keys ...string) {
	secretName, _ := utils.FindValue(secret, "metadata.name")
	writeSecretToFile(secret, secretName)
	writeValuesToFile(secret, secretName, keys...)
}

func writeSecretToFile(secret []byte, secretName any) {
	filenameSecret := fmt.Sprintf("generated/debug/%s.yaml", secretName)
	err := ioutil.WriteFile(filenameSecret, secret, 0644)
	if err != nil {
		log.Printf("WARNING: Failed to write %s.", secretName)
	}
}

func writeValuesToFile(secret []byte, secretName any, keys ...string) {
	prefixed := prefix("data:", keys...)
	values, _ := utils.FindValues(secret, prefixed...)
	filenameValues := fmt.Sprintf("generated/debug/%s-values.txt", secretName)
	fileValues, err := os.Create(filenameValues)
	if err != nil {
		log.Printf("WARNING: Failed to write %s.", filenameValues)
		return
	}
	defer fileValues.Close()

	for key, val := range values {
		decoded, err := utils.Base64Decode(val)
		if err != nil {
			log.Printf("WARNING: Failed to decode string %s.", val.String())
			continue
		}

		fmt.Fprintf(fileValues, "%s:\n", key)
		fmt.Fprintf(fileValues, "%s\n", decoded)
		fmt.Fprintf(fileValues, "==== END ======\n")
	}
}
