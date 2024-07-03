# 创建ssh的Client
## 使用user+passwd
使用示例：
```
// 配置用户名与密码
sshconfig := ssh.NewSSHConfig(ssh.WithUser("root"), ssh.WithPasswd("MyPassword"))

// 使用 NewHost 初始化 Host，配置ip与密码对应关系
host1 := ssh.NewHost([]string{"10.42.186.32", "10.42.186.33"}, sshconfig)

// 使用 NewSshHosts 初始化 SshHosts
sshHosts := ssh.NewSshHosts(host1)

// 获取指定ip的sshclient
ip0 := net.ParseIP("10.42.186.33")
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
```
示例中通过`ssh.WithUser("root"), ssh.WithPasswd("MyPassword")`配置了用户名与密码
## 使用privateKey
1. 使用默认`~/.ssh/id_rsa`，且没有密码，端口`22`，用户为`root`
```
sshconfig := ssh.NewSSHConfig()
```
2. 需要配置私钥文件路径，且有密码
```
sshconfig := ssh.NewSSHConfig(ssh.WithPrivateKey("path/to/id_rsa"),ssh.WithPrivateKeyPassword("123456"))
```