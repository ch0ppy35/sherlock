package testing

import (
	"sync"
	test "testing" // Yeah this is dumb

	"github.com/ch0ppy35/sherlock/internal/dns"
)

// TODO
func Test_queryDNSForHost(t *test.T) {
	type args struct {
		host    string
		server  string
		client  dns.TinyDNSClient
		results map[string]*dns.DNSRecords
		errors  map[string]error
		mu      *sync.Mutex
		wg      *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *test.T) {
			queryDNSForHost(tt.args.host, tt.args.server, tt.args.client, tt.args.results, tt.args.errors, tt.args.mu, tt.args.wg)
		})
	}
}
