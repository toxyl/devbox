# DevBox
...  is a command-line tool designed to simplify and streamline software development by running systems in a clean and isolated environment. With DevBox, you can create and manage isolated development boxes, enabling you to work with a variety of configurations without cluttering your primary system. DevBox is particularly useful for running multiple versions of a software stack, testing and developing software in a sandboxed environment, and simulating a network environment for testing purposes.

Please note that DevBox is still under development and breaking changes can occur. Furthermore, only Ubuntu DevBoxes have been tested on Ubuntu machines, and other distros might not work out of the box. Also, note that DevBoxes are always stored in the OS's temp directory. It is therefore vital to use the `store` or `workspace-store` commands if you want to preserve state and the `make` or `workspace-restore` commands to restore state.

# Installation
To install DevBox, run the following commands:

```bash
sudo apt install tmux tar xz-utils gzip
git clone https://github.com/toxyl/devbox
cd devbox
chmod +x build.sh
./build.sh
exec bash # to reload bash completions
```

# Commands
DevBox supports the following commands. For details on the usage and arguments of each command, run `sudo devbox <command>`.  
All commands can be completed using `<TAB>` and are context-aware. For example, when typing the name of a DevBox, pressing `<TAB>` will show a list of available DevBoxes. Similarly, when typing the name of a tarball file, `<TAB>` will complete the file path.

## DevBoxes
```bash
# Creates a `devbox`.
sudo devbox make                [devbox] [tarball]

# Enters a `devbox`.
sudo devbox enter               [devbox]

# Executes `command` in a `devbox`.
sudo devbox exec                [devbox] [command]

# Stores a `devbox` as `tarball`.
sudo devbox store               [devbox] [tarball]

# Removes a `devbox`.
sudo devbox destroy             [devbox]
```

## Workspaces
```bash
# Adds `devbox`es to the workspace `name`.
sudo devbox workspace-add       [name] [devbox_1<:delay>] <devbox_2> .. <devbox_n{:delay}> 

# Removes `devbox`es from the workspace `name`.
sudo devbox workspace-remove    [name] [devbox_1] <devbox_2> .. <devbox_n>

# Adds `IP`s to the workspace `name`.
sudo devbox workspace-ip-add    [name] [IP_1</prefix>] <IP_2> .. <IP_n{/prefix}>

# Removes `IP`s from the workspace `name`.
sudo devbox workspace-ip-remove [name] [IP_1</prefix>] <IP_2> .. <IP_n{/prefix}>

# Launches the workspace `name` in a tmux session.
sudo devbox workspace-launch    [name]

# Stores the workspace `name` as tarball.
sudo devbox workspace-store     [name] [tarball]

# Restores the workspace `name` from `tarball`.
sudo devbox workspace-restore   [name] [tarball]

# Completely removes the workspace `name`.
sudo devbox workspace-destroy   [name]
```

# Config
## DevBoxes
Each DevBox has a configuration file at `/config.yaml` which allows you to finetune the DevBox.  
Here's the default configuration:
```yaml
# The "limits" section defines the resource limits for the devboxes, including CPU usage, memory usage, swap space, and maximum number of processes.
limits:
  cpu: 0.1 # normalized percentage of host CPU (0..1)
  mem:
    hard: 0.10 # normalized percentage of host memory (0..1)
    soft: 0.75 # normalized percentage of `hard` (0..1) after which the OS might consider killing processes to reclaim memory
    swap: 0.10 # normalized percentage of `hard` (0..1) to grant as additional swap
  pids: 1024 # maximum number of processes allowed in the container

# The "options" section defines additional configuration options, such as mapping users and groups and binding all other devboxes.
options:
  map_users_and_groups: true # not all distributions need this, try setting it to false if you encouter startup errors
  bind_all: false # if enabled (true) all other devboxes will be mounted at /devboxes

# The "env" section sets environment variables for the devboxes, including the shell, terminal, home directory, and path.
env:
  shell: /bin/bash # this is the shell that will be used by the "enter" and "shell" commands
  term: xterm
  term_info: /usr/share/terminfo/
  home: "/root"
  path: "/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games:/usr/share/games:/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games:/usr/share/games:/usr/local/sbin:/usr/sbin:/sbin:/snap/bin:/snap/bin:/usr/sandbox/:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games:/usr/share/games:/usr/local/sbin:/usr/sbin:/sbin:/snap/bin:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games:/usr/share/games:/usr/local/sbin:/usr/sbin:/sbin:/snap/bin:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games"
  vars:
    MY_ENV_VAR: hello world
    MY_OTHER_ENV_VAR: 123

# The "commands" section defines scripts to execute before and after logging into the interactive shell.
commands:
  start: /usr/local/bin/start
  stop: /usr/local/bin/stop

# The "binds" section allows the binding of host directories to directories within the devboxes.
# Each entry is formatted as "src":"dst" pairs where 
# "src" is the directory within the DdevBox and 
# "dst" is the directory on the host to bind
binds: 
# For exmple, to bind /var/www/html on the host to /www on the DevBox you can enable this:
#  "/www": "/var/www/html"
```

## Workspaces
Each workspace has a configuration file at `/.workspace.yaml`which contains the configuration of all DevBoxes of the workspace. When launching a workspace the DevBox configuration from the workspace is written to `/config.yaml` of the DevBox, i.e. workspace settings always take precedence and replace DevBox settings. The workspace config will be updated with the configs of the DevBoxes once the last DevBox of the workspace is closed.

### Attaching And Detaching
After launching a workspace you'll be placed in a tmux session with all devboxes of the workspace. If you try to launch the same workspace again DevBox will just attach you to the existing session. If you want to detach all clients without closing the tmux session (e.g. to keep DevBox running in the background), you can use the `workspace-detach` command. 

## Storage
By default DevBox creates workspaces and DevBoxes in the temp directory (usually `/tmp`), so they will be cleaned up automatically on reboot. If you want to store workspaces and DevBoxes in a different location you can set that using the `storage-path-set` command. After switching the storage path you should run `exec bash` to update your Bash completions.  

Using this you can manage multiple sets of workspaces, for example:
```bash
# Switch storage path to /tmp/setA
sudo devbox storage-path-set /tmp/setA

# Create three new DevBoxes (will be stored in /tmp/setA/devbox/)
sudo devbox make bioA https://cloud-images.ubuntu.com/minimal/releases/bionic/release/ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz
sudo devbox make bioB https://cloud-images.ubuntu.com/minimal/releases/bionic/release/ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz
sudo devbox make bioC https://cloud-images.ubuntu.com/minimal/releases/bionic/release/ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz

# Create two new workspaces (will be stored in /tmp/setA/workspace)
sudo devbox workspace-add wsA bioA bioB
sudo devbox workspace-add wsB bioB bioC

tree -L 3 /tmp/setA/
# /tmp/setA
# ├── devbox
# │   ├── bioA
# │   │   ├── config.yaml
# │   │   ├── ...
# │   │   └── var
# │   ├── bioB
# │   │   ├── config.yaml
# │   │   ├── ...
# │   │   └── var
# │   ├── bioC
# │   │   ├── config.yaml
# │   │   ├── ...
# │   │   └── var
# │   └── ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz
# └── workspace
#     ├── wsA
#     │   ├── config.yaml
#     │   ├── bioA.tar.gz
#     │   └── bioB.tar.gz
#     └── wsA
#         ├── config.yaml
#         ├── bioB.tar.gz
#         └── bioC.tar.gz

# Switch storage path to /tmp/setB
sudo devbox storage-path-set /tmp/setB

# Create three new DevBoxes (will be stored in /tmp/setB/devbox/)
sudo devbox make bioD https://cloud-images.ubuntu.com/minimal/releases/bionic/release/ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz
sudo devbox make bioE https://cloud-images.ubuntu.com/minimal/releases/bionic/release/ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz
sudo devbox make bioF https://cloud-images.ubuntu.com/minimal/releases/bionic/release/ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz

# Create two new workspaces (will be stored in /tmp/setB/workspace)
sudo devbox workspace-add wsC bioD bioE
sudo devbox workspace-add wsD bioE bioF

tree -L 3 /tmp/setB/
# /tmp/setB
# ├── devbox
# │   ├── bioD
# │   │   ├── config.yaml
# │   │   ├── ...
# │   │   └── var
# │   ├── bioE
# │   │   ├── config.yaml
# │   │   ├── ...
# │   │   └── var
# │   ├── bioF
# │   │   ├── config.yaml
# │   │   ├── ...
# │   │   └── var
# │   └── ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz
# └── workspace
#     ├── wsC
#     │   ├── config.yaml
#     │   ├── bioD.tar.gz
#     │   └── bioE.tar.gz
#     └── wsD
#         ├── config.yaml
#         ├── bioE.tar.gz
#         └── bioF.tar.gz
```

# Usage Examples
## DevBoxes
```bash
# Make a DevBox
sudo devbox make bionic https://cloud-images.ubuntu.com/minimal/releases/bionic/release/ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz

# Enter the DevBox
sudo devbox enter bionic

# Execute a command in the DevBox
sudo devbox exec bionic 'echo "My hostname is: $(hostname)"'

# Store the DevBox
sudo devbox store bionic ~/bionic.tar.gz

# Destroy the DevBox
sudo devbox destroy bionic

# Restore the DevBox
sudo devbox make bionic ~/bionic.tar.gz 
```

## Workspaces
Before creating a workspace you first need to create one or more devboxes, as shown in the examples above.

```bash
# Create workspace with the DevBoxes 
# `bionic` (0 second startup delay) and 
# `kinetic` (5 second startup delay)
sudo devbox workspace-add my-workspace bionic kinetic:5

# Assign three IPs to the workspace
sudo devbox workspace-ip-add my-workspace 192.168.0.1 192.168.0.2 10.0.0.1

# Remove one IP from the workspace
sudo devbox workspace-ip-remove my-workspace 10.0.0.1

# Launch the workspace
sudo devbox workspace-launch my-workspace

# Store the workspace
sudo devbox workspace-store my-workspace ~/my-workspace.tar.gz

# Destroy the workspace
sudo devbox workspace-destroy my-workspace

# Restore the workspace
sudo devbox workspace-restore my-workspace ~/my-workspace.tar.gz
```

## Technologies Used
- **Linux namespaces** for isolation, with the exception of network interfaces which are shared with the host. Namespaces allow for creating a separate environment with its own process tree, file system, and more, ensuring that processes running within the DevBox are isolated from the host system and other DevBoxes.
- **Linux cgroups** (version 2) to set resource limits on the DevBoxes, such as CPU usage, memory usage, and the number of processes that can be run. These limits are controlled by systemd, which ensures that the DevBoxes do not exceed their allocated resources.
- **tar.gz** / **tar.xz** archives are used for storage, allowing for efficient and compressed storage of DevBox images and workspaces.
- DevBox uses **loopback aliases** to locally spoof IPs on the host that can be used by all DevBoxes. Aliases are created when launching a workspace and removed once it exits. 
- **YAML configuration files** are used to define the DevBoxes and workspaces. This allows for easy configuration and management of DevBoxes and workspaces.
- **Workspaces** are used to bundle DevBoxes together. A workspace is a collection of DevBoxes that are configured to work together. This can be useful for development teams that need to work on multiple components of an application.
- **tmux** is used to launch workspaces in a tiled session with mouse support, which allows for easy selection and resizing of panes.
- **Bind mounts** are used to map resources from the host to the DevBox. For example, one can map `/var/www/html` on the host to `/www` on the DevBox, allowing the DevBox to access files from the host's file system.
- **User and group mapping** can be used in two flavors: 
  - a) No users and groups are mapped to the DevBox, and only the `root` user is available. Every other user/group maps to `nobody`/`nogroup`.
  - b) All users and groups, except `root`, from the host are mapped to the DevBox.

# Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

# License
This project is licensed under the UNLICENSE - see the LICENSE file for details.
