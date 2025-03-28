package monitor

/*
	Sliver Implant Framework
	Copyright (C) 2021  Bishop Fox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/gsmith257-cyber/better-sliver-package/client/console"
	"github.com/gsmith257-cyber/better-sliver-package/protobuf/commonpb"
)

// MonitorStartCmd - Start monitoring threat intel for implants
func MonitorStartCmd(cmd *cobra.Command, con *console.SliverClient, args []string) {
	resp, err := con.Rpc.MonitorStart(context.Background(), &commonpb.Empty{})
	if err != nil {
		con.PrintErrorf("%s\n", err)
		return
	}
	if resp != nil && resp.Err != "" {
		con.PrintErrorf("%s\n", resp.Err)
		return
	}
	con.PrintInfof("Started monitoring threat intel platforms for implants hashes\n")
}
