# Microwaveo

一个小工具 微波炉加热一下dll  

1. 调用 go-donut 将dll/exe等转为shellcode
2. 使用go模板构建shellcode的加载器 最后输出exe
3. 也支持直接传入shellcode来构建最后的exe
4. 最后的exe支持 garble混淆

## 注意
**因为是使用go 构建exe所以需要go的环境**

## 使用

编译好的文件可以直接在 在这里 [releases](https://github.com/Ciyfly/microwaveo/releases)  下载 当然可以自己编译

```shell
./microwaveo --help
GLOBAL OPTIONS:
   --arch value, -a value             shellcode arch x32 x64 x84 default x64 (default: "x64")
   --encrypt value, -e value          encrypt the generated exe support aes default aes (default: "aes")
   --funcname value, --fn value       dll func name
   --help, -h                         show help (default: false)
   --input value, -i value            input file dll/exe/shellcode
   --obfuscate, --of                  obfuscate the generated exe using garble (default: false)
   --shellcodeFormat value, -s value  output shellcode format hex bin default bin (default: "bin")
   --version, -v                      print the version (default: false)
   --white value, -w value            bundled white files  file path
```

### 将dll转为exe

```shell
./microwaveo -i recar.dll -fn RunRecar
```

### 将exe控制为32位
```
./microwaveo -i recar.dll -fn RunRecar -a x32
```

### 使用garble混淆最后的exe
```
./microwaveo -i recar.dll -fn RunRecar --of
```
需要安装 garble
最简单的安装 使用 `go install mvdan.cc/garble@latest` 最后配置环境变量


## TODO

1. 思考是不是可以将加载器做成多个模板的形式来处理
2. 增加一些反沙箱的东西
3. 待定 有任何想法欢迎与我交流


