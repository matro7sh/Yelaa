all:
	install
	compile

install:
	go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest
	# go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
	# go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest
	# go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest

compile:
	go build -o Yelaa main.go
