# go-zone

This library provides methods for parsing DNS zone (master) files as described
in [RFC 1035 §5.1][1035§5.1]. A more comprehensive parser is available in
[https://pkg.go.dev/github.com/miekg/dns#ZoneParser][miekg].
The parser in this library will parse files that contain "loose" records
whereas the `miekg` parser strictly requires an origin to be defined.

This library was written to support [gdns][gdns], a REST API client for the
[Gandi LiveDNS][livedns] service, which needs to read individual records from
zone-like files so that they can be used to provide values to the remote API.

[1035§5.1]: https://datatracker.ietf.org/doc/html/rfc1035#section-5.1
[miekg]: https://pkg.go.dev/github.com/miekg/dns#ZoneParser
[gdns]: https://github.com/jsumners/gdns
[livedns]: https://api.gandi.net/docs/livedns/

## Example

```go
package main

import (
	"fmt"
	"strings"
	"github.com/jsumners/go-zone"
)

func main() {
	zoneData := "foo 300 in a 1.2.3.4\n"
	zp, _ := zone.NewZoneParser()
	z, _ := zp.Parse(strings.NewReader(zoneData))

	fmt.Println(z)
}
```

## Note On Looseness

Consider the record line:

```
a in ns
```

The line is meant to define a nameserver record for the server `a`. But it is
missing the value. Whereas a strict parser will refuse to parse this line,
this library will return a `ResourceRecord` with an empty values list:

```go
ResourceRecord {
	Name: "a",
	Class: "in",
	Type: "ns",
	Values: []string{},
}
```

For a more complete understanding of the consequences of the looseness of the
parser, review the [testdata/bind9](./testdata/bind9) fixtures and their
expected results. The expectations do not always conform to what [Bind][bind]
would allow. There are further details in the included
[Readme](testdata/bind9/Readme.md).

[bind]: https://github.com/isc-projects/bind9

## RFCs

+ https://datatracker.ietf.org/doc/html/rfc1034
+ https://datatracker.ietf.org/doc/html/rfc1035
+ https://datatracker.ietf.org/doc/html/rfc2308
+ https://datatracker.ietf.org/doc/html/rfc4034
