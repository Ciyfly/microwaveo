package file2shellcode

import (
	"bytes"

	"github.com/Binject/go-donut/donut"
)

var donutArch donut.DonutArch

func Build(filePath string, arch string, moduleName string, format string) (*bytes.Buffer, error) {
	switch arch {
	case "x32", "386":
		donutArch = donut.X32
	case "x64", "amd64":
		donutArch = donut.X64
	case "x84":
		donutArch = donut.X84
	}
	config := new(donut.DonutConfig)
	config.Arch = donutArch
	config.Entropy = uint32(3)
	config.OEP = uint64(0)
	config.InstType = donut.DONUT_INSTANCE_PIC
	config.Parameters = ""
	config.Runtime = ""
	config.URL = ""
	config.Class = ""
	config.Method = ""
	config.Domain = ""
	config.Bypass = 3
	config.Method = moduleName
	config.Compress = uint32(1)
	config.Verbose = false
	config.ExitOpt = uint32(1)
	if format == "hex" {
		config.Format = uint32(8)
	} else {
		config.Format = uint32(1)
	}
	// run
	return donut.ShellcodeFromFile(filePath, config)

}
