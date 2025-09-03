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

// 创建expect脚本进行SSH密码认证
func CreateExpectScript(host models.Host) (string, error) {
	scriptContent := fmt.Sprintf(`#!/usr/bin/expect -f
set timeout 30
spawn ssh -p %d %s@%s
expect {
    "yes/no" { send "yes\r"; exp_continue }
    "password:" { send "%s\r" }
}
interact
`, host.Port, host.Username, host.IP, host.Password)

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
				fmt.Printf("按任意键返回主菜单...\n")
				fmt.Scanln()
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
			fmt.Printf("按任意键返回主菜单...\n")

			fmt.Scanln()
			return
		}
	}

	// 添加端口参数
	if host.Port != 22 {
		sshArgs = append(sshArgs, "-p", strconv.Itoa(host.Port))
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
	fmt.Printf("按任意键返回主菜单...\n")

	fmt.Scanln()
}