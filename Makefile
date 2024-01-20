# 编译适用于 Linux 的可执行文件
linux:
	GOOS=linux GOARCH=amd64 go build -o output/bin/perp

# 编译适用于 macOS 的可执行文件
mac:
	GOOS=darwin GOARCH=amd64 go build -o perpetua main.go

# 编译适用于 Windows 的可执行文件
windows:
	GOOS=windows GOARCH=amd64 go build -o perpetua.exe main.go

# 默认目标为编译适用于 macOS 的可执行文件
default: linux
