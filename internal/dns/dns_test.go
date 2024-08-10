package dns

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/miekg/dns"
)

type MockTinyDNSClient struct {
	mockExchange func(*dns.Msg, string) (*dns.Msg, time.Duration, error)
}

func (m *MockTinyDNSClient) Exchange(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
	return m.mockExchange(msg, server)
}

func TestCompareRecords(t *testing.T) {
	type args struct {
		expected []string
		actual   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "All records match exactly",
			args: args{
				expected: []string{"example.com", "test.com"},
				actual:   []string{"example.com", "test.com"},
			},
			wantErr: false,
		},
		{
			name: "Some records are missing",
			args: args{
				expected: []string{"example.com", "test.com", "foo.com"},
				actual:   []string{"example.com", "test.com"},
			},
			wantErr: true,
		},
		{
			name: "Some unexpected records are present",
			args: args{
				expected: []string{"example.com", "test.com"},
				actual:   []string{"example.com", "test.com", "unexpected.com"},
			},
			wantErr: true,
		},
		{
			name: "Both missing and unexpected records are present",
			args: args{
				expected: []string{"example.com", "test.com", "foo.com"},
				actual:   []string{"example.com", "unexpected.com"},
			},
			wantErr: true,
		},
		{
			name: "Empty input cases",
			args: args{
				expected: []string{},
				actual:   []string{},
			},
			wantErr: false,
		},
		{
			name: "Empty expected records",
			args: args{
				expected: []string{},
				actual:   []string{"example.com"},
			},
			wantErr: true,
		},
		{
			name: "Empty actual records",
			args: args{
				expected: []string{"example.com"},
				actual:   []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CompareRecords(tt.args.expected, tt.args.actual); (err != nil) != tt.wantErr {
				t.Errorf("CompareRecords() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
			client := &MockTinyDNSClient{
				mockExchange: func(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
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
			client := &MockTinyDNSClient{
				mockExchange: func(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
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
