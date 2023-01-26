package parse

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"gihub.com/wanyuqin/gtp/internal/config"
	"gihub.com/wanyuqin/gtp/internal/object"
)

func Add(x, y int) int {
	return x + y
}

var templatePath = "./template/proto.tmpl"
var protoSuffix = "proto"

func Generate(basicSet config.BasicSet, objList object.ObjectList) {
	t := template.New("proto.tmpl")
	t.Funcs(template.FuncMap{
		"add": Add,
	})

	files, err := t.ParseFiles(templatePath)

	if err != nil {
		log.Fatal(err)
	}

	messages := objList.BuildParam()

	p := config.Proto{
		BasicSet: basicSet,
		Message:  messages,
	}
	fn := fmt.Sprintf("%s.%s", basicSet.FileName, protoSuffix)
	fn = filepath.Join(basicSet.OutputPath, fn)
	fi, err := os.Create(fn)
	defer fi.Close()
	_, err = os.Stat(fn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	err = files.Execute(fi, p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return

	}
}
