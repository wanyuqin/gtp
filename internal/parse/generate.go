package parse

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/sys/unix"

	"gihub.com/wanyuqin/go-to-proto/internal/config"
	"gihub.com/wanyuqin/go-to-proto/internal/object"
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
	di, err := os.Stat(basicSet.OutputPath)
	if err != nil && !errors.Is(err, unix.ENOENT) {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if di != nil && !di.IsDir() {
		err = errors.New("output path must be a directory")
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(basicSet.OutputPath, 0777)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return

		}

		fn = filepath.Join(basicSet.OutputPath, fn)

		fi, err := os.Create(fn)
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
}
