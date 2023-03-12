# config-reverter
Apply a config, revert it if the user doesn't interrupt the program within one minute

# Usage

The program takes three positional arguments: the new config to install, the actual path of the existing config where it must be written, and a command to run when the copying is done (restart a service, ...). Do note that a single command can be run, no chaining with `&&` or the like. My recommendation if you need one such thing would be to have a bash script with your commands in it, and then starting it (for example, `bash mycoolscript.sh` as the command). If the command fails for any reason, the configurtion is reverted.
An example:

```
config-reverter mycoolnetworkconfig /etc/network/interfaces 'systemctl restart networking'
```

This writes your changes to the interfaces configuration file, runs the command to restart networking, and waits for a minute. If it wasn't interrupted within the minute, then it reverts the changes and runs the command to restart networking again, hopefully restoring the access that you lost :sparkles:

# Building

Within the folder:

```
go build .
```

Quick note: this project is me trying to learn go. I probably butchered all code standards, sorry!
