package zone

import "strings"

type Zone struct {
	SOA     ResourceRecord
	Records []ResourceRecord
}

func (z *Zone) String() string {
	str := strings.Builder{}

	if z.SOA.IsEmpty() == false {
		str.WriteString(z.SOA.String())
	}
	for _, rr := range z.Records {
		str.WriteString(rr.String())
	}
	return str.String()
}
