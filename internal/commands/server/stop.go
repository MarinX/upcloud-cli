package server

import (
	"fmt"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/mapper"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

// StopCommand creates the "server stop" command
func StopCommand() commands.NewCommand {
	return &stopCommand{
		BaseCommand: commands.New("stop", "Stop a server"),
	}
}

type stopCommand struct {
	*commands.BaseCommand
	StopType string
}

// InitCommand implements Command.InitCommand
func (s *stopCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.Config()))

	//XXX: findout what to do with risky params (timeout actions)
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.StopType, "type", defaultStopType, "The type of stop operation. Available: soft, hard")
	s.AddFlags(flags)
}

func (s *stopCommand) ArgumentMapper() (mapper.Argument, error) {
	return mapper.CachingServer(s.Config().Service.Server())
}

func (s *stopCommand) Execute(exec commands.Executor, uuid string) (output.Command, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("stoping server %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))

	res, err := svc.StopServer(&request.StopServerRequest{
		UUID:     uuid,
		StopType: s.StopType,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	if s.Config().GlobalFlags.Wait {
		logline.SetMessage(fmt.Sprintf("%s: waiting to stop", msg))
		if err := exec.WaitFor(serverStateWaiter(uuid, upcloud.ServerStateStopped, msg, svc, logline), s.Config().ClientTimeout()); err != nil {
			return nil, err
		}

		logline.SetMessage(fmt.Sprintf("%s: server stoped", msg))
	} else {
		logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
	}

	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
