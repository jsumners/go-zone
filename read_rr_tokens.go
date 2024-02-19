package zone

import "github.com/spf13/cast"

// readRRTokens iterates a set of line tokens and adds them to a
// specified [ResourceRecord] by expected position.
func readRRTokens(tokens [][]byte, rr *ResourceRecord) {
	if isClassToken.Match(tokens[0]) {
		rr.Class = string(tokens[0])
		if isTtlToken.Match(tokens[1]) {
			// <class> <ttl> <type> <data>
			rr.TTL = cast.ToInt(string(tokens[1]))
			rr.Type = string(tokens[2])
			for _, t := range tokens[3:] {
				rr.Values = append(rr.Values, string(t))
			}
		} else {
			// <class> <type> <data>
			rr.Type = string(tokens[1])
			for _, t := range tokens[2:] {
				rr.Values = append(rr.Values, string(t))
			}
		}
	} else if isTtlToken.Match(tokens[0]) {
		rr.TTL = cast.ToInt(string(tokens[0]))
		if isClassToken.Match(tokens[1]) {
			// <ttl> <class> <type> <data>
			rr.Class = string(tokens[1])
			rr.Type = string(tokens[2])
			for _, t := range tokens[3:] {
				rr.Values = append(rr.Values, string(t))
			}
		} else {
			// <ttl> <type> <data>
			rr.Type = string(tokens[1])
			for _, t := range tokens[2:] {
				rr.Values = append(rr.Values, string(t))
			}
		}
	} else if isRecordType(tokens[0]) {
		// <type> <data>
		rr.Type = string(tokens[0])
		for _, t := range tokens[1:] {
			rr.Values = append(rr.Values, string(t))
		}
	}
}
