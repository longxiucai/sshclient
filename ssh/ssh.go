package ssh

import (
	"fmt"
	"net"
	"time"

	klog "k8s.io/klog/v2"

	netUtils "github.com/longxiucai/sshclient/ssh/utils/net"
)

type Interface interface {
	// CmdAsync exec command on remote host, and asynchronous return logs
	CmdAsync(cmd ...string) error
	// Cmd exec command on remote host, and return combined standard output and standard error
	Cmd(cmd string) ([]byte, error)
	// CmdToString exec command on remote host, and return spilt standard output and standard error
	CmdToString(cmd, spilt string) (string, error)

	Ping() error

	GetIP() net.IP
}

type SSHClient struct {
	IsStdout     bool
	Encrypted    bool
	User         string
	Password     string
	Port         string
	PkFile       string
	PkPassword   string
	Timeout      *time.Duration
	LocalAddress []net.Addr
	// Fs           fs.Interface
	IP net.IP
}

func newSSHClient(ssh *SshConfig, isStdout bool, ip net.IP) Interface {
	if ssh.User == "" {
		ssh.User = "root"
	}
	address, err := netUtils.GetLocalHostAddresses()
	if err != nil {
		klog.Warningf("failed to get local address: %v", err)
	}
	return &SSHClient{
		IsStdout:     isStdout,
		Encrypted:    ssh.Encrypted,
		User:         ssh.User,
		Password:     ssh.Passwd,
		Port:         ssh.Port,
		PkFile:       ssh.Pk,
		PkPassword:   ssh.PkPasswd,
		LocalAddress: address,
		// Fs:           fs.NewFilesystem(),
		IP: ip,
	}
}

// GetHostSSHClient is used to executed bash command and no std out to be printed.
func GetHostSSHClient(hostIP net.IP, sshConfig *SshHosts) (Interface, error) {
	for _, host := range sshConfig.Hosts {
		for _, ip := range host.IPS {
			if hostIP.Equal(net.ParseIP(ip)) {
				// if err := mergo.Merge(&host.SSH, &sshConfig.SSH); err != nil {
				// 	return nil, err
				// }
				return newSSHClient(&host.SshConfig, false, hostIP), nil
			}
		}
	}
	return nil, fmt.Errorf("failed to get host ssh client: host ip %s not in hosts ip list", hostIP)
}

// GetStdoutSSHClient is used to show std out when execute bash command.
func GetStdoutSSHClient(hostIP net.IP, sshConfig *SshHosts) (Interface, error) {
	for _, host := range sshConfig.Hosts {
		for _, ip := range host.IPS {
			if hostIP.Equal(net.ParseIP(ip)) {
				// if err := mergo.Merge(&host.SSH, &sshConfig.SSH); err != nil {
				// 	return nil, err
				// }
				return newSSHClient(&host.SshConfig, true, hostIP), nil
			}
		}
	}
	return nil, fmt.Errorf("failed to get host ssh client: host ip %s not in hosts ip list", hostIP)
}
