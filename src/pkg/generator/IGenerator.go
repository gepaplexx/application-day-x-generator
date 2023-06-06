package generator

import (
	"fmt"
	utils "gepaplexx/day-x-generator/pkg/util"
	"os"
	"strings"
	"text/template"
)

type IValueBuilder interface {
	GetValues(config map[string]utils.Value) (map[string]utils.Value, error)
}

type Stage string
type Generator struct {
	ValueBuilder IValueBuilder
	Stage        Stage
	Name         string
}

const (
	InitialClusterSetup Stage = "initial-cluster-setup"
	ClusterApplications Stage = "cluster-applications"
	VaultSetupScript    Stage = "vault-setup-script"
	ArgoBootstrap       Stage = "argo-bootstrap"
)

var stages = map[Stage]string{
	InitialClusterSetup: "initial-cluster-setup.yaml.tpl",
	ClusterApplications: "cluster-applications.yaml.tpl",
	VaultSetupScript:    "vault-setup-script.sh.tpl",
	ArgoBootstrap:       "argocd-bootstrap-app.yaml.tpl",
}

func Process(config []byte, generators []Generator) error {
	conf, err := utils.FindValuesFlatMap(config, "env")
	if err != nil {
		return err
	}
	env := conf["env"]

	utils.PrintActionHeader("Generate values yaml for initial-cluster-setup")
	err = processGeneric(config, generators, env, InitialClusterSetup)
	if err != nil {
		return err
	}
	utils.PrintActionHeader("Generate values yaml for cluster-applications")
	err = processGeneric(config, generators, env, ClusterApplications)
	if err != nil {
		return err
	}
	utils.PrintActionHeader("Generate vault setup script")
	err = processGeneric(config, generators, env, VaultSetupScript)
	if err != nil {
		return err
	}
	utils.PrintActionHeader("Generate Bootstrap Applications")
	err = processBootstrapApps(config, generators, env)
	if err != nil {
		return err
	}

	return nil
}

func processGeneric(config []byte, generators []Generator, env utils.Value, stage Stage) error {
	stageGenerators := findAllFor(generators, stage)
	vals := make(map[string]utils.Value)
	vals["env"] = env
	for _, currGen := range stageGenerators {
		currVals, err := utils.FindValuesFlatMap(config, buildSearchPath(currGen))
		if err != nil {
			return err
		}
		utils.PrintAction("Get values for " + currGen.Name)
		finalVals, err := currGen.ValueBuilder.GetValues(currVals)
		if err != nil {
			utils.PrintFailure()
			return err
		}
		utils.PrintSuccess()
		for k, v := range finalVals {
			vals[k] = v
		}
	}

	err := executeAndWriteTemplate(vals, stage, env.String(), "", stages[stage])
	if err != nil {
		return err
	}

	return nil
}

func processBootstrapApps(config []byte, generators []Generator, env utils.Value) error {
	bootstrapGenerators := findAllFor(generators, ArgoBootstrap)
	var apps = make(map[string]map[string]utils.Value)
	for _, currGen := range bootstrapGenerators {
		currVals, _ := utils.FindValuesFlatMap(config, buildSearchPath(currGen))
		utils.PrintAction("Get values for " + currGen.Name)
		finalVals, err := currGen.ValueBuilder.GetValues(currVals)
		if err != nil {
			utils.PrintFailure()
			return err
		}
		utils.PrintSuccess()
		apps[currGen.Name] = finalVals
	}

	for app, vals := range apps {
		err := executeAndWriteTemplate(vals, ArgoBootstrap, env.String(), app, stages[ArgoBootstrap])
		if err != nil {
			return err
		}
	}
	return nil
}

func executeAndWriteTemplate(values map[string]utils.Value, stage Stage, env string, genName string, templateName string) error {
	dest, err := createFile(buildFilename(stage, env, genName, templateName))
	if err != nil {
		return err
	}

	destTemplPath := fmt.Sprintf("templates/%s", templateName)
	utils.PrintAction("Read template '" + destTemplPath + "'")
	destTempl, err := template.ParseFiles(destTemplPath)
	if err != nil {
		utils.PrintFailure()
		fmt.Println(err)
		return err
	}
	utils.PrintSuccess()
	utils.PrintAction("Writing final file...")
	err = destTempl.Execute(dest, values)
	if err != nil {
		utils.PrintFailure()
		fmt.Println(err)
		return err
	}
	utils.PrintSuccess()

	dest.Close()

	return nil
}

func buildFilename(stage Stage, env string, genName string, templateName string) string {
	var prefix = ""
	if stage == ArgoBootstrap {
		prefix = genName
	}

	return fmt.Sprintf("%s/%s%s-%s", utils.TARGET_DIR, prefix, env, strings.TrimSuffix(templateName, ".tpl"))
}

func createFile(path string) (*os.File, error) {
	utils.PrintAction("Creating file for values...")
	file, err := os.Create(path)
	if err != nil {
		utils.PrintFailure()
		return nil, err
	}
	utils.PrintSuccess()

	return file, nil
}

func findAllFor(generators []Generator, stage Stage) []Generator {
	res := []Generator{}
	for _, currGen := range generators {
		if currGen.Stage == stage {
			res = append(res, currGen)
		}
	}

	return res
}

func buildSearchPath(gen Generator) string {
	if gen.Name == "" {
		return fmt.Sprintf("%s", gen.Stage)
	}

	return fmt.Sprintf("%s.%s", gen.Stage, gen.Name)
}
