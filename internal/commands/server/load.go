package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/storage"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type loadCommand struct {
	*commands.BaseCommand
	params loadParams
}

type loadParams struct {
	request.LoadCDROMRequest
}

// LoadCommand creates the "server load" command
func LoadCommand() commands.Command {
	return &loadCommand{
		BaseCommand: commands.New("load", "Load a CD-ROM into the server"),
	}
}

var defaultLoadParams = &loadParams{
	LoadCDROMRequest: request.LoadCDROMRequest{},
}

// InitCommand implements Command.InitCommand
func (s *loadCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.Config()))
	s.params = loadParams{LoadCDROMRequest: request.LoadCDROMRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.StorageUUID, "storage", defaultLoadParams.StorageUUID, "The UUID of the storage to be loaded in the CD-ROM device.")

	s.AddFlags(flagSet)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *loadCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.StorageUUID == "" {
			return nil, fmt.Errorf("storage is required")
		}

		serverSvc := s.Config().Service.Server()
		storageSvc := s.Config().Service.Storage()
		strg, err := storage.SearchSingleStorage(s.params.StorageUUID, storageSvc)
		if err != nil {
			return nil, err
		}
		s.params.StorageUUID = strg.UUID

		return Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.LoadCDROMRequest
				req.ServerUUID = uuid
				return &req
			},
			Service:    serverSvc,
			ExactlyOne: true,
			Handler: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.LoadCDROMRequest)
					return fmt.Sprintf("Loading %q as a CD-ROM of server %q", req.StorageUUID, req.ServerUUID)
				},
				MaxActions: maxServerActions,
				Action: func(req interface{}) (interface{}, error) {
					return storageSvc.LoadCDROM(req.(*request.LoadCDROMRequest))
				},
			},
		}.Send(args)
	}
}
