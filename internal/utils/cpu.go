package utils

import (
	"errors"
	"fmt"
	global "github.com/IUnlimit/perpetua/internal"
	"github.com/IUnlimit/perpetua/internal/logger"
	"github.com/bytedance/gopkg/util/gopool"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func RunExec(mx *sync.Mutex) error {
	execName := "Lagrange.OneBot"
	if IsWinPlatform() {
		execName += ".exe"
	}
	cmdDir := global.ParentPath
	execPath := global.LgrFolder + execName
	cmd := exec.Command(execPath)

	cmd.Dir = cmdDir
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	// 将错误输出与标准输出连接至同一管道
	cmd.Stderr = cmd.Stdout
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	// 将命令行输入复制到 stdin 管道中
	gopool.Go(func() {
		_, err := io.Copy(in, os.Stdin)
		if err != nil {
			log.Fatalf("Failed to copy stdin: %v", err)
		}
	})

	if err = cmd.Start(); err != nil {
		return err
	}
	global.OneBotProcess = cmd.Process
	if mx != nil {
		mx.Unlock()
	}

	var n int
	hook := logger.Hook
	bytes := make([]byte, 8*1024)
	for {
		n, err = out.Read(bytes)
		if err != nil {
			break
		}
		err = hook.ExecLogWrite(string(bytes[:n]))
		if err != nil {
			log.Warnf("Write exec log error: %v", err)
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func IsWinPlatform() bool {
	return strings.Contains(strings.ToLower(runtime.GOOS), "win")
}

func RandomAvailablePort(rangePort bool, start int, end int) (int, error) {
	if !rangePort {
		port, err := tryListen(0)
		if err != nil {
			return 0, err
		}
		return port, nil
	}

	portRange := end - start
	if start < 512 || end > 65535 || portRange < 0 {
		return 0, errors.New(fmt.Sprintf("invalid port range [%d-%d]", start, end))
	}
	if portRange == 0 {
		return start, nil
	}

	for i := start; i <= end; i++ {
		port, err := tryListen(i)
		if err == nil {
			return port, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("unavailable port in range [%d-%d]", start, end))
}

func tryListen(port int) (int, error) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, err
	}
	_ = listen.Close()

	port = listen.Addr().(*net.TCPAddr).Port
	return port, nil
}
