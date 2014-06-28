// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package jekyll

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v1"
)

const (
	// rubyDateFormat is the default string representation for ruby Time values
	rubyDateFormat = "2006-01-02 15:04:05 -0700"
)

// front matter delimiter
var delim = []byte("---\n")

// Page is a jekyll page or post.
type Page struct {
	// Name is the file name for this page.
	Name string

	// FrontMatter is the parsed YAML metadata from the top of the page.
	FrontMatter map[string]interface{}
}

// NewPage parses the Jekyll file f into a new Page.
func NewPage(f *os.File) (*Page, error) {
	p := &Page{Name: filepath.Base(f.Name())}

	err := p.parseFrontMatter(f)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// parseFrontMatter reads from r, extracts the front matter (YAML at the top of
// the file between '---\n' delimeters) and populates p.FrontMatter.  If r does
// not contain properly formatted front matter, an error is returned.
func (p *Page) parseFrontMatter(r io.Reader) error {
	buf := bufio.NewReader(r)
	peek, err := buf.Peek(len(delim))
	if err != nil {
		return err
	}

	if bytes.Equal(peek, delim) {
		buf.Read(make([]byte, len(delim))) // throw away

		var fm bytes.Buffer
		for {
			line, err := buf.ReadBytes('\n')
			if err != nil {
				// io.EOF is treated as error as well
				return err
			}
			if bytes.Equal(line, delim) {
				break
			}
			fm.Write(line)
		}

		// unmarshall yaml
		err := yaml.Unmarshal(fm.Bytes(), &p.FrontMatter)
		if err != nil {
			return err
		}
	}

	return nil
}

// Time returns the published time of p.  It first looks for a 'date' key in
// p.FrontMatter, then for a date embedded in p.Name.
func (p *Page) Time() time.Time {
	// parse date from front matter
	if d, ok := p.FrontMatter["date"]; ok {
		if date, ok := d.(string); ok {
			if t, err := time.Parse(time.RFC3339, date); err == nil {
				return t
			}
			if t, err := time.Parse(rubyDateFormat, date); err == nil {
				return t
			}
		}
	}

	// fallback to filename
	if p.Name != "" {
		if np := strings.SplitN(p.Name, "-", 4); len(np) >= 3 {
			date := strings.Join(np[0:3], "-")
			if t, err := time.Parse("2006-01-02", date); err == nil {
				return t
			}
		}
	}

	return time.Time{}
}
