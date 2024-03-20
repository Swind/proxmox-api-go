package qemu

import (
	"fmt"

	"github.com/Telmate/proxmox-api-go/cli"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/spf13/cobra"
)

var qemu_cloneCmd = &cobra.Command{
	Use:   "clone GUESTID",
	Short: "Clones the specified guest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		vmid := cli.ValidateIntIDset(args, "GuestID")
		node, _ := cmd.Flags().GetString("node")
		newid, _ := cmd.Flags().GetInt("newid")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		sourceVmr := proxmox.NewVmRef(vmid)
		sourceVmr.SetNode(node)

		c := cli.NewClient()
		params := map[string]interface{}{
			"newid":       newid,
			"name":        name,
			"description": description,
		}
		_, err = c.CloneQemuVm(sourceVmr, params)
		if err == nil {
			fmt.Println("true")
		} else {
			fmt.Println(err)
		}

		return
	},
}

func init() {
	qemuCmd.AddCommand(qemu_cloneCmd)

	// node name
	qemu_cloneCmd.Flags().StringP("node", "", "", "Node name")
	qemu_cloneCmd.Flags().IntP("newid", "", -1, "New VM ID")

	qemu_cloneCmd.Flags().StringP("name", "", "", "New VM Name")
	qemu_cloneCmd.Flags().StringP("description", "", "", "New VM Description")

	qemu_cloneCmd.MarkFlagRequired("node")
	qemu_cloneCmd.MarkFlagRequired("newid")
}
