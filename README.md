# Yelaa

Obtain a clean-cut architecture at the launch of a mission and make some tests

# Requirements

You need to have the chrome binary in your path:
```
google-chrome
```

# How to install

Manually :
```bash
git clone https://github.com/CMEPW/Yelaa.git
cd Yelaa
make compile
```

Or if you have set your GO path and all the requirements installed :
```bash
go install github.com/CMEPW/Yelaa@latest
```

In a Docker-container:
```bash
# Build docker container
make docker

# Or
docker build -t yelaa \
		--build-arg USER_ID=$(id -g) \
		--build-arg GROUP_ID=$(id -u) \
		.

# create a file with your target
echo "Some web addresses..." > targets.txt

# run the container like so
docker run \
    --security-opt seccomp=unconfined \
    -v $PWD:/home/yelaa_user \
    yelaa \
    checkAndScreen -t /home/yelaa_user/targets.txt
```

In Kali:
```bash

wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb

sudo  apt install ./google-chrome-stable_current_amd64.deb

wget https://github.com/CMEPW/Yelaa/releases/download/v1.7.1/Yelaa_1.7.1_Linux_x86_64.tar.gz

tar -xvf Yelaa_1.7.1_Linux_x86_64.tar.gz
./Yelaa -h
```

# How to use
>-s is optional
You can run `Yelaa create -c <client> -s <PathToSharedFolder>`

## How to run scan

`Yelaa scan -t <PathToTargetFile>`

## Use http / socks proxy

```bash
# using a http proxy
Yelaa scan -p http://localhost:8080 -target ./targets.txt`

# or, socks5 proxy
Yelaa scan -p socks5://localhost:9050 -target ./targets.txt`
```

>Flag `-k` is available to skip tls configuration

>Please prefer using socks5 as much as possible, as socks4 can fail depending on your go version

## How to run osint on a domain

`Yelaa osint -t ./targets.txt -p http://localhost:8080 --path /tmp`

or

`./Yelaa osint -d <domain>`

This command use the default browser to open the dork page
To run osint command on several domains run `Yelaa osint -t targets.txt`

## How to run httpx then gowitness

`Yelaa checkAndScreen -t domains.txt`

## Low fruits : Infrastructure Penetration Testing

```bash
# run scan on ports 80, 443, 8080 & 8443
nmap -T4 -Pn -p 80,443,8080,8443 --open -oA EvilCorp-24 192.168.1.0/24

# fetch tcp open ports & put them in web-targets.txt
cat *.gnmap | grep -i "open/tcp" | cut -d " " -f2 | sort -u > web-targets.txt

# run check-and-screen to quickly map infra
./Yelaa checkAndScreen -t ./web-targets.txt
```

## Help

```
Yelaa -h
 __   __         _
 \ \ / /   ___  | |   __ _    __ _
  \ V /   / _ \ | |  / _` |  / _` |
   | |   |  __/ | | | (_| | | (_| |
   |_|    \___| |_|  \__,_|  \__,_|
Obtain a clean-cut architecture at the launch of a mission and make some tests

Usage:
  create -c [client name] [flags]
  create [command]

Available Commands:
  checkAndScreen Run httpx and gowitness
  completion     Generate the autocompletion script for the specified shell
  help           Help about any command
  osint          Run subfinder, dnsx and httpx to find ips and subdomains of a specific domain
  scan           It will run Nuclei templates, dirsearch and more.

Flags:
  -c, --client string         Client name
      --dry-run               Run in dry-run mode
  -e, --excludedType string   excluded type
  -h, --help                  help for create
  -k, --insecure              Allow insecure certificate
      --nuclei                Enable nuclei with the command
      --path string           Output path (default "/home/$USER/.yelaa")
  -p, --proxy string          Add HTTP proxy
      --rate-limit int32      Rate limitation for nuclei and gobuster (default 100)
  -s, --shared string         path to shared folder

Use "create [command] --help" for more information about a command.
All temporary file have been succesfully removed
```

> This script will create a default structure using `create` command, as well as a cherytree database with payloads for external testing and useful commands for internal testing

## run with Proxychains

> this is not the recommanded way to use a proxy! You can just specify a proxy with the `-p` option!

If you *must* run Yelaa through Proxychains, it is possible but will require a bit of tweaking.
The reason for that is that Yelaa is statically compiled, and `Proxychains` uses `LD_PRELOAD` tricks to set a proxy.
You will have to compile Yelaa dynamically, using `gcc-go` (you will have to [install it yourself](https://go.dev/doc/install/gccgo) before compiling):

```bash
git clone https://github.com/CMEPW/Yelaa.git

cd Yelaa

make dynamic
```

# Contributors

| [<img src="https://github.com/darkweak.png?size=85" width=85><br><sub>Darkweak</sub>](https://github.com/darkweak) | [<img src="https://github.com/jenaye.png?size=85" width=85><br><sub>Mike Houziaux</sub>](https://github.com/jenaye) | [<img src="https://github.com/jarrault.png?size=85" width=85><br><sub>Julien</sub>](https://github.com/jarrault) | [<img src="https://github.com/TomChv.png?size=85" width=85><br><sub>Tom Chauveau</sub>](https://github.com/TomChv) | [<img src="https://github.com/bogdzn.png?size=85" width=85><br><sub>bogdan</sub>](https://github.com/bogdzn)| [<img src="h[ttps://github.com/bogdzn](https://github.com/VidsSkids.png?size=85" width=85><br><sub>VidsSkids</sub>]([https://github.com/bogdzn](https://github.com/VidsSkids))
| :---: | :---: | :---: | :---: | :---: | :---: |
