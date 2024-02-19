package zone

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"strings"
	"testing"
)

//go:embed testdata/*
var testdataFS embed.FS

func Test_WithSkipIncludes(t *testing.T) {
	fn := WithSkipIncludes(false)
	err := fn(&ZoneParser{})
	assert.ErrorIs(t, err, ErrNotImplemented)

	fn = WithSkipIncludes(true)
	err = fn(&ZoneParser{})
	assert.NoError(t, err)
}

func Test_WithDefaultTtl(t *testing.T) {
	zp, _ := NewZoneParser(WithDefaultTtl(1000))
	reader := strings.NewReader("foo in a 1.2.3.4\n")
	found, err := zp.Parse(reader)
	assert.Nil(t, err)
	expected := ResourceRecord{
		Name:   "foo",
		Class:  "in",
		Type:   "a",
		TTL:    1000,
		Values: []string{"1.2.3.4"},
	}
	assert.Equal(t, expected.String(), found.String())
}

func Test_Fixtures(t *testing.T) {
	zp, _ := NewZoneParser()
	fixtures, err := readFixtures("testdata")
	require.Nil(t, err)
	defer closeFixtures(fixtures)

	for name, fix := range fixtures {
		t.Logf("testing fixture: %s", name)
		found, err := zp.Parse(fix.input)
		assert.Nil(t, err)
		assert.Equal(t, fix.expected, found.String())
	}
}

func Test_PreferSoaMinTtl(t *testing.T) {
	zp, _ := NewZoneParser(WithPreferSoaMinTtl(true))
	fixtures, err := readFixtures("testdata/prefer_soa_min_ttl")
	require.Nil(t, err)
	defer closeFixtures(fixtures)

	for name, fix := range fixtures {
		t.Logf("testing fixture: %s", name)
		found, err := zp.Parse(fix.input)
		assert.Nil(t, err)
		assert.Equal(t, fix.expected, found.String())
	}
}

func Test_Bind9_Fixtures(t *testing.T) {
	zp, _ := NewZoneParser()
	fixtures, err := readFixtures("testdata/bind9")
	require.Nil(t, err)
	defer closeFixtures(fixtures)

	for name, fix := range fixtures {
		t.Logf("testing fixture: %s", name)
		found, err := zp.Parse(fix.input)
		assert.Nil(t, err)
		assert.Equal(t, fix.expected, found.String())
	}
}

type fixture struct {
	input    fs.File
	expected string
}

func readFixtures(dir string) (map[string]fixture, error) {
	result := make(map[string]fixture)
	files, err := testdataFS.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.HasSuffix(name, ".txt") {
			fd, err := testdataFS.Open(dir + "/" + name)
			if err != nil {
				return nil, err
			}
			result[name] = fixture{input: fd}
			continue
		}

		if strings.HasSuffix(name, ".expected") {
			fixtureName := strings.TrimSuffix(name, ".expected")
			fixture, ok := result[fixtureName]
			if ok == false {
				continue
			}
			expected, err := testdataFS.ReadFile(dir + "/" + name)
			if err != nil {
				return nil, err
			}
			fixture.expected = string(expected)
			result[fixtureName] = fixture
			continue
		}
	}

	return result, nil
}

func closeFixtures(fixtures map[string]fixture) {
	for _, fix := range fixtures {
		fix.input.Close()
	}
}
