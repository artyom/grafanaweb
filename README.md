# Grafanaweb

This is a self-contained web server and reverse proxy to host [Grafana][] used
with [InfluxDB][] backend. It is presumed that both InfluxDB and this server
run on the same host or in same safe network.

[Grafana]: http://grafana.org
[InfluxDB]: http://influxdb.com

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

Then you can serve your setup via http like this:

	grafanaweb -listen=:80 -root=/path/to/grafana/files

For running in HTTPS mode you'll need certificate and key in PEM format. If you
don't have one, try this command:

	go run $(go env GOROOT)/src/pkg/crypto/tls/generate_cert.go -host=YOUR_HOST_IP_OR_HOSTNAME

Then run server in ssl mode with self-signed certificate created:

	grafanaweb -listen=:443 -ssl -key /path/to/key.pem -cert /path/to/cert.pem
