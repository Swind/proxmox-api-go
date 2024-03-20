package qemu

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Telmate/proxmox-api-go/cli"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/spf13/cobra"
)

var qemu_configCmd = &cobra.Command{
	Use:   "config GUESTID",
	Short: "configs the specified guest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		vmid := cli.ValidateIntIDset(args, "GuestID")
		config, _ := cmd.Flags().GetString("config")

		var configReader io.Reader
		if config == "" {
			// Read from stdin
			configReader = cmd.InOrStdin()
		} else {
			// Read from file
			configReader, err = os.Open(config)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		configBytes, err := io.ReadAll(configReader)
		if err != nil {
			fmt.Println(err)
			return
		}

		configQemu := proxmox.ConfigQemu{}
		json.Unmarshal(configBytes, &configQemu)
		sourceVmr := proxmox.NewVmRef(vmid)

		c := cli.NewClient()
		err = c.CheckVmRef(sourceVmr)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = configQemu.Update(true, sourceVmr, c)
		if err != nil {
			fmt.Println(err)
		}

		return
	},
}

func init() {
	qemuCmd.AddCommand(qemu_configCmd)

	// node name
	qemu_configCmd.Flags().StringP("config", "", "", "Config file path")
}
