VERSION ?= v0.3.0
LDFLAGS = -ldflags "-X 'github.com/IUnlimit/perpetua/internal/conf.Version=$(VERSION)' -s -w"
OUTPUT  = output

# 编译适用于 Linux 的可执行文件
linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT)/perp

# 编译适用于 macOS 的可执行文件
mac:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT)/perp

# 编译适用于 macOS ARM 的可执行文件
mac-arm:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT)/perp

# 编译适用于 Windows 的可执行文件
windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT)/perp.exe

# 全平台编译
all: clean
	@mkdir -p $(OUTPUT)
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT)/perp-linux-amd64
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT)/perp-darwin-amd64
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT)/perp-darwin-arm64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT)/perp-windows-amd64.exe
	@echo "Build complete: $(OUTPUT)/"

# Docker 镜像构建
docker:
	docker build -t perpetua:$(VERSION) .

# 清理构建产物
clean:
	rm -rf $(OUTPUT)

# 默认目标
default: linux
