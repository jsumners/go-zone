package zone

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseOriginLine(t *testing.T) {
	tests := []string{
		"$ORIGIN example.com.",
		"$ORIGIN example.com. ; with a comment",
	}

	for _, test := range tests {
		found := parseOriginLine([]byte(test))
		assert.Equal(t, "example.com.", found)
	}
}

func Test_parseRecordLine(t *testing.T) {
	line := []byte("IN 300 TXT \"a=b \" \"c=d\"")
	found := parseRecordLine(line)
	expected := ResourceRecord{
		Class:  "IN",
		Type:   "TXT",
		TTL:    300,
		Values: []string{"\"a=b \"", "\"c=d\""},
	}
	assert.Equal(t, expected, found)

	line = []byte("\t\t\tIN\t\tNS\t\tdns1.example.com.")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Class:  "IN",
		Type:   "NS",
		Values: []string{"dns1.example.com."},
	}
	assert.Equal(t, expected, found)

	line = []byte("\t\t\tIN\t\tNS\t\tdns1.example.com. ; comment")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Class:  "IN",
		Type:   "NS",
		Values: []string{"dns1.example.com."},
	}
	assert.Equal(t, expected, found)

	line = []byte("300 IN NS dns1")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Class:  "IN",
		Type:   "NS",
		TTL:    300,
		Values: []string{"dns1"},
	}
	assert.Equal(t, expected, found)

	line = []byte("300 NS dns1")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Type:   "NS",
		TTL:    300,
		Values: []string{"dns1"},
	}
	assert.Equal(t, expected, found)

	line = []byte("dns1 in a 1.1.1.1")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Name:   "dns1",
		Class:  "in",
		Type:   "a",
		Values: []string{"1.1.1.1"},
	}
	assert.Equal(t, expected, found)

	line = []byte("dns1 300 in a 1.1.1.1")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Name:   "dns1",
		Class:  "in",
		Type:   "a",
		TTL:    300,
		Values: []string{"1.1.1.1"},
	}
	assert.Equal(t, expected, found)

	line = []byte("dns1 in 300 a 1.1.1.1")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Name:   "dns1",
		Class:  "in",
		Type:   "a",
		TTL:    300,
		Values: []string{"1.1.1.1"},
	}
	assert.Equal(t, expected, found)

	line = []byte("NS dns.example.com")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Type:   "NS",
		Values: []string{"dns.example.com"},
	}
	assert.Equal(t, expected, found)

	// Bad line from Bind9 master2.data test file:
	line = []byte("a\t\tin\tns")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Name:  "a",
		Class: "in",
		Type:  "ns",
	}
	assert.Equal(t, expected, found)

	// Missing leading owner:
	line = []byte("\tin\tns\tns.example.com")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Name:   "",
		Class:  "in",
		Type:   "ns",
		Values: []string{"ns.example.com"},
	}
	assert.Equal(t, expected, found)

	// Bad class field:
	line = []byte("a\t\tany\tns\tns.vix.com.")
	found = parseRecordLine(line)
	expected = ResourceRecord{
		Name:   "a",
		Class:  "any",
		Type:   "ns",
		Values: []string{"ns.vix.com."},
	}
	assert.Equal(t, expected, found)
}

func Test_parseSoaLine(t *testing.T) {
	// All fields, with a leading space:
	line := []byte(" @ 300 IN SOA ns.example.com. foo.example.net. 123456 1000 1000 84000 3600")
	record := parseSoaLine(line)
	assert.Equal(t, record, ResourceRecord{
		Name:  "@",
		Class: "IN",
		Type:  "SOA",
		TTL:   300,
		Values: []string{
			"ns.example.com.",
			"foo.example.net.",
			"123456",
			"1000",
			"1000",
			"84000",
			"3600",
		},
	})

	// All fields, swapped class and ttl:
	line = []byte("@ IN 300 SOA ns.example.com. foo.example.net. 123456 1000 1000 84000 3600")
	record = parseSoaLine(line)
	assert.Equal(t, record, ResourceRecord{
		Name:  "@",
		Class: "IN",
		Type:  "SOA",
		TTL:   300,
		Values: []string{
			"ns.example.com.",
			"foo.example.net.",
			"123456",
			"1000",
			"1000",
			"84000",
			"3600",
		},
	})

	// All fields, parentheses added:
	line = []byte(" @ 300 IN SOA ns.example.com. foo.example.net. ( 123456 1000 1000 84000 3600 )")
	record = parseSoaLine(line)
	assert.Equal(t, record, ResourceRecord{
		Name:  "@",
		Class: "IN",
		Type:  "SOA",
		TTL:   300,
		Values: []string{
			"ns.example.com.",
			"foo.example.net.",
			"123456",
			"1000",
			"1000",
			"84000",
			"3600",
		},
	})

	// All fields, parentheses added with comment:
	line = []byte(" @ 300 IN SOA ns.example.com. foo.example.net. ( 123456 1000 1000 84000 3600 ) ; comment")
	record = parseSoaLine(line)
	assert.Equal(t, record, ResourceRecord{
		Name:  "@",
		Class: "IN",
		Type:  "SOA",
		TTL:   300,
		Values: []string{
			"ns.example.com.",
			"foo.example.net.",
			"123456",
			"1000",
			"1000",
			"84000",
			"3600",
		},
	})

	// Missing ttl:
	line = []byte("@ IN SOA ns.example.com. foo.example.net. 123456 1000 1000 84000 3600")
	record = parseSoaLine(line)
	assert.Equal(t, record, ResourceRecord{
		Name:  "@",
		Class: "IN",
		Type:  "SOA",
		TTL:   0,
		Values: []string{
			"ns.example.com.",
			"foo.example.net.",
			"123456",
			"1000",
			"1000",
			"84000",
			"3600",
		},
	})

	// Missing class:
	line = []byte("@ 300 SOA ns.example.com. foo.example.net. 123456 1000 1000 84000 3600")
	record = parseSoaLine(line)
	assert.Equal(t, record, ResourceRecord{
		Name:  "@",
		Class: "IN",
		Type:  "SOA",
		TTL:   300,
		Values: []string{
			"ns.example.com.",
			"foo.example.net.",
			"123456",
			"1000",
			"1000",
			"84000",
			"3600",
		},
	})

	// Missing class and ttl:
	line = []byte("@ SOA ns.example.com. foo.example.net. 123456 1000 1000 84000 3600")
	record = parseSoaLine(line)
	assert.Equal(t, record, ResourceRecord{
		Name:  "@",
		Class: "IN",
		Type:  "SOA",
		TTL:   0,
		Values: []string{
			"ns.example.com.",
			"foo.example.net.",
			"123456",
			"1000",
			"1000",
			"84000",
			"3600",
		},
	})
}

func Test_parseTtlLine(t *testing.T) {
	line := []byte("$TTL 300 ; 5 minutes")
	ttl := parseTtlLine(line)
	assert.Equal(t, 300, ttl)

	line = []byte("$TTL 300")
	ttl = parseTtlLine(line)
	assert.Equal(t, 300, ttl)
}
