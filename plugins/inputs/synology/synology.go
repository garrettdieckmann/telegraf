package synology

import (
	"fmt"
	"strconv"

	sapi "github.com/garrettdieckmann/synologyapi"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

// Connection contains fields parsed from TOML configuration, for connection to NAS
type Connection struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Account  string `toml:"account"`
	Password string `toml:"password"`
}

// Description returns Synology plugin description
func (s *Connection) Description() string {
	return "Synology NAS plugin, using Synology API"
}

// SampleConfig returns sampleconfig for Input Plugin
func (s *Connection) SampleConfig() string {
	return `
## Synology NAS connection information
host = "192.168.1.1"
port = "5000"
account = "operator"
password = "password1"
`
}

// Gather returns metrics collection from Synology NAS API
func (s *Connection) Gather(acc telegraf.Accumulator) error {
	tags := map[string]string{}
	tags["host"] = s.Host

	synas, err := sapi.NewConnection(s.Host, s.Port, s.Account, s.Password)
	if err != nil {
		return fmt.Errorf("Error in connection: %v", err)
	}

	sysresp, err := synas.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("Error in GetSystemInfo: %v", err)
	}

	acc.AddFields("SynologyMeasurement", map[string]interface{}{"OneMinLoad": strconv.Itoa(sysresp.CPU.OneMinLoad)}, tags)

	return nil
}

func init() {
	inputs.Add("synology", func() telegraf.Input { return &Connection{} })
}
