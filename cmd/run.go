/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/exfly/bmc"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Manifest struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

type ManifestConfig struct {
	Architecture string `json:"architecture"`
	Config       struct {
		Env        []string    `json:"Env"`
		Cmd        []string    `json:"Cmd"`
		WorkingDir string      `json:"WorkingDir"`
		OnBuild    interface{} `json:"OnBuild"`
	} `json:"config"`
}

var tar = archiver.Tar{
	OverwriteExisting: true,
	MkdirAll:          true,
}

var (
	runFromTar     string
	runToDir       string
	runRootfs      bool
	runContainerID string
	runConfig      string
	runRawCmd      string
	runTerminal    bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		baseDir := runToDir
		fromDir := filepath.Join(baseDir, "origin")
		toDir := filepath.Join(baseDir, "rootfs")

		// STEP: untar docker image
		if runRootfs {
			log.Println("build rootfs")
			err := tar.Unarchive(
				runFromTar,
				fromDir,
			)
			if err != nil {
				return errors.Wrap(err, "")
			}

			err = prepareRootfs(ctx, runTerminal, baseDir, fromDir, toDir, strings.Split(runRawCmd, " "))
			if err != nil {
				return errors.Wrap(err, "")
			}
		}

		runcPath := filepath.Join(baseDir, "runc")
		if !fileExists(runcPath) {
			runcExec, err := bmc.BinDataFs.ReadFile("bin/runc")
			if err != nil {
				return errors.Wrap(err, "")
			}

			err = ioutil.WriteFile(runcPath, runcExec, 0711)
			if err != nil {
				return errors.Wrap(err, "")
			}
		}

		log.Println("run in new namespace")

		baseDirAbs, err := filepath.Abs(baseDir)
		if err != nil {
			return errors.Wrap(err, "")
		}
		err = os.Setenv("PATH", baseDirAbs)
		if err != nil {
			return errors.Wrap(err, "")
		}
		err = execCmdInteract(ctx, baseDir, "runc", "run", runContainerID)
		if err != nil {
			return errors.Wrap(err, "")
		}
		return nil
	},
}

func prepareRootfs(ctx context.Context, terminal bool, baseDir, fromDir, toDir string, firstCmd []string) error {
	manifestFile := "manifest.json"
	manifestConfig := make([]string, 0, 10)

	// STEP: 解压 layer
	content, err := ioutil.ReadFile(filepath.Join(fromDir, manifestFile))
	if err != nil {
		return errors.Wrap(err, "")
	}
	var manifests []Manifest
	err = json.Unmarshal(content, &manifests)
	if err != nil {
		return errors.Wrap(err, "")
	}

	for _, manifest := range manifests {
		manifestConfig = append(manifestConfig, manifest.Config)

		for _, layer := range manifest.Layers {
			log.Printf("untar file %v -> %v", layer, toDir)
			err = tar.Unarchive(filepath.Join(fromDir, layer), toDir)
			if err != nil {
				return errors.Wrap(err, "")
			}
		}
	}

	// STEP: 初始化 hostNetwork 网络配置
	err = CopyFile(filepath.Join(toDir, "/etc/resolv.conf"), "/etc/resolv.conf")
	if err != nil {
		return errors.Wrap(err, "")
	}
	err = CopyFile(filepath.Join(toDir, "/etc/hosts"), "/etc/hosts")
	if err != nil {
		return errors.Wrap(err, "")
	}

	// STEP: prepare runc container spec
	envs := make([]string, 0, 10)
	for _, c := range manifestConfig {
		content, err := ioutil.ReadFile(filepath.Join(fromDir, c))
		if err != nil {
			return errors.Wrap(err, "")
		}
		var config ManifestConfig
		err = json.Unmarshal(content, &config)
		if err != nil {
			return errors.Wrap(err, "")
		}

		envs = append(envs, config.Config.Env...)
	}

	config, err := bmc.BinDataFs.ReadFile("config.json")
	if err != nil {
		return errors.Wrap(err, "")
	}
	var newConfig map[string]interface{}
	err = json.Unmarshal(config, &newConfig)
	if err != nil {
		return errors.Wrap(err, "unmarshal config.json template failed")
	}

	newConfig["process"].(map[string]interface{})["env"] = envs
	newConfig["process"].(map[string]interface{})["args"] = firstCmd
	newConfig["process"].(map[string]interface{})["terminal"] = terminal

	updatedConfig, err := json.Marshal(newConfig)
	if err != nil {
		return errors.Wrap(err, "")
	}

	configPath := filepath.Join(baseDir, "config.json")
	err = ioutil.WriteFile(configPath, updatedConfig, 0664)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flags := runCmd.Flags()
	flags.StringVarP(&runFromTar, "from", "f", "deploy.tar", "docker image tar")
	flags.StringVarP(&runToDir, "to", "t", "./tmp", "to base dir")
	flags.BoolVarP(&runRootfs, "build-rootfs", "", false, "build rootfs")
	flags.StringVarP(&runContainerID, "container-id", "", "bmc-container-id", "container id")
	flags.StringVarP(&runConfig, "config", "", "config.json", "runc container spec config")
	flags.StringVarP(&runRawCmd, "cmd", "", "sh", "container first cmd")
	flags.BoolVarP(&runTerminal, "terminal", "", false, "open terminal")
}

func CopyFile(dest, src string) error {
	bytesRead, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = ioutil.WriteFile(dest, bytesRead, 0644)
	return errors.Wrap(err, "")
}
