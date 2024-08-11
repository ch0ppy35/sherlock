package dns

import (
	"fmt"
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

type TinyDNSClient interface {
	Exchange(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error)
}

type MockTinyDNSClient struct {
	MockExchange func(*dns.Msg, string) (*dns.Msg, time.Duration, error)
}

func (m *MockTinyDNSClient) Exchange(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
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

// CompareRecords compares the expected and actual DNS records and returns an error if they don't match.
func CompareRecords(expected []string, actual []string) error {
	expectedMap := make(map[string]struct{}, len(expected))
	for _, val := range expected {
		expectedMap[val] = struct{}{}
	}

	actualMap := make(map[string]struct{}, len(actual))
	for _, val := range actual {
		actualMap[val] = struct{}{}
	}

	matchedRecords := []string{}
	unexpectedRecords := []string{}
	missingRecords := []string{}

	for _, val := range actual {
		if _, found := expectedMap[val]; !found {
			unexpectedRecords = append(unexpectedRecords, val)
		} else {
			matchedRecords = append(matchedRecords, val)
		}
	}

	for _, val := range expected {
		if _, found := actualMap[val]; !found {
			missingRecords = append(missingRecords, val)
		}
	}

	if len(matchedRecords) > 0 {
		fmt.Println("Matched records:")
		for _, record := range matchedRecords {
			fmt.Printf("  %s\n", record)
		}
	} else {
		fmt.Println("Matched records: None Found")
	}

	if len(unexpectedRecords) > 0 {
		fmt.Println("Unexpected records:")
		for _, record := range unexpectedRecords {
			fmt.Printf("  %s\n", record)
		}
	} else {
		fmt.Println("Unexpected records: None")
	}

	if len(missingRecords) > 0 {
		fmt.Println("Missing records:")
		for _, record := range missingRecords {
			fmt.Printf("  %s\n", record)
		}
	} else {
		fmt.Println("Missing records: None")
	}

	if len(unexpectedRecords) > 0 || len(missingRecords) > 0 {
		return fmt.Errorf("mismatched records found")
	}
	return nil
}

// QueryDNS fetches DNS records of various types for a given domain.
func QueryDNS(domain string, dnsServer string, client TinyDNSClient) (*DNSRecords, error) {
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

// QueryDNSRecord queries a specific DNS record type and processes the results using a setter function.
func QueryDNSRecord(client TinyDNSClient, domain, server string, qtype uint16, setter func(dns.RR)) error {
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
