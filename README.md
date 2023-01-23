# go-to-proto

## 简介

go-to-proto is a tool to convert the structure in the go file into the grpc-gateway format proto


## command line 
```shell
ptg is a tool for generating proto file through go

Usage:
  ptg [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    
  help        Help about any command

Flags:
  -h, --help      help for ptg
  -v, --version   version for ptg

Use "ptg [command] --help" for more information about a command.
```


### generate command
```shell

Usage:
  ptg generate [flags]

Examples:

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


Flags:
  -c, --config string   specify generate config
  -d, --dir string      specify go file directory
  -f, --file string     specify a single go file
  -h, --help            help for generate

```
### config yaml
```shell
packageName: "default"
goPackageName: "/pb"
apiVersion: "v1"
outputPath: "./proto"
```

