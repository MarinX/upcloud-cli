package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/storage"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateServer(t *testing.T) {

	var Storage1 = upcloud.Storage{
		UUID:   UUID1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	var StorageDef = upcloud.Storage{
		UUID:   UUID2,
		Title:  "Debian GNU/Linux 10 (Buster)",
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	var Storage3 = upcloud.Storage{
		UUID:   UUID3,
		Title:  Title2,
		Access: "public",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   10,
		Tier:   "maxiops",
	}
	var storages = &upcloud.Storages{
		Storages: []upcloud.Storage{
			Storage1,
			StorageDef,
			Storage3,
		},
	}
	var serverDetailsMaint = upcloud.ServerDetails{
		Server: upcloud.Server{
			UUID:  UUID1,
			State: upcloud.ServerStateMaintenance,
		},
		VideoModel: "vga",
		Firewall:   "off",
	}

	serverDetailsStarted := serverDetailsMaint
	serverDetailsStarted.State = upcloud.ServerStateStarted

	for _, test := range []struct {
		name            string
		args            []string
		createServerReq request.CreateServerRequest
		error           string
	}{
		{
			name: "use default values",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--password-delivery", "email",
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				PasswordDelivery: "email",
				LoginUser:        &request.LoginUser{CreatePassword: "yes"},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "",
					Storage: StorageDef.UUID,
					Title:   "example.com-osDisk",
					Size:    10,
					Tier:    upcloud.StorageTierMaxIOPS,
					Type:    upcloud.StorageTypeDisk,
				}},
			},
		},
		{
			name: "server OS set, size larger than the minimum",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--os", Storage1.UUID,
				"--os-storage-size", "100",
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				PasswordDelivery: "none",
				LoginUser:        &request.LoginUser{CreatePassword: "yes"},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "",
					Storage: Storage1.UUID,
					Title:   "example.com-osDisk",
					Size:    100,
					Tier:    upcloud.StorageTierMaxIOPS,
					Type:    upcloud.StorageTypeDisk,
				}},
			},
		},
		{
			name: "flags mapped to the correct field",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--avoid-host", "1234",
				"--host", "5678",
				"--boot-order", "cdrom,network",
				"--user-data", "example.com",
				"--cores", "12",
				"--memory", "4096",
				"--plan", "custom",
				"--password-delivery", "sms",
				"--simple-backup", "00,monthlies",
				"--time-zone", "EET",
				"--video-model", "VM",
				"--firewall",
				"--metadata",
				"--username", "johndoe",
				"--remote-access-enabled",
				"--remote-access-type", upcloud.RemoteAccessTypeVNC,
				"--remote-access-password", "secret",
			},
			createServerReq: request.CreateServerRequest{
				Hostname:             "example.com",
				Title:                "test-server",
				Zone:                 "uk-lon1",
				AvoidHost:            1234,
				Host:                 5678,
				BootOrder:            "cdrom,network",
				UserData:             "example.com",
				CoreNumber:           12,
				MemoryAmount:         4096,
				Plan:                 "custom",
				PasswordDelivery:     "sms",
				SimpleBackup:         "00,monthlies",
				TimeZone:             "EET",
				VideoModel:           "VM",
				Firewall:             "on",
				Metadata:             upcloud.FromBool(true),
				RemoteAccessEnabled:  upcloud.FromBool(true),
				RemoteAccessType:     upcloud.RemoteAccessTypeVNC,
				RemoteAccessPassword: "secret",
				LoginUser:            &request.LoginUser{CreatePassword: "yes", Username: "johndoe"},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "",
					Storage: StorageDef.UUID,
					Title:   "example.com-osDisk",
					Size:    10,
					Tier:    upcloud.StorageTierMaxIOPS,
					Type:    upcloud.StorageTypeDisk,
				}},
			},
		},
		{
			name: "multiple storages",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--password-delivery", "email",
				"--storage", fmt.Sprintf("action=create,address=virtio,type=disk,size=20,title=new-storage"),
				"--storage", fmt.Sprintf("action=clone,storage=%s,title=three-clone", Storage3.Title),
				"--storage", fmt.Sprintf("action=attach,storage=%s,type=cdrom", Storage1.Title),
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				PasswordDelivery: "email",
				LoginUser:        &request.LoginUser{CreatePassword: "yes"},
				StorageDevices: request.CreateServerStorageDeviceSlice{
					request.CreateServerStorageDevice{
						Action:  "clone",
						Address: "",
						Storage: StorageDef.UUID,
						Title:   "example.com-osDisk",
						Size:    10,
						Tier:    upcloud.StorageTierMaxIOPS,
						Type:    upcloud.StorageTypeDisk,
					},
					request.CreateServerStorageDevice{
						Action:  "create",
						Address: "virtio",
						Title:   "new-storage",
						Size:    20,
						Type:    upcloud.StorageTypeDisk,
					},
					request.CreateServerStorageDevice{
						Action:  "clone",
						Storage: Storage3.UUID,
						Title:   "three-clone",
					},
					request.CreateServerStorageDevice{
						Action:  "attach",
						Storage: Storage1.UUID,
						Type:    upcloud.StorageTypeCDROM,
					},
				},
			},
		},
		{
			name: "with networks",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--password-delivery", "email",
				"--network", "family=IPv4,type=utility",
				"--network", "family=IPv6,type=public",
				"--network", "family=IPv6,type=private",
			},
			createServerReq: request.CreateServerRequest{
				VideoModel:       "vga",
				TimeZone:         "UTC",
				Plan:             "1xCPU-2GB",
				Hostname:         "example.com",
				Title:            "test-server",
				Zone:             "uk-lon1",
				PasswordDelivery: "email",
				LoginUser:        &request.LoginUser{CreatePassword: "yes"},
				StorageDevices: request.CreateServerStorageDeviceSlice{request.CreateServerStorageDevice{
					Action:  "clone",
					Address: "",
					Storage: StorageDef.UUID,
					Title:   "example.com-osDisk",
					Size:    10,
					Tier:    upcloud.StorageTierMaxIOPS,
					Type:    upcloud.StorageTypeDisk,
				}},
				Networking: &request.CreateServerNetworking{Interfaces: request.CreateServerInterfaceSlice{
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: "IPv4"}},
						Type:        upcloud.NetworkTypeUtility,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: "IPv6"}},
						Type:        upcloud.NetworkTypePublic,
					},
					request.CreateServerInterface{
						IPAddresses: request.CreateServerIPAddressSlice{request.CreateServerIPAddress{Family: "IPv6"}},
						Type:        upcloud.NetworkTypePrivate,
					},
				}},
			},
		},
		{
			name: "networks type missing",
			args: []string{
				"--hostname", "example.com",
				"--title", "test-server",
				"--zone", "uk-lon1",
				"--network", "family=IPv4,type=utility",
				"--network", "family=IPv6,type=public",
				"--network", "family=IPv6",
				"--password-delivery", "sms",
			},
			error: "network type is required",
		},
		{
			name: "hostname is missing",
			args: []string{
				"--title", "title",
				"--zone", "zone",
			},
			error: "hostname and zone are both required",
		},
		{
			name: "zone is missing",
			args: []string{
				"--title", "title",
				"--hostname", "hostname",
			},
			error: "hostname and zone are both required",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedServers = nil
			storage.CachedStorages = nil

			mServerService := MockServerService{}
			mServerService.On("CreateServer", &test.createServerReq).Return(&serverDetailsMaint, nil)
			mServerService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: serverDetailsMaint.UUID}).Return(&serverDetailsStarted, nil)

			mStorageService := MockStorageService{}
			mStorageService.On("GetStorages", mock.Anything).Return(storages, nil)

			cc := commands.BuildCommand(CreateCommand(&mServerService, &mStorageService), nil, config.New(viper.New()))
			cc.SetFlags(test.args)

			_, err := cc.MakeExecuteCommand()([]string{})

			if test.error != "" {
				assert.Equal(t, test.error, err.Error())
			} else {
				mStorageService.AssertNumberOfCalls(t, "GetStorages", 1)
				mServerService.AssertNumberOfCalls(t, "CreateServer", 1)
				mServerService.AssertNumberOfCalls(t, "GetServerDetails", 1)
			}
		})
	}

}