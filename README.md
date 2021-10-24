# Yelaa

Obtain a clean-cut architecture at the launch of a mission and make some tests

# Preview

![preview](img/preview.png)

## How to use 
>-s is optionnal
You can run `yelaa create -c <client> -s <PathToSharedFolder>`

## How to run scan 

`yelaa scan -target <PathToTargetFile>`

## Use http proxy

`yelaa scan -p http://localhost:8080 -target ./targets.txt`

## Help 

``` 
./yelaaa -h
Obtain a clean-cut architecture at the launch of a mission and make some tests

Usage:
  create -c [client name] [flags]
  create [command]

Available Commands:
  help        Help about any command
  scan        It will run Nuclei templates, sslscan, dirsearch and more.

Flags:
  -c, --client string         Client name
  -e, --excludedType string   excluded type
  -h, --help                  help for create
  -s, --shared string         path to shared folder

Use "create [command] --help" for more information about a command.
``` 

>this script will create a default structure, as well as a cherytree database with payloads for external testing and useful commands for internal testing.