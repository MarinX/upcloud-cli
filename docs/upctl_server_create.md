## upctl server create

Create a server

```
upctl server create [flags]
```

### Examples

```
upctl server create --title myapp --zone fi-hel1 --hostname myapp --password-delivery email
upctl server create --wait --title myapp --zone fi-hel1 --hostname myapp --password-delivery email
upctl server create --title "My Server" --zone fi-hel1 --hostname myapp --password-delivery email
upctl server create --zone fi-hel1 --hostname myapp --password-delivery email --plan 2xCPU-4GB
upctl server create --zone fi-hel1 --hostname myapp --password-delivery email --plan custom --cores 2 --memory 4096
upctl server create --zone fi-hel1 --hostname myapp --password-delivery email --os "Debian GNU/Linux 10 (Buster)"
upctl server create --zone fi-hel1 --hostname myapp --ssh-keys /path/to/publickey --network type=private,network=037a530b-533e-4cef-b6ad-6af8094bb2bc,ip-address=10.0.0.1
```

### Options

```
      --avoid-host int                  Use this to make sure VMs do not reside on specific host. Refers to value from host -attribute. Useful when building HA-environments.
      --boot-order string               The boot device order, disk / cdrom / network or comma separated combination.
      --cores plan                      Number of cores. Only allowed if plan option is set to "custom".
      --create-password[=true]          Create an admin password.
      --enable-firewall[=true]          Enable firewall.
      --enable-metadata[=true]          Enable metadata service.
      --enable-remote-access[=true]     Enable remote access.
      --host int                        Use this to start a VM on a specific private cloud host. Refers to value from host -attribute. Only available in private clouds.
      --hostname string                 Server hostname.
      --memory plan                     Memory amount in MiB. Only allowed if plan option is set to "custom".
      --network stringArray             A network interface for the server, multiple can be declared.
                                        Usage: --network family=IPv4,type=public
                                        
                                        --network type=private,network=037a530b-533e-4cef-b6ad-6af8094bb2bc,ip-address=10.0.0.1
      --os string                       Server OS to use (will be the first storage device). The value should be title or UUID of an either public or private template. Set to empty to fully customise the storages. (default "Ubuntu Server 20.04 LTS (Focal Fossa)")
      --os-storage-size os              OS storage size in GiB. This is only applicable if os is also set. Zero value makes the disk equal to the minimum size of the template.
      --password-delivery string        Defines how password is delivered. Available: email, sms (default "none")
      --plan cores                      Server plan name. See "server plans" command for valid plans. Set to "custom" and use cores and `memory` options for flexible plan. (default "1xCPU-2GB")
      --remote-access-password string   Defines the remote access password.
      --remote-access-type string       Set a remote access type. Available: vnc, spice
      --simple-backup string            Simple backup rule. Format (HHMM,{dailies,weeklies,monthlies}). Example: 2300,dailies
      --ssh-keys strings                Add one or more SSH keys to the admin account. Accepted values are SSH public keys or filenames from where to read the keys.
      --storage stringArray             A storage connected to the server, multiple can be declared.
                                        Usage: --storage action=attach,storage=01000000-0000-4000-8000-000020010301,type=cdrom
      --time-zone string                Time zone to set the RTC to. (default "UTC")
      --title string                    A short, informational description.
      --user-data string                Defines URL for a server setup script, or the script body itself.
      --username string                 Admin account username.
      --video-model string              Video interface model of the server. Available: vga, cirrus (default "vga")
      --wait[=true]                     Wait for server to be in started state before returning.
      --zone string                     Zone where to create the server.
  -h, --help                            help for create
```

### Options inherited from parent commands

```
  -t, --client-timeout duration   Client timeout to use in API calls.
      --config string             Configuration file path.
      --debug                     Print out more verbose debug logs.
      --force-colours[=true]      Force coloured output despite detected terminal support.
      --no-colours[=true]         Disable coloured output despite detected terminal support. Colours can also be disabled by setting NO_COLOR environment variable.
  -o, --output string             Output format (supported: json, yaml and human) (default "human")
```

### SEE ALSO

* [upctl server](upctl_server.md)	 - Manage servers

###### Auto generated by spf13/cobra on 17-Oct-2022
