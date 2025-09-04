package ssh

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/daihao4371/hostmanager/internal/models"
)

// 检查主机连通性
func CheckHostStatus(host models.Host) string {
	address := net.JoinHostPort(host.IP, strconv.Itoa(host.Port))
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return "offline"
	}
	conn.Close()
	return "online"
}

// 检查expect工具是否可用
func CheckExpectAvailable() bool {
	_, err := exec.LookPath("expect")
	return err == nil
}

// 检查 lrzsz 工具是否可用（sz/rz 命令）
func CheckLrzszAvailable() bool {
	_, szErr := exec.LookPath("sz")
	_, rzErr := exec.LookPath("rz")
	return szErr == nil && rzErr == nil
}

// 检查 Zmodem 支持的综合状态
func CheckZmodemSupport() (bool, string) {
	if !CheckLrzszAvailable() {
		return false, "系统缺少 lrzsz 工具包，请安装：brew install lrzsz (macOS) 或 apt install lrzsz (Ubuntu)"
	}
	return true, ""
}

// 创建expect脚本进行SSH密码认证（支持Zmodem）
func CreateExpectScript(host models.Host) (string, error) {
	// 构建SSH参数，支持Zmodem时添加必要选项
	sshArgs := fmt.Sprintf("-p %d", host.Port)
	if host.IsZmodemEnabled() {
		// 启用 Zmodem 支持需要的 SSH 选项
		sshArgs += " -o RequestTTY=yes -o RemoteCommand=\"exec \\$SHELL -l\""
	}

	scriptContent := fmt.Sprintf(`#!/usr/bin/expect -f
set timeout 30
spawn ssh %s %s@%s
expect {
    "yes/no" { send "yes\r"; exp_continue }
    "password:" { send "%s\r" }
}
interact
`, sshArgs, host.Username, host.IP, host.Password)

	tmpFile, err := os.CreateTemp("", "ssh_expect_*.exp")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(scriptContent)
	if err != nil {
		return "", err
	}

	err = os.Chmod(tmpFile.Name(), 0755)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// SSH连接函数
func Connect(host models.Host, onConnect func(models.Host)) {
	// 添加到连接历史
	if onConnect != nil {
		onConnect(host)
	}

	var cmd *exec.Cmd

	// 构建SSH连接命令
	sshArgs := []string{}

	// 处理认证方式
	if host.AuthType == "key" && host.KeyPath != "" {
		sshArgs = append(sshArgs, "-i", host.KeyPath)
	} else if host.AuthType == "password" || (host.AuthType == "key" && host.KeyPath == "" && host.Password != "") {
		if !CheckExpectAvailable() {
			fmt.Printf("错误: 系统缺少 expect 工具来支持密码认证\n")
			fmt.Printf("请手动输入密码进行连接\n")
		} else {
			scriptPath, err := CreateExpectScript(host)
			if err != nil {
				fmt.Printf("创建expect脚本失败: %v\n", err)
				// 不在这里等待输入，让UI层处理
				return
			}

			defer func() {
				os.Remove(scriptPath)
			}()

			cmd = exec.Command("expect", scriptPath)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Printf("\n🔗 正在连接到 %s (%s@%s:%d)...\n", host.Name, host.Username, host.IP, host.Port)
			fmt.Printf("💡 提示: 连接断开后将自动返回主菜单\n")
			fmt.Printf("═══════════════════════════════════════════════════════════\n")
			err = cmd.Run()
			if err != nil {
				fmt.Printf("连接失败: %v\n", err)
			}
			fmt.Printf("\n📋 与 %s 的连接已断开\n", host.Name)
			// 不在这里等待输入，让UI层处理
			return
		}
	}

	// 添加端口参数
	if host.Port != 22 {
		sshArgs = append(sshArgs, "-p", strconv.Itoa(host.Port))
	}

	// 添加 Zmodem 支持参数
	if host.IsZmodemEnabled() {
		// 检查 Zmodem 支持状态
		if supported, msg := CheckZmodemSupport(); !supported {
			fmt.Printf("⚠️  Zmodem 不可用: %s\n", msg)
		} else {
			fmt.Printf("📁 Zmodem 文件传输已启用 (sz/rz 命令可用)\n")
			// 添加必要的SSH选项以支持Zmodem
			sshArgs = append(sshArgs, "-o", "RequestTTY=yes")
		}
	}

	// 添加目标地址
	sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", host.Username, host.IP))

	// 构建SSH命令
	cmd = exec.Command("ssh", sshArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("\n🔗 正在连接到 %s (%s@%s:%d)...\n", host.Name, host.Username, host.IP, host.Port)
	fmt.Printf("💡 提示: 连接断开后将自动返回主菜单\n")
	fmt.Printf("═══════════════════════════════════════════════════════════\n")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
	}
	fmt.Printf("\n📋 与 %s 的连接已断开\n", host.Name)
	// 不在这里等待输入，让UI层统一处理
}