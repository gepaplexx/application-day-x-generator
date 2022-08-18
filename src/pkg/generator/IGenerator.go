package generator

import (
	"fmt"
	utils "gepaplexx/day-x-generator/pkg/util"

	"os"
	"text/template"
)

type IValueBuilder interface {
	GetValues(config map[string]utils.Value) (map[string]utils.Value, error)
}

type Stage string

const (
	InitialClusterSetup    Stage = "initial-cluster-setup"
	ClusterSetupCheckpoint Stage = "cluster-setup-checkpoint"
	ClusterApplications    Stage = "cluster-applications"
	AlertManagerConfig     Stage = "default-alertmanager-config"
)

type Generator struct {
	ValueBuilder IValueBuilder
	Stage        Stage
	Name         string
}

type Value = utils.Value

func Process(config []byte, generators []Generator) error {
	initialClusterSetupVals := make(map[string]utils.Value)
	clusterSetupCheckpointVals := make(map[string]utils.Value)
	clusterApplicationsVals := make(map[string]utils.Value)
	alertManagerConfigVals := make(map[string]utils.Value)

	conf, err := utils.FindValuesFlatMap(config, "env")
	if err != nil {
		return err
	}

	initialClusterSetupVals["env"] = conf["env"]
	clusterSetupCheckpointVals["env"] = conf["env"]
	clusterApplicationsVals["env"] = conf["env"]

	utils.PrintActionHeader("BUILD VALUE YAML FILES")
	for _, gen := range generators {
		conf, err = utils.FindValuesFlatMap(config, "env", fmt.Sprintf("%s.%s", gen.Stage, gen.Name))
		if err != nil {
			return err
		}

		utils.PrintAction("Get Values for " + gen.Name)
		res, err := gen.ValueBuilder.GetValues(conf)
		if err != nil {
			utils.PrintFailure()
			return err
		}
		utils.PrintSuccess()

		switch gen.Stage {
		case InitialClusterSetup:
			for k, v := range res {
				initialClusterSetupVals[k] = v
			}
		case ClusterSetupCheckpoint:
			for k, v := range res {
				clusterSetupCheckpointVals[k] = v
			}
		case ClusterApplications:
			for k, v := range res {
				clusterApplicationsVals[k] = v
			}
		case AlertManagerConfig:
			for k, v := range res {
				alertManagerConfigVals[k] = v
			}
		}
	}

	if len(initialClusterSetupVals) != 0 {
		err = executeAndWriteTemplate(InitialClusterSetup, initialClusterSetupVals)
		if err != nil {
			return err
		}
	}

	if len(clusterSetupCheckpointVals) != 0 {
		err = executeAndWriteTemplate(ClusterSetupCheckpoint, clusterSetupCheckpointVals)
		if err != nil {
			return err
		}
	}

	if len(clusterApplicationsVals) != 0 {
		err = executeAndWriteTemplate(ClusterApplications, clusterApplicationsVals)
		if err != nil {
			return err
		}
	}

	if len(clusterApplicationsVals) != 0 {
		err = executeAndWriteTemplate(AlertManagerConfig, alertManagerConfigVals)
		if err != nil {
			return err
		}
	}

	return nil
}

func executeAndWriteTemplate(stage Stage, values map[string]utils.Value) error {
	dest, err := createFile(fmt.Sprintf("generated/%s-%s.yaml", values["env"], stage))
	if err != nil {
		return err
	}

	destTemplPath := fmt.Sprintf("templates/%s.yaml.tpl", stage)
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
