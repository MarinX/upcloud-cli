package router

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	service service.Network
	req     request.ModifyRouterRequest
}

// ModifyCommand creates the "router modify" command
func ModifyCommand(service service.Network) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a router"),
		service:     service,
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getRouterArgCompletionFunction(s.service))
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.req.Name, "name", "", "Name of the router. [Required]")
	s.AddFlags(fs)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.req.Name == "" {
			return nil, fmt.Errorf("name is required")
		}

		return routerRequest{
			BuildRequest: func(uuid string) interface{} {
				s.req.UUID = uuid
				return &s.req
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.ModifyRouterRequest).UUID },
				MaxActions:    maxRouterActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Modifying router",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyRouter(req.(*request.ModifyRouterRequest))
				},
			},
		}.send(args)
	}
}