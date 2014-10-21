package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"strconv"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/codegangsta/cli"
)

type CliRpcService struct {
	listener net.Listener
	stopCh   chan struct{}
	Pinged   bool
	RpcCmd   *CliRpcCmd
}

type CliRpcCmd struct {
	ReturnData        interface{}
	coreCommandRunner *cli.App
}

func NewRpcService(commandRunner *cli.App) (*CliRpcService, error) {
	rpcService := &CliRpcService{
		RpcCmd: &CliRpcCmd{
			ReturnData:        new(interface{}),
			coreCommandRunner: commandRunner,
		},
	}

	err := rpc.Register(rpcService.RpcCmd)
	if err != nil {
		return nil, err
	}

	return rpcService, nil
}

func (cli *CliRpcService) Stop() {
	close(cli.stopCh)
	cli.listener.Close()
}

func (cli *CliRpcService) Port() string {
	return strconv.Itoa(cli.listener.Addr().(*net.TCPAddr).Port)
}

func (cli *CliRpcService) Start() error {
	var err error

	cli.stopCh = make(chan struct{})

	cli.listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := cli.listener.Accept()
			if err != nil {
				select {
				case <-cli.stopCh:
					return
				default:
					fmt.Println(err)
				}
			} else {
				go rpc.ServeConn(conn)
			}
		}
	}()

	return nil
}

func (cmd *CliRpcCmd) CallCoreCommand(args []string, retVal *string) error {
	err := cmd.coreCommandRunner.Run(args)
	return err
}

func (cmd *CliRpcCmd) SetPluginMetadata(pluginMetadata plugin.PluginMetadata, retVal *bool) error {
	cmd.ReturnData = interface{}(pluginMetadata)
	*retVal = true
	return nil
}
