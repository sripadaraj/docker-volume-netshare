## Installation

#### From Source

```
$ go get github.com/ContainX/docker-volume-netshare
$ go build
```

#### From Binaries

Binaries are available through GitHub releases.  You can download the appropriate binary, package and version from the [Releases](https://github.com/ContainX/docker-volume-netshare/releases) page

#### On Ubuntu / Debian

The method below will install the sysvinit and /etc/default options that can be overwritten during service start.

1. Install the Package

```
  $ wget https://github.com/ContainX/docker-volume-netshare/releases/download/v0.20/docker-volume-netshare_0.20_amd64.deb
  $ sudo dpkg -i docker-volume-netshare_0.20_amd64.deb
```

2. Modify the startup options in `/etc/default/docker-volume-netshare`
3. Start the service `service docker-volume-netshare start`


## Usage


### Launching in Samba/CIFS mode


**1. Run the plugin - can be added to systemd or run in the background**

```
  $ sudo docker-volume-netshare cifs --username user --password pass --domain domain --security security
```

**2. Launch a container**

```
  // In CIFS the "//" is omitted and handled by netshare
  $ docker run -it --volume-driver=cifs -v cifshost/share:/mount ubuntu /bin/bash
```

**1. Run the plugin - can be added to systemd or run in the background**

```
  $ sudo docker-volume-netshare cifs
```

**2. Create a Volume**

This will create a new volume via the Docker daemon which will call `Create` in netshare passing in the corresponding user, pass and domain info.

```
  $ docker volume create -d cifs --name cifshost/share --opt username=user --opt password=pass --opt domain=domain --opt security=security --opt fileMode=0777 --opt dirMode=0777
```

**3. Launch a container**

```
  // cifs/share matches the volume as defined in Step #2 using docker volume create
  $ docker run -it -v cifshost/share:/mount ubuntu /bin/bash
```
