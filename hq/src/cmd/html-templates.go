package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const templatesPath = "./templates"

func main() {
	filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
		fileExtName := strings.ToLower(filepath.Ext(path))
		if info.IsDir() || fileExtName != ".html" {
			return nil
		}
		bts, err := ioutil.ReadFile(path)
		if err == nil {
			s := `
			package templates
			func init(){
			Data["` + info.Name() + `"]=[]byte(` + "`" + string(bts) + "`" + `)
			}`
			d, _ := format.Source([]byte(s))
			ioutil.WriteFile(path+".go", d, 0666)
		} else {
			fmt.Println(string(bts), err)
		}
		return nil
	})
}
