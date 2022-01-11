# Yelaa

Obtain a clean-cut architecture at the launch of a mission and make some tests

# Requirements

You need to have theses binaries in your path:
```
nuclei, google-chrome
```

# How to install

Manually :
```bash
git clone https://github.com/CMEPW/Yelaa.git
cd Yelaa
make install
make compile
```

Or if you have set your GO path and all the requirements installed :
```bash
go install github.com/CMEPW/Yelaa@latest
```

# How to use 
>-s is optionnal
You can run `Yelaa create -c <client> -s <PathToSharedFolder>`

## How to run scan 

`Yelaa scan -target <PathToTargetFile>`

## Use http proxy

`Yelaa scan -p http://localhost:8080 -target ./targets.txt`

>Flag `-k` is available to skip tls configuration

## How to run osint on a domain

`Yelaa osint -d example.com`

To run osint command on several domains run `Yelaa osint -t domains.txt`

## How to run httpx then gowitness

`Yelaa checkAndScreen -t domains.txt`

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
  help            Help about any command
  osint           Run subfinder, dnsx and httpx to find ips and subdomains of a specific domain
  scan            It will run Nuclei templates, sslscan, dirsearch and more.

Flags:
  -c, --client string         Client name
  -e, --excludedType string   excluded type
  -h, --help                  help for create
  -k, --insecure              Allow insecure certificate
  -p, --proxy string          Add HTTP proxy
  -s, --shared string         path to shared folder

Use "create [command] --help" for more information about a command.

``` 

>this script will create a default structure using `create` command, as well as a cherytree database with payloads for external testing and useful commands for internal testing

# Contributors

[darkweak](https://github.com/darkweak)
[jenaye](https://github.com/jenaye)
[jarrault](https://github.com/jarrault)
