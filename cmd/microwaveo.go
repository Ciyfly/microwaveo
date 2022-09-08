package main

import (
	"fmt"
	"microwaveo/conf"
	"microwaveo/pkg/cupboard"
	"microwaveo/pkg/logger"
	"microwaveo/pkg/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt, os.Kill, syscall.SIGKILL)
	go func() {
		s := <-c
		fmt.Println(fmt.Sprintf("recv signal: %d", s))
		fmt.Println("ctrl+c exit")
		os.Exit(0)
	}()
}

func init() {
	logger.Init()
	conf.Init()
}
func main() {
	SetupCloseHandler()
	app := cli.NewApp()
	app.Name = "mcrowaveo"
	app.Usage = "mcrowaveo -i test.dll "
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "input file dll/exe/shellcode",
		},
		&cli.StringFlag{
			Name:    "arch",
			Aliases: []string{"a"},
			Value:   "x64",
			Usage:   "shellcode arch x32 x64 x84 default x64",
		},
		&cli.StringFlag{
			Name:    "funcname",
			Aliases: []string{"fn"},
			Usage:   "dll func name",
		},
		&cli.StringFlag{
			Name:    "shellcodeFormat",
			Aliases: []string{"s"},
			Value:   "bin",
			Usage:   "output shellcode format hex bin default bin",
		},
		&cli.BoolFlag{
			Name:    "obfuscate",
			Aliases: []string{"of"},
			Usage:   "obfuscate the generated exe using garble",
		},
		//encrypt
		&cli.StringFlag{
			Name:    "encrypt",
			Aliases: []string{"e"},
			Value:   "aes",
			Usage:   "encrypt the generated exe support aes default aes",
		},
		&cli.StringFlag{
			Name:    "white",
			Aliases: []string{"w"},
			Usage:   "bundled white files  file path",
		},
	}
	app.Action = RunMain

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalf("cli.RunApp err: %s", err.Error())
		return
	}
}

func RunMain(c *cli.Context) error {

	input := c.String("input")
	white := c.String("white")
	arch := c.String("arch")
	funcName := c.String("funcname")
	obfuscate := c.Bool("obfuscate")
	encrypt := c.String("encrypt")
	shellcodeFormat := c.String("shellcodeFormat")
	if input == "" {
		logger.Fatal("You need to enter the dll exe shellcode file path specified by -i")
		os.Exit(-1)
	}
	if utils.FileIsExist(input) == false {
		logger.Fatal("input file not exist")
		os.Exit(-1)
	}
	args := &cupboard.Cmdargs{
		Input:           input,
		Arch:            arch,
		FuncName:        funcName,
		ShellcodeFormat: shellcodeFormat,
		White:           white,
		Obfuscate:       obfuscate,
		Encrypt:         encrypt,
	}
	conf.EnvironmentalTestGo()
	if obfuscate {
		conf.EnvironmentalTestGarble()
	}
	cupboard.Build(args)
	return nil
}
