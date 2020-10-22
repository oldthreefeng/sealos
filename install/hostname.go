package install

import "fmt"


var (
	masterPrefix = "k8s-sealos-master-00"
	nodePrefix = "k8s-sealos-node-00"
)

func (s *SealosInstaller) SetHostname() {
	for k, v := range s.Masters {
		prefix := fmt.Sprintf("%s-%d", masterPrefix, k)
		setHostname(prefix, v)
	}
	for k, v := range s.Nodes {
		prefix := fmt.Sprintf("%s-%d", nodePrefix, k)
		setHostname(prefix, v)
	}
}

func setHostname(hostPrefix string, host string) {
	cmd := fmt.Sprintf("hostnamectl set-hostname %s", hostPrefix)
	SSHConfig.CmdAsync(host, cmd)
}
