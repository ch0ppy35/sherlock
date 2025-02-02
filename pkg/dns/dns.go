package dns

import (
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type DNSRecords struct {
	ARecords     []string
	AAAARecords  []string
	CNAMERecords []string
	MXRecords    []MXRecord
	TXTRecords   []string
	NSRecords    []string
}

type MXRecord struct {
	Host string
	Pref uint16
}

// IDNSClient is a small 'github.com/miekg/dns.client' implementation for easy testing
type IDNSClient interface {
	Exchange(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error)
}

type MockIDNSClient struct {
	MockExchange func(*dns.Msg, string) (*dns.Msg, time.Duration, error)
}

func (m *MockIDNSClient) Exchange(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
	return m.MockExchange(msg, server)
}

func (r *DNSRecords) addARecord(rr dns.RR) {
	if a, ok := rr.(*dns.A); ok {
		r.ARecords = append(r.ARecords, a.A.String())
	}
}

func (r *DNSRecords) addAAAARecord(rr dns.RR) {
	if aaaa, ok := rr.(*dns.AAAA); ok {
		r.AAAARecords = append(r.AAAARecords, aaaa.AAAA.String())
	}
}

func (r *DNSRecords) addCNAMERecord(rr dns.RR) {
	if cname, ok := rr.(*dns.CNAME); ok {
		r.CNAMERecords = append(r.CNAMERecords, cname.Target)
	}
}

func (r *DNSRecords) addMXRecord(rr dns.RR) {
	if mx, ok := rr.(*dns.MX); ok {
		r.MXRecords = append(r.MXRecords, MXRecord{
			Host: mx.Mx,
			Pref: mx.Preference,
		})
	}
}

func (r *DNSRecords) addTXTRecord(rr dns.RR) {
	if txt, ok := rr.(*dns.TXT); ok {
		r.TXTRecords = append(r.TXTRecords, txt.Txt...)
	}
}

func (r *DNSRecords) addNSRecord(rr dns.RR) {
	if ns, ok := rr.(*dns.NS); ok {
		r.NSRecords = append(r.NSRecords, ns.Ns)
	}
}

// QueryDNS fetches DNS records of various types for a given domain
func QueryDNS(domain string, dnsServer string, client IDNSClient) (*DNSRecords, error) {
	records := &DNSRecords{}
	server := dnsServer + ":53"

	// Map of DNS record types to their corresponding handler functions
	queryTypes := map[uint16]func(rr dns.RR){
		dns.TypeA:     records.addARecord,
		dns.TypeAAAA:  records.addAAAARecord,
		dns.TypeCNAME: records.addCNAMERecord,
		dns.TypeMX:    records.addMXRecord,
		dns.TypeTXT:   records.addTXTRecord,
		dns.TypeNS:    records.addNSRecord,
	}

	for qtype, setter := range queryTypes {
		if err := QueryDNSRecord(client, domain, server, qtype, setter); err != nil {
			return nil, fmt.Errorf("failed to query DNS records: %w", err)
		}
	}

	return records, nil
}

// QueryDNSRecord queries a specific DNS record type and processes the results using a setter function
func QueryDNSRecord(client IDNSClient, domain string, server string, qtype uint16, setter func(dns.RR)) error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), qtype)
	resp, _, err := client.Exchange(msg, server)
	if err != nil {
		return err
	}

	for _, answer := range resp.Answer {
		setter(answer)
	}

	return nil
}

// QueryAndExtract handles the DNS query and extracts the relevant records
func QueryAndExtract(client IDNSClient, testType, dnsServer, domain string) ([]string, error) {
	qtype, err := GetQueryTypeFromString(testType)
	if err != nil {
		return nil, fmt.Errorf("invalid query type: %v", err)
	}

	records, err := QueryDNS(domain, dnsServer, client)
	if err != nil {
		return nil, fmt.Errorf("failed to query DNS: %v", err)
	}

	result, err := ExtractRecords(records, qtype)
	if err != nil {
		return nil, fmt.Errorf("failed to extract records: %v", err)
	}

	return result, nil
}

// GetQueryTypeFromString maps the string test type to the corresponding DNS query type
func GetQueryTypeFromString(testType string) (uint16, error) {
	switch strings.ToLower(testType) {
	case "a":
		return dns.TypeA, nil
	case "aaaa":
		return dns.TypeAAAA, nil
	case "cname":
		return dns.TypeCNAME, nil
	case "mx":
		return dns.TypeMX, nil
	case "txt":
		return dns.TypeTXT, nil
	case "ns":
		return dns.TypeNS, nil
	default:
		return 0, fmt.Errorf("unsupported test type, supported types: a, aaaa, cname, mx, txt, ns")
	}
}

// ExtractRecords extracts records from the DNS query results based on the query type
func ExtractRecords[T uint16 | string](records *DNSRecords, qtype T) ([]string, error) {
	var queryType uint16
	switch qt := any(qtype).(type) {
	case string:
		resolvedType, err := GetQueryTypeFromString(qt)
		if err != nil {
			return []string{}, fmt.Errorf("something went wrong determining the query type: %v", err)
		}
		queryType = resolvedType
	case uint16:
		queryType = qt
	}

	switch queryType {
	case dns.TypeA:
		return records.ARecords, nil
	case dns.TypeAAAA:
		return records.AAAARecords, nil
	case dns.TypeCNAME:
		return records.CNAMERecords, nil
	case dns.TypeMX:
		hosts := []string{}
		for _, mx := range records.MXRecords {
			hosts = append(hosts, mx.Host)
		}
		return hosts, nil
	case dns.TypeTXT:
		return records.TXTRecords, nil
	case dns.TypeNS:
		return records.NSRecords, nil
	}
	return []string{}, nil
}
