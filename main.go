package zone

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"io"
	"regexp"
	"strings"
)

const bracketCloseByte = byte(')')
const bracketOpenByte = byte('(')
const commentStartByte = byte(';')
const escapeByte = byte('\\')
const quoteByte = byte('"')
const defaultTtl = 86400

var commentStartBytes = []byte{commentStartByte}
var originLineBytes = []byte("$ORIGIN")
var ttlLineBytes = []byte("$TTL")
var includeLineBytes = []byte("$INCLUDE")
var generateLineBytes = []byte("$GENERATE")

// isClassToken matches the classes defined in
// https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.4
//
// It also recognizes the QCLASS "any" defined in
// https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.5
var isClassToken = regexp.MustCompile(`in|IN|ch|CH|hs|HS|cs|CS|any|ANY`)
var isTtlToken = regexp.MustCompile(`^[0-9]+$`)
var isSoaLine = regexp.MustCompile(`\s+(SOA|soa)\s+`)

// ZoneParser reads zone [master files] into [Zone] objects.
//
// [master files]: https://datatracker.ietf.org/doc/html/rfc1035#autoid-48
type ZoneParser struct {
	defaultTtl      int
	preferSoaMinTtl bool
	skipIncludes    bool
}

type Option func(zp *ZoneParser) error

func NewZoneParser(opts ...Option) (*ZoneParser, error) {
	zoneParser := &ZoneParser{
		skipIncludes: true,
		defaultTtl:   defaultTtl,
	}

	for _, opt := range opts {
		err := opt(zoneParser)
		if err != nil {
			return nil, err
		}
	}

	return zoneParser, nil
}

// WithDefaultTtl allows defining the default TTL that will be used when no
// $TTL directive has been found. The default is `86_400`.
// If `WithPreferSoaMinTtl(true)` is used, and a SOA record is present, then
// this value will be ignored.
func WithDefaultTtl(value int) Option {
	return func(zp *ZoneParser) error {
		zp.defaultTtl = value
		return nil
	}
}

// WithPreferSoaMinTtl will _always_ use the minimum TTL value from the SOA
// line when value is `true`. Any `$TTL` directives will be ignored.
func WithPreferSoaMinTtl(value bool) Option {
	return func(zp *ZoneParser) error {
		zp.preferSoaMinTtl = value
		return nil
	}
}

func WithSkipIncludes(value bool) Option {
	return func(zp *ZoneParser) error {
		if value == false {
			return fmt.Errorf("parse $INCLUDE directive: %w", ErrNotImplemented)
		}
		zp.skipIncludes = value
		return nil
	}
}

// Parse reads the given reader line-by-line as a zone file.
// All comments are discarded.
func (zp *ZoneParser) Parse(reader io.Reader) (*Zone, error) {
	result := &Zone{
		Records: make([]ResourceRecord, 0),
	}

	r := bufio.NewReader(reader)
	var currentOrigin string
	var currentTtl int
	var lastRecord ResourceRecord
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if bytes.Equal(line, []byte("\n")) ||
			bytes.HasPrefix(line, commentStartBytes) ||
			bytes.HasPrefix(line, generateLineBytes) {
			continue
		}

		if isContinuedLine(line) {
			line, err = readContinuedLine(r, line)
			if err != nil {
				return nil, err
			}
		}

		if bytes.Equal(line[0:7], originLineBytes) {
			currentOrigin = parseOriginLine(line)
			continue
		}

		if bytes.Equal(line[0:4], ttlLineBytes) {
			if zp.preferSoaMinTtl == true {
				continue
			}
			currentTtl = parseTtlLine(line)
			continue
		}

		if bytes.Equal(line[0:8], includeLineBytes) {
			if zp.skipIncludes == true {
				continue
			}
			// TODO: support handling $INCLUDE lines
		}

		if isSoaLine.Match(line) == true {
			record := parseSoaLine(line)
			if record.Name == "@" && currentOrigin != "" {
				record.Name = currentOrigin
			}
			if record.TTL == 0 {
				if zp.preferSoaMinTtl {
					minTtl := cast.ToInt(record.Values[len(record.Values)-1])
					record.TTL = minTtl
					currentTtl = minTtl
				} else if currentTtl > 0 {
					record.TTL = currentTtl
				} else {
					record.TTL = defaultTtl
				}
			}
			result.SOA = record
			lastRecord = record
			continue
		}

		record := parseRecordLine(line)
		if record.Name == "" {
			if lastRecord.Name != "" {
				record.Name = lastRecord.Name
			} else {
				record.Name = currentOrigin
			}
		}
		if strings.HasSuffix(record.Name, ".") == false {
			if currentOrigin != "" {
				record.Name = record.Name + "." + currentOrigin
			}
		}
		if record.TTL == 0 {
			if currentTtl > 0 {
				record.TTL = currentTtl
			} else {
				record.TTL = zp.defaultTtl
			}
		}

		result.Records = append(result.Records, record)
		lastRecord = record
	}

	return result, nil
}
