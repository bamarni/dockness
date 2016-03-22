# dockness

DNS for Docker machines, allows to access them with the following domain format : `{machine}.docker`.

## How it works

It spins up a simplistic DNS server, only listening for questions about A records.

Behind the scene it will run `docker-machine ip {machine}` in order to resolve the IP address of a given machine.

## Installation

### Prebuilt binaries

Prebuilt binaries are available in the [releases](https://github.com/bamarni/dockness/releases).

### From source (requires Go)

    go get github.com/bamarni/dockness

## Usage

    dockness [options...]

    Options:
      -tld          Top-level domain to use (defaults to "docker")
      -port         Port to listen on (defaults to "53")
      -server-only  Server only, doesn't try to create a resolver configuration file
      -user         Execute the "docker-machine ip" command as a different user (defaults to "SUDO_USER")

### Mac OSX

To develop on Mac you probably have a local VM, using VirtualBox for example.
However this machine gets assigned a dynamic IP address.

The program can be up and running in one command :

    > sudo dockness
    2016/02/18 10:39:52 Creating configuration file at /etc/resolver/docker...
    2016/02/18 10:39:52 Listening on :53...

### Linux

Even though Linux users might not need a development VM, it can be useful for remote machines
(eg. `ssh staging2.docker`, ...).

Here and unlike Mac OSX, there is no quick trick to make it work out of the box. Linux distributions come
with more advanced DNS management where manually tweaking `resolv.conf` is usually not an option.

What should work in most cases is to use [Dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html),
which provides a lightweight DNS server.

If you have it installed, you can run :

    > dockness -port 10053
    2016/02/18 10:40:43 Listening on :10053...

Then let Dnsmasq know about the resolver by running those commands :

    > echo "server=/docker/127.0.0.1#10053" | sudo tee -a /etc/dnsmasq.conf
    > sudo /etc/init.d/dnsmasq restart

## Configure dockness as a service

As it's not very convenient to run the program manually in a terminal, you can instead set it as a service.
Doing so, it will be running in the background automatically when booting your computer.

### Mac OSX

Run this command : `echo "nameserver 127.0.0.1\nport 10053" | sudo tee /etc/resolver/docker`

It tells your Mac that the resolver for `.docker` TLD listens locally on port 10053.

You can now create the appropriate service configuration file at `~/Library/LaunchAgents/local.dockness.plist` :

``` xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>Disabled</key>
        <false/>
        <key>EnvironmentVariables</key>
        <dict>
                <key>PATH</key>
                <string>/usr/local/bin</string>
        </dict>
        <key>Label</key>
        <string>local.dockness</string>
        <key>ProgramArguments</key>
        <array>
                <string>/path/to/dockness</string>
                <string>-server-only</string>
                <string>-port</string>
                <string>10053</string>
        </array>
        <key>RunAtLoad</key>
        <true/>
</dict>
</plist>
```

You'll have to adapt 2 parameters :
- `/path/to/dockness`
- `/usr/local/bin`, which is the directory containing your `docker-machine` executable

Finally, the service can be enabled :

    launchctl load ~/Library/LaunchAgents/local.dockness.plist

### Linux

Here again, it will depend on your Linux distribution.
We'll take as example [Systemd](https://freedesktop.org/wiki/Software/systemd/),
which is nowadays the default init system in Ubuntu/Debian.

Create the following file at `/etc/systemd/system/dockness.service`:

    TODO

