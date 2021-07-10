package server

import (
	"github.com/spf13/cobra"
)

// AddServerListCmd Command for listing the members in the cluster
func AddServerListCmd(parentCmd *cobra.Command) {
	var serverAddress string
	var port int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all the members in cluster",

		//RunE: func(cmd *cobra.Command, args []string) error {
		//	resp, err := server.SendMessage(&server.MemberServer{
		//		Port:    port,
		//		Address: serverAddress,
		//	}, &server.Request{
		//		ID:        server.ListMembers,
		//		Arguments: nil,
		//	}, nil)
		//	if err != nil {
		//		return err
		//	}
		//	marshal, err := json.MarshalIndent(resp.Resp, "", "   ")
		//	if err != nil {
		//		return err
		//	}
		//	fmt.Println(string(marshal))
		//	return nil
		//},
	}
	cmd.Flags().StringVar(&serverAddress, "server", "127.0.0.1", "IP address or host name")
	_ = cmd.MarkFlagRequired("server")

	cmd.Flags().IntVar(&port, "port", 9191, "Port to use for communication")
	_ = cmd.MarkFlagRequired("port")

	parentCmd.AddCommand(cmd)
}
