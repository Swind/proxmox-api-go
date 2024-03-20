package qemu

import (
	"fmt"

	"github.com/Telmate/proxmox-api-go/cli"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/spf13/cobra"
)

var qemu_resizeCmd = &cobra.Command{
	Use:   "resize GUESTID",
	Short: "resizes the specified guest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		vmid := cli.ValidateIntIDset(args, "GuestID")
		disk, _ := cmd.Flags().GetString("disk")
		size, _ := cmd.Flags().GetString("size")

		sourceVmr := proxmox.NewVmRef(vmid)

		c := cli.NewClient()
		err = c.CheckVmRef(sourceVmr)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = c.ResizeQemuDiskRaw(sourceVmr, disk, size)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("true")
		}

		return
	},
}

func init() {
	qemuCmd.AddCommand(qemu_resizeCmd)

	// node name
	qemu_resizeCmd.Flags().StringP("disk", "", "", "Disk name")
	qemu_resizeCmd.Flags().StringP("size", "", "", "New VM Size")

	qemu_resizeCmd.MarkFlagRequired("size")
	qemu_resizeCmd.MarkFlagRequired("disk")
}
