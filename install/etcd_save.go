package install

import (
	"fmt"
	"github.com/wonderivan/logger"
	"os"
	"strings"
)

// 建议 nodes和master 采用hosts进行封装一下子.
type Hosts struct {
	// Name of the host provisioned via docker machine
	NodeName string
	// IP:port that is fully resolvable and used for SSH communication
	Address string
	// Optional - Internal address that will be used for components communication
	InternalAddress string
	// SSH usesr
	User string
	//SSH password
	Password string
	// SSH Private KeyPassword
	SSHKeyPassword string
	// SSH Private Key Path
	SSHKeyPath string
	// Node Labels
	Labels map[string]string `yaml:"labels" json:"labels,omitempty"`
}

type EtcdBackFlags struct {
	Name      string
	Dir       string
	EtcdHosts []string
	SealConfig
}

func GetEtcdBackFlags() *EtcdBackFlags {
	e := &EtcdBackFlags{}
	err := e.Load("")
	if err != nil {
		logger.Error(err)
		e.ShowDefaultConfig()
		os.Exit(0)
	}
	// get Etcd host
	e.EtcdHosts = e.Masters
	e.Dir = EtcdBackDir
	e.Name = SnapshotName

	return e
}


// 只需要在master上备份一次即可， 然后复制snapshot到各etcd节点。
func SnapshotEtcd(e *EtcdBackFlags) {
	cmdMkdir := fmt.Sprintf("mkdir -p %s || true", e.Dir)
	CmdWorkSpace(e.Masters[0], cmdMkdir, TMPDIR)
	host := reFormatHostToIp(e.Masters[0])
	err := SnapshotEtcdDefaultSave(host, e.Name, e.Dir)
	if err != nil {
		logger.Error("etcd back error: ", err)
		os.Exit(-1)
	}
	if len(e.Masters) > 1 {
		path := fmt.Sprintf("%s/%s", e.Dir, e.Name)
		SendPackage(path, e.Masters[1:], e.Dir, nil, nil)
	}
	err = HealthCheck(reFormatHostToIp(e.Masters[0]))
	if err != nil {
		logger.Info("health check is failed")
	}
}

func SnapshotEtcdDefaultSave(host, snapshotName, dir string) error {
	// use default
	eCmd := getDefaultCmd()
	endpoints := fmt.Sprintf("%s:2379", host)
	cmd := fmt.Sprintf(`%s--endpoints %s snapshot save %s`, eCmd, endpoints, snapshotName)
	fmt.Println(cmd)
	if err := CmdWork(host, cmd, dir); err != nil {
		return err
	}
	return nil
}

func reFormatHostToIp(host string) string {
	if strings.Contains(host, ":") {
		s := strings.Split(host, ":")
		return s[0]
	}
	return host
}