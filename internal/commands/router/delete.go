package router

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingRouter
}

// DeleteCommand creates the "delete router" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a router"),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	// TODO: reimplement
	// s.ArgCompletion(getRouterArgCompletionFunction(s.service))
}

// MaximumExecutions implements NewCommand.MaximumExecutions
func (s *deleteCommand) MaximumExecutions() int {
	return maxRouterActions
}

// Execute implements command.NewCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Deleting router %s", arg)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	err := exec.Network().DeleteRouter(&request.DeleteRouterRequest{UUID: arg})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.Marshaled{Value: nil}, nil
}
