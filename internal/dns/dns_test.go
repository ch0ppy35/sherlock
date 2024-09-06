package dns

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestQueryDNSRecord(t *testing.T) {
	tests := []struct {
		name          string
		domain        string
		qtype         uint16
		expectedError bool
		mockResponse  *dns.Msg
		mockError     error
		expectedSet   []dns.RR
	}{
		{
			name:   "Valid A record query",
			domain: "example.com",
			qtype:  dns.TypeA,
			expectedSet: []dns.RR{
				&dns.A{Hdr: dns.RR_Header{Name: "example.com."}, A: net.ParseIP("10.0.0.1")},
			},
			mockResponse: &dns.Msg{
				Answer: []dns.RR{
					&dns.A{Hdr: dns.RR_Header{Name: "example.com."}, A: net.ParseIP("10.0.0.1")},
				},
			},
		},
		{
			name:   "Valid MX record query",
			domain: "example.com",
			qtype:  dns.TypeMX,
			expectedSet: []dns.RR{
				&dns.MX{Hdr: dns.RR_Header{Name: "example.com."}, Mx: "mail.example.com.", Preference: 10},
			},
			mockResponse: &dns.Msg{
				Answer: []dns.RR{
					&dns.MX{Hdr: dns.RR_Header{Name: "example.com."}, Mx: "mail.example.com.", Preference: 10},
				},
			},
		},
		{
			name:          "Query returns no answers",
			domain:        "example.com",
			qtype:         dns.TypeA,
			expectedError: false,
			mockResponse:  &dns.Msg{},
		},
		{
			name:          "Query returns an error",
			domain:        "example.com",
			qtype:         dns.TypeA,
			expectedError: true,
			mockError:     fmt.Errorf("network error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockIDNSClient{
				MockExchange: func(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
					if tt.mockError != nil {
						return nil, 0, tt.mockError
					}
					return tt.mockResponse, 0, nil
				},
			}

			var receivedRecords []dns.RR
			setter := func(rr dns.RR) {
				receivedRecords = append(receivedRecords, rr)
			}

			err := QueryDNSRecord(client, tt.domain, "8.8.8.8", tt.qtype, setter)
			if (err != nil) != tt.expectedError {
				t.Errorf("QueryDNSRecord() error = %v, expectedError %v", err, tt.expectedError)
			}

			if !reflect.DeepEqual(receivedRecords, tt.expectedSet) {
				t.Errorf("QueryDNSRecord() received records = %v, expected %v", receivedRecords, tt.expectedSet)
			}
		})
	}
}

func TestQueryDNS(t *testing.T) {
	tests := []struct {
		name          string
		domain        string
		mockResponses map[uint16]*dns.Msg
		mockError     error
		expected      *DNSRecords
		expectError   bool
	}{
		{
			name:   "Valid A record query",
			domain: "example.com",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeA: {
					Answer: []dns.RR{
						&dns.A{Hdr: dns.RR_Header{Name: "example.com."}, A: net.ParseIP("10.0.0.1")},
					},
				},
			},
			expected: &DNSRecords{
				ARecords: []string{"10.0.0.1"},
			},
		},
		{
			name:   "Valid AAAA record query",
			domain: "example.com",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeAAAA: {
					Answer: []dns.RR{
						&dns.AAAA{Hdr: dns.RR_Header{Name: "example.com."}, AAAA: net.ParseIP("2001:db8::1")},
					},
				},
			},
			expected: &DNSRecords{
				AAAARecords: []string{"2001:db8::1"},
			},
		},
		{
			name:   "Valid CNAME record query",
			domain: "www.example.com",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeCNAME: {
					Answer: []dns.RR{
						&dns.CNAME{Hdr: dns.RR_Header{Name: "www.example.com."}, Target: "example.com."},
					},
				},
			},
			expected: &DNSRecords{
				CNAMERecords: []string{"example.com."},
			},
		},
		{
			name:   "Valid MX record query",
			domain: "example.com",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeMX: {
					Answer: []dns.RR{
						&dns.MX{Hdr: dns.RR_Header{Name: "example.com."}, Mx: "mail.example.com.", Preference: 10},
					},
				},
			},
			expected: &DNSRecords{
				MXRecords: []MXRecord{
					{Host: "mail.example.com.", Pref: 10},
				},
			},
		},
		{
			name:   "Valid TXT record query",
			domain: "example.com",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeTXT: {
					Answer: []dns.RR{
						&dns.TXT{Hdr: dns.RR_Header{Name: "example.com."}, Txt: []string{"v=spf1 include:_spf.example.com ~all"}},
					},
				},
			},
			expected: &DNSRecords{
				TXTRecords: []string{"v=spf1 include:_spf.example.com ~all"},
			},
		},
		{
			name:   "Valid NS record query",
			domain: "example.com",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeNS: {
					Answer: []dns.RR{
						&dns.NS{Hdr: dns.RR_Header{Name: "example.com."}, Ns: "ns1.example.com."},
					},
				},
			},
			expected: &DNSRecords{
				NSRecords: []string{"ns1.example.com."},
			},
		},
		{
			name:   "Query returns no answers",
			domain: "example.com",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeA: {},
			},
			expected: &DNSRecords{},
		},
		{
			name:        "Query returns an error",
			domain:      "example.com",
			mockError:   fmt.Errorf("network error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockIDNSClient{
				MockExchange: func(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
					if tt.mockError != nil {
						return nil, 0, tt.mockError
					}
					if resp, ok := tt.mockResponses[msg.Question[0].Qtype]; ok {
						return resp, 0, nil
					}
					return &dns.Msg{}, 0, nil
				},
			}

			records, err := QueryDNS(tt.domain, "8.8.8.8", client)
			if (err != nil) != tt.expectError {
				t.Errorf("QueryDNS() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !reflect.DeepEqual(records, tt.expected) {
				t.Errorf("QueryDNS() = %v, expected %v", records, tt.expected)
			}
		})
	}
}

func TestQueryAndExtract(t *testing.T) {
	tests := []struct {
		name          string
		testType      string
		domain        string
		dnsServer     string
		mockResponses map[uint16]*dns.Msg
		mockError     error
		expected      []string
		expectError   bool
	}{
		{
			name:      "Valid A record extraction",
			testType:  "a",
			domain:    "example.com",
			dnsServer: "8.8.8.8",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeA: {
					Answer: []dns.RR{
						&dns.A{Hdr: dns.RR_Header{Name: "example.com."}, A: net.ParseIP("10.0.0.1")},
					},
				},
			},
			expected:    []string{"10.0.0.1"},
			expectError: false,
		},
		{
			name:      "Valid MX record extraction",
			testType:  "mx",
			domain:    "example.com",
			dnsServer: "8.8.8.8",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeMX: {
					Answer: []dns.RR{
						&dns.MX{Hdr: dns.RR_Header{Name: "example.com."}, Mx: "mail.example.com.", Preference: 10},
					},
				},
			},
			expected:    []string{"mail.example.com."},
			expectError: false,
		},
		{
			name:      "Query with no answers",
			testType:  "a",
			domain:    "example.com",
			dnsServer: "8.8.8.8",
			mockResponses: map[uint16]*dns.Msg{
				dns.TypeA: {},
			},
			expected:    []string{},
			expectError: false,
		},
		{
			name:        "Query returns an error",
			testType:    "a",
			domain:      "example.com",
			dnsServer:   "8.8.8.8",
			mockError:   fmt.Errorf("network error"),
			expectError: true,
		},
		{
			name:        "Invalid query type",
			testType:    "invalid",
			domain:      "example.com",
			dnsServer:   "8.8.8.8",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockIDNSClient{
				MockExchange: func(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
					if tt.mockError != nil {
						return nil, 0, tt.mockError
					}
					if resp, ok := tt.mockResponses[msg.Question[0].Qtype]; ok {
						return resp, 0, nil
					}
					return &dns.Msg{}, 0, nil
				},
			}

			got, err := QueryAndExtract(client, tt.testType, tt.dnsServer, tt.domain)
			if (err != nil) != tt.expectError {
				t.Errorf("QueryAndExtract() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("QueryAndExtract() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestGetQueryTypeFromString(t *testing.T) {
	tests := []struct {
		name     string
		testType string
		want     uint16
		wantErr  bool
	}{
		{
			name:     "Valid A record type",
			testType: "a",
			want:     dns.TypeA,
			wantErr:  false,
		},
		{
			name:     "Valid AAAA record type",
			testType: "aaaa",
			want:     dns.TypeAAAA,
			wantErr:  false,
		},
		{
			name:     "Valid CNAME record type",
			testType: "cname",
			want:     dns.TypeCNAME,
			wantErr:  false,
		},
		{
			name:     "Valid MX record type",
			testType: "mx",
			want:     dns.TypeMX,
			wantErr:  false,
		},
		{
			name:     "Valid TXT record type",
			testType: "txt",
			want:     dns.TypeTXT,
			wantErr:  false,
		},
		{
			name:     "Valid NS record type",
			testType: "ns",
			want:     dns.TypeNS,
			wantErr:  false,
		},
		{
			name:     "Invalid record type",
			testType: "invalid",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "Empty record type",
			testType: "",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "Mixed case record type",
			testType: "Mx",
			want:     dns.TypeMX,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetQueryTypeFromString(tt.testType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueryType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetQueryType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractRecords(t *testing.T) {
	type args[T uint16 | string] struct {
		records *DNSRecords
		qtype   T
	}
	tests := []struct {
		name string
		args any
		want []string
	}{
		{
			name: "Extract A records (uint16)",
			args: args[uint16]{
				records: &DNSRecords{
					ARecords: []string{"10.0.0.1"},
				},
				qtype: dns.TypeA,
			},
			want: []string{"10.0.0.1"},
		},
		{
			name: "Extract A records (string)",
			args: args[string]{
				records: &DNSRecords{
					ARecords: []string{"10.0.0.1"},
				},
				qtype: "A",
			},
			want: []string{"10.0.0.1"},
		},
		{
			name: "Extract AAAA records (uint16)",
			args: args[uint16]{
				records: &DNSRecords{
					AAAARecords: []string{"2001:db8::1"},
				},
				qtype: dns.TypeAAAA,
			},
			want: []string{"2001:db8::1"},
		},
		{
			name: "Extract AAAA records (string)",
			args: args[string]{
				records: &DNSRecords{
					AAAARecords: []string{"2001:db8::1"},
				},
				qtype: "AAAA",
			},
			want: []string{"2001:db8::1"},
		},
		{
			name: "Extract CNAME records (uint16)",
			args: args[uint16]{
				records: &DNSRecords{
					CNAMERecords: []string{"example.com."},
				},
				qtype: dns.TypeCNAME,
			},
			want: []string{"example.com."},
		},
		{
			name: "Extract CNAME records (string)",
			args: args[string]{
				records: &DNSRecords{
					CNAMERecords: []string{"example.com."},
				},
				qtype: "CNAME",
			},
			want: []string{"example.com."},
		},
		{
			name: "Extract MX records (uint16)",
			args: args[uint16]{
				records: &DNSRecords{
					MXRecords: []MXRecord{
						{Host: "mail.example.com.", Pref: 10},
					},
				},
				qtype: dns.TypeMX,
			},
			want: []string{"mail.example.com."},
		},
		{
			name: "Extract MX records (string)",
			args: args[string]{
				records: &DNSRecords{
					MXRecords: []MXRecord{
						{Host: "mail.example.com.", Pref: 10},
					},
				},
				qtype: "MX",
			},
			want: []string{"mail.example.com."},
		},
		{
			name: "Extract TXT records (uint16)",
			args: args[uint16]{
				records: &DNSRecords{
					TXTRecords: []string{"v=spf1 include:_spf.example.com ~all"},
				},
				qtype: dns.TypeTXT,
			},
			want: []string{"v=spf1 include:_spf.example.com ~all"},
		},
		{
			name: "Extract TXT records (string)",
			args: args[string]{
				records: &DNSRecords{
					TXTRecords: []string{"v=spf1 include:_spf.example.com ~all"},
				},
				qtype: "TXT",
			},
			want: []string{"v=spf1 include:_spf.example.com ~all"},
		},
		{
			name: "Extract NS records (uint16)",
			args: args[uint16]{
				records: &DNSRecords{
					NSRecords: []string{"ns1.example.com."},
				},
				qtype: dns.TypeNS,
			},
			want: []string{"ns1.example.com."},
		},
		{
			name: "Extract NS records (string)",
			args: args[string]{
				records: &DNSRecords{
					NSRecords: []string{"ns1.example.com."},
				},
				qtype: "NS",
			},
			want: []string{"ns1.example.com."},
		},
		{
			name: "Record type not found",
			args: args[uint16]{
				records: &DNSRecords{},
				qtype:   dns.TypeA,
			},
			want: nil,
		},
		{
			name: "Invalid DNS query type (string)",
			args: args[string]{
				records: &DNSRecords{},
				qtype:   "invalidtype",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.args.(type) {
			case args[uint16]:
				got := ExtractRecords(v.records, v.qtype)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ExtractRecords() = %v, want %v", got, tt.want)
				}
			case args[string]:
				got := ExtractRecords(v.records, v.qtype)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ExtractRecords() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
