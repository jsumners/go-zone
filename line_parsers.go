package zone

import (
	"bytes"
	"github.com/spf13/cast"
	"strings"
)

func parseOriginLine(line []byte) string {
	data := bytes.TrimSpace(stripComment(line))
	fields := strings.Fields(string(data))
	return fields[1]
}

// parseRecordLine parses a record line into a [ResourceRecord]. It handles
// all permutations of a record line as determined by:
//
//  1. <domain-name><rr> [<comment>]
//  2. <blank><rr> [<comment>]
//
// Where <rr> has possible forms:
//  1. [<ttl>] [<class>] <type> <data>
//  2. [<class>] [<ttl>] <type> <data>
//
// If an input line is invalid, it will parse as many fields as exist in the
// line and return a record with the appropriate fields filled in.
func parseRecordLine(line []byte) ResourceRecord {
	result := ResourceRecord{}
	tokens := tokenizeLine(line)

	if isClassToken.Match(tokens[0]) || isTtlToken.Match(tokens[0]) {
		// The line looks like one of:
		// 1. "in a 1.1.1.1"
		// 2. "300 a 1.1.1.1"
		// 3. "300 in a 1.1.1.1"
		// 4. "in 300 a 1.1.1.1"
		readRRTokens(tokens, &result)
	} else if isClassToken.Match(tokens[1]) {
		// The line looks like one of:
		// 1. "a in a 1.1.1.1" (note that the first "a" is a domain)
		// 2. "a in 300 a 1.1.1.1"
		result.Name = string(tokens[0])
		readRRTokens(tokens[1:], &result)
	} else if isRecordType(tokens[0]) {
		readRRTokens(tokens, &result)
	} else {
		result.Name = string(tokens[0])
		readRRTokens(tokens[1:], &result)
	}

	return result
}

func parseSoaLine(line []byte) ResourceRecord {
	result := ResourceRecord{
		Type: "SOA",
	}
	str := strings.
		NewReplacer("(", "", ")", "").
		Replace(
			string(stripComment(line)),
		)
	fields := strings.Fields(str)

	// Maximum number of fields in a SOA record: 11.
	// Following the BNF in https://datatracker.ietf.org/doc/html/rfc1035#section-5.1:
	// <0>: domain
	// [1]: ttl or class
	// [2]: class or ttl
	// <3>: type
	// <4-10>: data
	//
	// Therefore:
	// 1. If there are 11 total fields, we have a value for every field. Fields 1
	// and 2 must be checked for text or number to determine class or ttl.
	// 2. If there are 10 total fields, either ttl or class is missing. Field 1
	// must be checked for text or number to determine class or ttl.
	// 3. If there are 9 total fields, both ttl and class are missing.

	if isClassToken.Match([]byte(fields[0])) {
		// We seem to have encountered a SOA line that is missing a leading owner
		// field. So we will force one in.
		result.Name = "@"
		fields = append([]string{"@"}, fields...)
	} else {
		// Domain name
		result.Name = fields[0]
	}

	switch len(fields) {
	case 11:
		if isTtl(fields[1]) == true {
			result.TTL = cast.ToInt(fields[1])
			result.Class = fields[2]
		} else {
			result.Class = fields[1]
			result.TTL = cast.ToInt(fields[2])
		}
		result.Values = []string{fields[4], fields[5], fields[6], fields[7], fields[8], fields[9], fields[10]}

	case 10:
		if isTtl(fields[1]) == true {
			result.TTL = cast.ToInt(fields[1])
			result.Class = "IN"
		} else {
			result.Class = fields[1]
		}
		result.Values = []string{fields[3], fields[4], fields[5], fields[6], fields[7], fields[8], fields[9]}

	case 9:
		result.Class = "IN"
		result.Values = []string{fields[2], fields[3], fields[4], fields[5], fields[6], fields[7], fields[8]}
	}

	return result
}

func parseTtlLine(line []byte) int {
	data := bytes.TrimSpace(stripComment(line))
	fields := strings.Fields(string(data))
	return cast.ToInt(fields[1])
}
