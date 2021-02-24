package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/cli/internal/commands"
)

var (
	maxStorageActions = 10
	// CachedStorages stores the cached list of storages in order to not hit the service more than once
	// TODO: remove the cross-command dependencies
	CachedStorages []upcloud.Storage
)

const positionalArgHelp = "<UUID or Title...>"

// BaseStorageCommand creates the base "storage" command
func BaseStorageCommand() commands.Command {
	return &storageCommand{commands.New("storage", "Manage storages")}
}

type storageCommand struct {
	*commands.BaseCommand
}

func storageStateColor(state string) text.Colors {
	switch state {
	case upcloud.StorageStateOnline, upcloud.StorageStateSyncing:
		return text.Colors{text.FgGreen}
	case upcloud.StorageStateError:
		return text.Colors{text.FgHiRed, text.Bold}
	case upcloud.StorageStateMaintenance:
		return text.Colors{text.FgYellow}
	case upcloud.StorageStateCloning, upcloud.StorageStateBackuping:
		return text.Colors{text.FgHiMagenta, text.Bold}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

func importStateColor(state string) text.Colors {
	switch state {
	case "completed":
		return text.Colors{text.FgGreen}
	case "failed":
		return text.Colors{text.FgHiRed, text.Bold}
	case "pending", "importing":
		return text.Colors{text.FgYellow}
	case "cancelling":
		return text.Colors{text.FgHiMagenta, text.Bold}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

func matchStorages(storages []upcloud.Storage, searchVal string) []*upcloud.Storage {
	var r []*upcloud.Storage
	for _, storage := range storages {
		storage := storage
		if storage.Title == searchVal || storage.UUID == searchVal {
			r = append(r, &storage)
		}
	}
	return r
}

func searchStorage(storagesPtr *[]upcloud.Storage, service service.Storage, uuidOrTitle string, unique bool) ([]*upcloud.Storage, error) {
	if storagesPtr == nil || service == nil {
		return nil, fmt.Errorf("no storages or service passed")
	}
	storages := *storagesPtr
	if len(CachedStorages) == 0 {
		res, err := service.GetStorages(&request.GetStoragesRequest{})
		if err != nil {
			return nil, err
		}
		storages = res.Storages
		*storagesPtr = storages
	}
	matched := matchStorages(storages, uuidOrTitle)
	if len(matched) == 0 {
		return nil, fmt.Errorf("no storage with uuid, name or title %q was found", uuidOrTitle)
	}
	if len(matched) > 1 && unique {
		return nil, fmt.Errorf("multiple storages matched to query %q, use UUID to specify", uuidOrTitle)
	}
	return matched, nil
}

func searchAllStorages(terms []string, service service.Storage, unique bool) ([]string, error) {
	return commands.SearchResources(
		terms,
		func(id string) (interface{}, error) {
			return searchStorage(&CachedStorages, service, id, unique)
		},
		func(in interface{}) string { return in.(*upcloud.Storage).UUID })
}

// SearchSingleStorage returns exactly one storage where title or uuid matches uuidOrTitle
// TODO: remove the cross-command dependencies
func SearchSingleStorage(uuidOrTitle string, service service.Storage) (*upcloud.Storage, error) {
	matchedResults, err := searchStorage(&CachedStorages, service, uuidOrTitle, true)
	if err != nil {
		return nil, err
	}
	return matchedResults[0], nil
}

func getStorageDetailsUUID(in interface{}) string {
	return in.(*upcloud.StorageDetails).UUID
}

type storageRequest struct {
	ExactlyOne   bool
	BuildRequest func(storage string) (interface{}, error)
	Service      service.Storage
	Handler      ui.Handler
}

func (s storageRequest) send(args []string) (interface{}, error) {
	if s.ExactlyOne && len(args) != 1 {
		return nil, fmt.Errorf("single storage uuid is required")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("at least one storage uuid is required")
	}

	storages, err := searchAllStorages(args, s.Service, true)
	if err != nil {
		return nil, err
	}

	var requests []interface{}
	for _, storage := range storages {
		req, err := s.BuildRequest(storage)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return s.Handler.Handle(requests)
}

func getStorageArgumentCompletionFunction(s service.Storage) func(toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(toComplete string) ([]string, cobra.ShellCompDirective) {
		storages, err := s.GetStorages(&request.GetStoragesRequest{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range storages.Storages {
			vals = append(vals, v.UUID, v.Title)
		}
		return commands.MatchStringPrefix(vals, toComplete, false), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
}