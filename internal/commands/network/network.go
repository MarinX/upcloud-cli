package network

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

const maxNetworkActions = 10
const positionalArgHelp = "<UUID/Name...>"

// BaseNetworkCommand creates the base "network" command
func BaseNetworkCommand() commands.Command {
	return &networkCommand{commands.New("network", "Manage network")}
}

type networkCommand struct {
	*commands.BaseCommand
}

func getNetworkUUID(in interface{}) string {
	return in.(*upcloud.Network).UUID
}

// SearchUniqueNetwork returns exactly one network with name or uuid matching *term*
func SearchUniqueNetwork(term string, service service.Network) (*upcloud.Network, error) {
	result, err := SearchNetwork(term, service, true)
	if err != nil {
		return nil, err
	}
	if len(result) > 1 {
		return nil, fmt.Errorf("multiple networks matched to query %q, use UUID to specify", term)
	}
	return result[0], nil
}

var cachedNetworks []upcloud.Network

// SearchNetwork returns all networks whose name or uuid matches term.
// It will get the available networks from service once and cache the results on future calls
func SearchNetwork(term string, service service.Network, unique bool) ([]*upcloud.Network, error) {
	var result []*upcloud.Network

	if len(cachedNetworks) == 0 {
		networks, err := service.GetNetworks()
		if err != nil {
			return nil, err
		}
		cachedNetworks = networks.Networks
	}

	for _, n := range cachedNetworks {
		network := n
		if network.UUID == term || network.Name == term {
			result = append(result, &network)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no network was found with %s", term)
	}
	if len(result) > 1 && unique {
		return nil, fmt.Errorf("multiple networks matched to query %q, use UUID to specify", term)
	}

	return result, nil
}

func searchAllNetworks(terms []string, service service.Network, unique bool) ([]string, error) {
	return commands.SearchResources(
		terms,
		func(id string) (interface{}, error) {
			return SearchNetwork(id, service, unique)
		},
		func(in interface{}) string { return in.(*upcloud.Network).UUID })
}

type networkRequest struct {
	ExactlyOne    bool
	BuildRequest  func(uuid string) interface{}
	Service       service.Network
	HandleContext ui.HandleContext
}

func (s networkRequest) send(args []string) (interface{}, error) {
	if s.ExactlyOne && len(args) != 1 {
		return nil, fmt.Errorf("single network uuid or name is required")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("at least one network uuid or name is required")
	}

	servers, err := searchAllNetworks(args, s.Service, true)
	if err != nil {
		return nil, err
	}

	var requests []interface{}
	for _, server := range servers {
		requests = append(requests, s.BuildRequest(server))
	}

	return s.HandleContext.Handle(requests)
}

func getArgCompFn(s service.Network) func(toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(toComplete string) ([]string, cobra.ShellCompDirective) {
		networks, err := s.GetNetworks()
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range networks.Networks {
			vals = append(vals, v.UUID, v.Name)
		}
		return commands.MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
	}
}

func handleNetwork(in string) (*upcloud.IPNetwork, error) {
	result := &upcloud.IPNetwork{}
	var dhcp string
	var dhcpDefRout string
	var dns string

	args, err := commands.Parse(in)
	if err != nil {
		return nil, err
	}

	fs := &pflag.FlagSet{}
	fs.StringVar(&dns, "dhcp-dns", dns, "Defines if the gateway should be given as default route by DHCP. Defaults to yes on public networks, and no on other ones.")
	fs.StringVar(&result.Address, "address", result.Address, "Sets address space for the network.")
	fs.StringVar(&result.Family, "family", result.Address, "IP address family. Currently only IPv4 networks are supported.")
	fs.StringVar(&result.Gateway, "gateway", result.Gateway, "Gateway address given by the DHCP service. Defaults to first address of the network if not given.")
	fs.StringVar(&dhcp, "dhcp", dhcp, "Toggles DHCP service for the network.")
	fs.StringVar(&dhcpDefRout, "dhcp-default-route", dhcpDefRout, "Defines if the gateway should be given as default route by DHCP. Defaults to yes on public networks, and no on other ones.")

	err = fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if dhcp != "" {
		switch dhcp {
		case "true":
			result.DHCP = upcloud.FromBool(true)
		case "false":
			result.DHCP = upcloud.FromBool(false)
		default:
			return nil, fmt.Errorf("%s is an invalid value for dhcp, it can be true of false", dhcp)
		}
	}

	if dhcpDefRout != "" {
		if dhcpDefRout == "false" {
			result.DHCPDefaultRoute = upcloud.FromBool(false)
		}
		if dhcpDefRout == "true" {
			result.DHCPDefaultRoute = upcloud.FromBool(true)
		}
		return nil, fmt.Errorf("%s is an invalid value for dhcp default rout, it can be true of false", dhcp)
	}

	if dns != "" {
		result.DHCPDns = strings.Split(dns, ",")
	}

	return result, nil
}