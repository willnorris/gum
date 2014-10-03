// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package jekyll

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v1"
	"willnorris.com/go/newbase60"
)

const (
	// rubyDateFormat is the default string representation for ruby Time values
	rubyDateFormat = "2006-01-02 15:04:05 -0700"
)

// front matter delimiter
var delim = []byte("---\n")

// Page is a jekyll page or post.
//
// TODO: certain functions such as Slug() (and maybe others?) make assumptions
// that only apply to posts, not pages.
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

// Slug returns the slug for the page, based on the filename format:
// YYYY-MM-DD-slug.ext.
func (p *Page) Slug() string {
	n := strings.TrimSuffix(p.Name, filepath.Ext(p.Name))
	return strings.SplitN(n, "-", 4)[3]
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

// ShortURLs fetches the short URLs for the page.  This is determined by the
// short_url property in the page's front matter, as well as the wordpress_id
// property.
//
// Other short URLs may exist for the page that can only be calculated with
// knowledge of other pages/posts.  For example, short URLs that identify the
// the Nth post of type T on a particular date.
func (p *Page) ShortURLs() ([]string, error) {
	var urls []string

	if s, ok := p.FrontMatter["short_url"]; ok {
		switch v := s.(type) {
		case string:
			urls = append(urls, v)
		case []interface{}:
			for _, i := range v {
				if u, ok := i.(string); ok {
					urls = append(urls, u)
				}
			}
		default:
			return nil, fmt.Errorf("unable to parse short_url: %v", s)
		}
	}

	// this is very specific to posts imported from a WordPress blog that
	// used the Hum plugin.
	if w, ok := p.FrontMatter["wordpress_id"]; ok {
		id, ok := w.(int)
		if !ok {
			return nil, fmt.Errorf("unable to parse wordpress_id: %v", w)
		}

		// newer newbase60-encoded style
		u := fmt.Sprintf("/b/%s", newbase60.EncodeInt(id))
		urls = append(urls, u)

		// really old style shortlinks
		u = fmt.Sprintf("/p/%d", id)
		urls = append(urls, u)
	}

	return urls, nil
}

// Permalink returns the permalink path for the page.
func (p *Page) Permalink(template string) string {
	if perm, ok := p.FrontMatter["permalink"]; ok {
		if s, ok := perm.(string); ok {
			return s
		}
	}

	t := p.Time()

	u := template
	u = strings.Replace(u, ":year", fmt.Sprint(t.Year()), -1)
	u = strings.Replace(u, ":short_year", fmt.Sprint(t.Year()%100), -1)
	u = strings.Replace(u, ":month", fmt.Sprintf("%02d", int(t.Month())), -1)
	u = strings.Replace(u, ":i_month", fmt.Sprint(int(t.Month())), -1)
	u = strings.Replace(u, ":day", fmt.Sprintf("%02d", t.Day()), -1)
	u = strings.Replace(u, ":i_day", fmt.Sprint(t.Day()), -1)
	u = strings.Replace(u, ":title", p.Slug(), -1)

	return u
}
