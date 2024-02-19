package zone

import (
	"github.com/spf13/cast"
	"strings"
)

type ResourceRecord struct {
	Name   string
	Class  string
	Type   string
	TTL    int
	Values []string
}

func (rr *ResourceRecord) String() string {
	str := strings.Builder{}
	if rr.Name != "" {
		str.WriteString(rr.Name + " ")
	}
	if rr.TTL > 0 {
		str.WriteString(cast.ToString(rr.TTL) + " ")
	}
	if rr.Class != "" {
		str.WriteString(rr.Class + " ")
	}
	if rr.Type != "" {
		str.WriteString(rr.Type + " ")
	}
	for _, v := range rr.Values {
		str.WriteString(v + " ")
	}
	return strings.TrimSpace(str.String()) + "\n"
}

func (rr *ResourceRecord) IsEmpty() bool {
	return rr.TTL == 0 &&
		rr.Class == "" &&
		rr.Name == "" &&
		rr.Type == "" &&
		len(rr.Values) == 0
}
