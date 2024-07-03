package ssh

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
func NewSSHConfig(user, passwd string, options ...Option) SshConfig {
	ssh := SshConfig{
		User:   user,
		Passwd: passwd,
		Port:   "22", // Default port
	}

	for _, option := range options {
		option(&ssh)
	}

	return ssh
}

// Option functions for SSH
func WithEncrypted(encrypted bool) Option {
	return func(ssh *SshConfig) {
		ssh.Encrypted = encrypted
	}
}

func WithPrivateKey(pk string) Option {
	return func(ssh *SshConfig) {
		ssh.Pk = pk
	}
}

func WithPrivateKeyPassword(pkPasswd string) Option {
	return func(ssh *SshConfig) {
		ssh.PkPasswd = pkPasswd
	}
}

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
