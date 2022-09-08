package conf

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"microwaveo/pkg/utils"
	"os"
	"os/exec"
	"path"
)

//go:embed template/template.tmpl
var t []byte

//go:embed template/template_aes.tmpl
var tAes []byte

//go:embed template/template_aes_white.tmpl
var tAesWhite []byte

func initDir(p string) {
	if !utils.FileIsExist(p) {
		os.Mkdir(p, 0666)
	}
}

func initTemplate(tName string, content []byte) {
	tmplPath := path.Join(utils.GetCurrentDirectory(), "conf", "template", tName)
	if !utils.FileIsExist(tmplPath) {
		ioutil.WriteFile(tmplPath, content, 0666)
	}
}

func Init() {
	confParh := path.Join(utils.GetCurrentDirectory(), "conf")
	initDir(confParh)
	templateParh := path.Join(confParh, "template")
	initDir(templateParh)
	staticPath := path.Join(templateParh, "static")
	initDir(staticPath)
	// template
	initTemplate("template.tmpl", t)
	initTemplate("template_aes.tmpl", tAes)
	initTemplate("template_aes_white.tmpl", tAesWhite)
}

func EnvironmentalTestGo() {
	// 需要go 环境 还有garble
	goCmd := exec.Command("go", "version")
	goErr := goCmd.Run()
	if goErr != nil {
		fmt.Println("you need to install go for build exe: https://studygolang.com/dl")
		os.Exit(-1)
	}
}

func EnvironmentalTestGarble() {
	garbleCmd := exec.Command("garble", "version")
	garbleErr := garbleCmd.Run()
	if garbleErr != nil {
		fmt.Println("You need to install garble for compilation: go install mvdan.cc/garble@latest")
		os.Exit(-1)
	}
}
