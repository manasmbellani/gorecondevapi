# gorecondevapi
Recon DEV API for retrieving subdomains and IPs via recon.dev website

## Install
```
go get -u -v github.com/manasmbellani/gorecondevapi/src/gorecondevapi
```

## Run
To scan IPs, domains via `recon.dev` API, execute the following command:
```
gorecondevapi -apiKey $RECON_DEV_API_KEY -domain $DOMAIN > /tmp/assets.txt
```

Domains and IPs are written one per line in `/tmp/assets.txt`