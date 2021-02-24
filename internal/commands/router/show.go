package router

import (
	"encoding/json"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/network"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"sync"
)

// ShowCommand creates the "router show" command
func ShowCommand(service service.Network) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show current router"),
		service:     service,
	}
}

type showCommand struct {
	*commands.BaseCommand
	service service.Network
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getRouterArgCompletionFunction(s.service))
}

type routerWithNetworks struct {
	router   *upcloud.Router
	networks []*upcloud.Network
}

// MarshalJSON implements json.Marshaler
func (c *routerWithNetworks) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.router)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("one router uuid or name is required")
		}
		r, err := searchRouter(args[0], s.service, true)
		if err != nil {
			return nil, err
		}

		var networks []*upcloud.Network
		var wg sync.WaitGroup
		var getNetworkError error

		for _, n := range r[0].AttachedNetworks {
			wg.Add(1)
			go func(rn upcloud.RouterNetwork) {
				defer wg.Done()
				nw, err := network.SearchUniqueNetwork(rn.NetworkUUID, s.service)
				if err != nil {
					getNetworkError = err
				}
				networks = append(networks, nw)
			}(n)
		}
		wg.Wait()
		if getNetworkError != nil {
			return nil, getNetworkError
		}
		return &routerWithNetworks{
			router:   r[0],
			networks: networks,
		}, nil
	}
}

// HandleOutput implements Command.HandleOutput
func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	routerWithNetworks := out.(*routerWithNetworks)
	r := routerWithNetworks.router
	networks := routerWithNetworks.networks

	l := ui.NewListLayout(ui.ListLayoutDefault)

	dCommon := ui.NewDetailsView()
	dCommon.Append(
		table.Row{"UUID:", ui.DefaultUUUIDColours.Sprint(r.UUID)},
		table.Row{"Name:", r.Name},
		table.Row{"Type:", r.Type},
	)
	l.AppendSection("Common", dCommon.Render())

	if len(networks) > 0 {
		tIPRouter := ui.NewDataTable("UUID", "Name", "Router", "Type", "Zone")
		for _, n := range networks {
			tIPRouter.Append(table.Row{
				ui.DefaultUUUIDColours.Sprint(n.UUID),
				n.Name,
				ui.DefaultUUUIDColours.Sprint(n.Router),
				n.Type,
				n.Zone,
			})
		}
		l.AppendSection("Networks:", tIPRouter.Render())
	} else {
		l.AppendSection("Networks:", "no network found for this router")
	}
	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}