package id

import (
	"fmt"

	"github.com/Telmate/proxmox-api-go/cli"
	"github.com/spf13/cobra"
)

var id_nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Returns the lowest available ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		currentId := cli.ValidateIntIDset(args, "ID")
		c := cli.NewClient()
		id, err := c.GetNextID(currentId)
		if err != nil {
			return
		}
		fmt.Fprintf(idCmd.OutOrStdout(), "%d\n", id)
		return
	},
}

func init() {
	idCmd.AddCommand(id_nextCmd)
}
