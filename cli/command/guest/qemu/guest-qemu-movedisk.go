package qemu

import (
	"fmt"

	"github.com/Telmate/proxmox-api-go/cli"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/spf13/cobra"
)

var qemuMoveDiskCmd = &cobra.Command{
	Use:   "movedisk GUESTID",
	Short: "movedisks the specified guest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		vmid := cli.ValidateIntIDset(args, "GuestID")
		disk, _ := cmd.Flags().GetString("disk")
		targetStorage, _ := cmd.Flags().GetString("storage")

		sourceVmr := proxmox.NewVmRef(vmid)

		c := cli.NewClient()
		err = c.CheckVmRef(sourceVmr)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = c.MoveQemuDisk(sourceVmr, disk, targetStorage)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("true")
		}

		return
	},
}

func init() {
	qemuCmd.AddCommand(qemuMoveDiskCmd)

	// node name
	qemuMoveDiskCmd.Flags().StringP("disk", "", "", "Disk name")
	qemuMoveDiskCmd.Flags().StringP("storage", "", "", "Target Storage Name")

	qemuMoveDiskCmd.MarkFlagRequired("size")
	qemuMoveDiskCmd.MarkFlagRequired("storage")
}
