package server

import (
	"fmt"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/mapper"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

// StartCommand creates the "server start" command
func StartCommand() commands.NewCommand {
	return &startCommand{
		BaseCommand: commands.New("start", "Start a server"),
	}
}

type startCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (s *startCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.Config()))
}
func (s *startCommand) ArgumentMapper() (mapper.Argument, error) {
	return mapper.CachingServer(s.Config().Service.Server())
}

func (s *startCommand) Execute(exec commands.Executor, uuid string) (output.Command, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("starting server %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))

	res, err := svc.StartServer(&request.StartServerRequest{
		UUID: uuid,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	if s.Config().GlobalFlags.Wait {
		logline.SetMessage(fmt.Sprintf("%s: waiting to start", msg))
		if err := exec.WaitFor(serverStateWaiter(uuid, upcloud.ServerStateStarted, msg, svc, logline), s.Config().ClientTimeout()); err != nil {
			return nil, err
		}

		logline.SetMessage(fmt.Sprintf("%s: server started", msg))
	} else {
		logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
	}

	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
