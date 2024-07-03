package ssh

import (
	"os"
	"path/filepath"

	"k8s.io/klog/v2"
)

type SshConfig struct {
	Encrypted bool   `mapstructure:"encrypted" yaml:"encrypted,omitempty" json:"encrypted,omitempty"`
	User      string `mapstructure:"user" yaml:"user,omitempty" json:"user,omitempty"`
	Passwd    string `mapstructure:"passwd" yaml:"passwd,omitempty" json:"passwd,omitempty"`
	Pk        string `mapstructure:"pk" yaml:"pk,omitempty" json:"pk,omitempty"`
	PkPasswd  string `mapstructure:"pkPasswd" yaml:"pkPasswd,omitempty" json:"pkPasswd,omitempty"`
	Port      string `mapstructure:"port" yaml:"port,omitempty" json:"port,omitempty"`
}

type Host struct {
	IPS       []string `mapstructure:"ips" yaml:"ips,omitempty" json:"ips,omitempty"`
	SshConfig `mapstructure:"ssh" yaml:"ssh,omitempty" json:"ssh,omitempty"`
}

type SshHosts struct {
	Hosts []Host `mapstructure:"hosts" yaml:"hosts" json:"hosts"`
}

// Option is a function that configures an SSH instance.
type Option func(*SshConfig)

// NewSSH initializes and returns an SSH instance with required and optional parameters.
func NewSSHConfig(options ...Option) SshConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		klog.Errorf("Failed to get home directory: %v", err)
	}
	keyPath := filepath.Join(homeDir, ".ssh", "id_rsa")
	ssh := SshConfig{
		User: "root", // Default user
		Port: "22",   // Default port
		Pk:   keyPath,
	}

	for _, option := range options {
		option(&ssh)
	}
	klog.Infof("NewSSHConfig: %+v", ssh)
	return ssh
}

// Option functions for SSH,config user
func WithUser(user string) Option {
	return func(ssh *SshConfig) {
		ssh.User = user
	}
}

// Option functions for SSH,config user's password
func WithPasswd(passwd string) Option {
	return func(ssh *SshConfig) {
		ssh.Passwd = passwd
	}
}

// Option functions for SSH,config encrypted
func WithEncrypted(encrypted bool) Option {
	return func(ssh *SshConfig) {
		ssh.Encrypted = encrypted
	}
}

// Option functions for SSH,config private key
func WithPrivateKey(pk string) Option {
	return func(ssh *SshConfig) {
		ssh.Pk = pk
	}
}

// Option functions for SSH,config password of the private key
func WithPrivateKeyPassword(pkPasswd string) Option {
	return func(ssh *SshConfig) {
		ssh.PkPasswd = pkPasswd
	}
}

// Option functions for SSH,config port
func WithPort(port string) Option {
	return func(ssh *SshConfig) {
		ssh.Port = port
	}
}

// NewHost initializes and returns a Host instance
func NewHost(ips []string, ssh SshConfig) Host {
	return Host{
		IPS:       ips,
		SshConfig: ssh,
	}
}

// NewSSHConfig initializes and returns an SSHConfig instance
func NewSshHosts(hosts ...Host) SshHosts {
	return SshHosts{Hosts: hosts}
}
