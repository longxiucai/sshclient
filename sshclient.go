package main

import (
	"fmt"
	"net"

	"github.com/longxiucai/sshclient/ssh"
)

func main() {
	// 配置用户名与密码
	sshconfig := ssh.NewSSHConfig("root", "MyPassword")

	// 使用 NewHost 初始化 Host，配置ip与密码对应关系
	host1 := ssh.NewHost([]string{"10.42.186.32", "10.42.186.33"}, sshconfig)

	// 使用 NewSshHosts 初始化 SshHosts
	sshHosts := ssh.NewSshHosts(host1)

	// 获取指定ip的sshclient
	ip0 := net.ParseIP("10.42.186.33") // 将ip地址转换成net.IP
	sshClient0, err := ssh.GetHostSSHClient(ip0, &sshHosts)
	if err != nil {
		panic(err)
	}

	// 进行后续操作...
	if err := ssh.WaitSSHReady(sshClient0, 3); err != nil {
		panic(err)
	}
	res, err := sshClient0.CmdToString("ls /root/", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	ip1 := net.ParseIP("10.42.186.32")
	sshClient, err := ssh.GetStdoutSSHClient(ip1, &sshHosts)
	if err != nil {
		panic(err)
	}
	sshClient.CmdAsync("tail -f /var/log/messages")
}
