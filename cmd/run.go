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
			err := os.RemoveAll(toDir)
			if err != nil {
				log.Printf("remove dir %v failed %v", toDir, err)
			}
			err = tar.Unarchive(
				runFromTar,
				fromDir,
			)
			if err != nil {
				return errors.Wrap(err, "")
			}

			err = prepareRootfs(ctx, fromDir, toDir)
			if err != nil {
				return errors.Wrap(err, "")
			}
		}

		configPath := filepath.Join(baseDir, "config.json")
		// STEP: prepare runc container spec
		if !fileExists(configPath) {
			config, err := bmc.BinDataFs.ReadFile("config.json")
			if err != nil {
				return errors.Wrap(err, "")
			}
			err = ioutil.WriteFile(configPath, config, 0664)
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

func prepareRootfs(ctx context.Context, fromDir, toDir string) error {
	manifestFile := "manifest.json"

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
}

func CopyFile(dest, src string) error {
	bytesRead, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = ioutil.WriteFile(dest, bytesRead, 0644)
	return errors.Wrap(err, "")
}
