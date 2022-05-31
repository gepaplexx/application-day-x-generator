package sealedSecrets

import (
	utils "gepaplexx/day-x-generator/pkg/util"
	"log"
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
