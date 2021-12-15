# Yelaa

Obtain a clean-cut architecture at the launch of a mission and make some tests

# Requirements

Go theses binary in your path : 
```
nuclei, dirsearch, testssl, subfinder, dnsx
```


> You can set alias like this `dirsearch='python /home/jenaye/softs/dirsearch/dirsearch.py'` 

## How to use 
>-s is optionnal
You can run `yelaa create -c <client> -s <PathToSharedFolder>`

## How to run scan 

`yelaa scan -target <PathToTargetFile>`

## Use http proxy

`yelaa scan -p http://localhost:8080 -target ./targets.txt`

>Flag `-k` is available to skip tls configuration

## Help 

``` 
./yelaa -h 
 __   __         _
 \ \ / /   ___  | |
  \ V /   / _ \ | |  / _` |  / _` |
   | |   |  __/ | | | (_| | | (_| |
   |_|    \___| |_|  \__,_|  \__,_|
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

# Preview command create

![preview](img/preview.png)

# Preview command scan 

![pou](img/preview-scan.png)

# Contributors

[darkweak](https://github.com/darkweak)
[jenaye](https://github.com/jenaye)
[jarrault](https://github.com/jarrault)