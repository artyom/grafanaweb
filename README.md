# Grafanaweb

This is a self-contained web server and reverse proxy to host [Grafana][] used
with [InfluxDB][] backend. It is presumed that both InfluxDB and this server
run on the same host or in same safe network.

[Grafana]: http://grafana.org
[InfluxDB]: http://influxdb.com

## Installing

To install `grafanaweb` and `bcryptpasswd` (see below) commands, execute:

	go get -v github.com/artyom/grafanaweb/...

`grafanaweb` usage:

	Usage of grafanaweb:
	  -auth=false: use basic authentication (needs -authfile to be set)
	  -authfile="": path to file with "username:brypt_hash" credential records
	  -cert="": path to cert.pem file (only if -ssl used)
	  -key="": path to key.pem file (only if -ssl used)
	  -listen="127.0.0.1:8080": address to listen
	  -prefix="/db/": url prefix to db (should finish with /)
	  -proxy="http://127.0.0.1:8086": url to proxy prefix matched requests
	  -root="/var/lib/grafana": path to grafana files
	  -ssl=false: use https instead of http (needs both -cert and -key options set)

## Configuring Grafana

Consider you have InfluxDB API listening on 127.0.0.1:8086. You can place
grafana files on the same host, modify its `config.js` file to something like
this (note relative URLs):

	datasources: {
		influxdb: {
			type: 'influxdb',
			url: "/db/test",
			username: 'admin',
			password: 'admin',
		},
		grafana: {
			type: 'influxdb',
			url: "/db/grafana",
			username: 'admin',
			password: 'admin',
			grafanaDB: true
		},
	},

## Basic Setup (development, trusted network, etc.)

You can host your setup via http like this:

	grafanaweb -listen=:80 -root=/path/to/grafana/files

## Hardened setup (internet-facing host): https + authentication

For running in HTTPS mode you'll need certificate and key in PEM format. If you
don't have one, try this command:

	go run $(go env GOROOT)/src/pkg/crypto/tls/generate_cert.go -host=YOUR_HOST_IP_OR_HOSTNAME

Create credentials file, entering password for "admin" user (enter your
password to stdin):

	bcryptpasswd admin >> /path/to/htpasswd.bcrypt

Then run server in ssl mode with self-signed certificate created and http basic authentication enabled:

	grafanaweb -listen=:443 \
		-ssl -key=/path/to/key.pem -cert=/path/to/cert.pem \
		-auth -authfile=/path/to/htpasswd.bcrypt
