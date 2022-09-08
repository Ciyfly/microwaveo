package loader

import (
	"context"
	"fmt"
	"io/ioutil"
	"microwaveo/pkg/logger"
	"microwaveo/pkg/utils"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"time"
)

type TmplValue struct {
	FileName  string
	Obfuscate bool
	Encrypt   string
	AesKey    string
	White     string
	Arch      string
}

func Build(tv TmplValue) {
	var tmplPath string
	goCodeDir := path.Join(utils.GetCurrentDirectory(), "conf", "template")
	if tv.Encrypt == "aes" {
		tmplPath = path.Join(goCodeDir, "template_aes.tmpl")
	} else {
		tmplPath = path.Join(goCodeDir, "template.tmpl")
	}
	if tv.White != "" {
		tmplPath = path.Join(goCodeDir, "template_aes_white.tmpl")
	}
	logger.Printf("use loader tmpl: %s", tmplPath)
	goCodePath := path.Join(goCodeDir, "tmp.go")
	outputpath := path.Join(utils.GetCurrentDirectory(), tv.FileName+".exe")
	tpl, err := template.ParseFiles(tmplPath)

	if err != nil {
		logger.Fatalf("parse tmpl error: %s", err.Error())
	}
	if utils.FileIsExist(goCodePath) {
		os.Remove(goCodePath)
	}
	f, err := os.Create(goCodePath)
	if err != nil {
		logger.Fatalf("os create exe error: %s,%s", err.Error(), goCodePath)
	}
	tpl.Execute(f, tv)
	// f.Close()
	f.Close()
	time.Sleep(2 * time.Second)

	// go build tmpl
	var cmd *exec.Cmd
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
	cmd = exec.CommandContext(ctx, "go", "build", "-ldflags", "-s -w -H windowsgui", "-o", outputpath, goCodePath)

	// 如果开启了混淆使用 garble来build
	if tv.Obfuscate {
		logger.Print("use garble build exe, It will take a long time, please wait")
		cmd = exec.CommandContext(ctx, "garble", "build", "-ldflags", "-s -w -H windowsgui", "-o", outputpath, goCodePath)
	}
	closer, err := cmd.StdoutPipe()
	defer func() {
		cancelFunc()
		_ = closer.Close()
		_ = cmd.Wait()
	}()
	cmd.Dir = goCodeDir
	// arch env
	var arch = "amd64"
	if tv.Arch == "x32" {
		arch = "386"
	}
	logger.Printf("arch: %s", arch)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOARCH=%s", arch))
	// windows
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", "windows"))
	cmd.Env = append(cmd.Env, fmt.Sprintf("CGO_ENABLED=%s", "0"))

	err = cmd.Start()
	if err != nil {
		logger.Fatalf("go build to exe error: %s", err.Error())
	}
	bytes, err := ioutil.ReadAll(closer)
	if err != nil {
		logger.Fatal("go build to exe error")
	}
	if string(bytes) != "" {
		logger.Printf("build tmpl error %s", strings.TrimSpace(string(bytes)))
	}
	err = os.Remove(goCodePath)
	if utils.FileIsExist(outputpath) {
		logger.Printf("generated exe: %s", outputpath)
	} else {
		logger.Printf("Generated exe error Could be a template or environment issue")
	}
}
