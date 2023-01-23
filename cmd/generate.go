package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"gihub.com/wanyuqin/gtp/internal/config"
	"gihub.com/wanyuqin/gtp/internal/object"
	"gihub.com/wanyuqin/gtp/internal/parse"
	"gihub.com/wanyuqin/gtp/utils"
)

type GenerateOptions struct {
	File   string
	Dir    string
	Config string
}

var FilePathNotFound = errors.New("file path not found")

var goSuffix = ".go"

var generateExample = `
The generate command must specify a directory or a file.
When both the directory and the file are specified, 
the directory will be the main one, and the specified content of the file will be ignored. 
The configuration file is optional. If no configuration file is specified, 
the default configuration file will be used.

# specified file
ptg generate -f ./example/user.go

# specified directory
ptg generate -d ./example/

# specified config file
ptg generate -f ./example/user.go -c ./example/example-config.yamlk
`

func NewGenerateCmd() *cobra.Command {
	options := GenerateOptions{}
	generateCmd := &cobra.Command{
		Use:     "generate",
		Example: generateExample,
		Short:   "",
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(options.Check())

		},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(options.Run())
		},
	}

	generateCmd.Flags().StringVarP(&options.File, "file", "f", "", "specify a single go file")
	generateCmd.Flags().StringVarP(&options.Dir, "dir", "d", "", "specify go file directory")
	generateCmd.Flags().StringVarP(&options.Config, "config", "c", "", "specify generate config")

	return generateCmd
}

func (g *GenerateOptions) Check() error {
	if g.Dir == "" && g.File == "" {
		return FilePathNotFound
	}

	if g.Config == "" {
		fmt.Fprintln(os.Stdout, "no configuration file specified, use the default configuration file")
	}
	if g.Dir != "" {
		g.File = ""
	}

	if g.File != "" {
		fileInfo, err := os.Stat(g.File)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return errors.New(fmt.Sprintf("%s can`t be a folder", g.File))
		}

		if !strings.HasSuffix(fileInfo.Name(), goSuffix) {
			return errors.New(fmt.Sprintf("%s not a go file", fileInfo.Name()))
		}
	}

	if g.Dir != "" {
		fileInfo, err := os.Stat(g.Dir)
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			return errors.New(fmt.Sprintf("%s must be a folder", g.Dir))
		}

	}

	return nil
}

func (g *GenerateOptions) Run() error {

	// load	config
	basicSet, err := config.InitConfig(g.Config)
	if err != nil {
		return err
	}
	// get go file
	filePaths := make([]string, 0)
	if g.Dir != "" {
		filePaths = getGoFilePath(g.Dir)
	}
	if g.File != "" {
		filePaths = append(filePaths, g.File)
	}

	// generate
	for _, v := range filePaths {
		fi, err := os.Stat(v)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		index := strings.LastIndex(fi.Name(), goSuffix)
		if index < 0 {
			continue
		}

		basicSet.FileName = fi.Name()[:index]
		objectList, err := object.ParseGo(v)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		parse.Generate(basicSet, objectList)

	}

	return nil
}

func getGoFilePath(path string) []string {
	filePath := make([]string, 0)
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, goSuffix) && !info.IsDir() {
			filePath = append(filePath, path)
		}
		return nil
	})

	return filePath
}
