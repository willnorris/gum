// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package jekyll

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

const (
	configFile = "_config.yml"
	postsDir   = "_posts"
)

// Site identifies a Jekyll site.
type Site struct {
	// base is the base path for the jekyll site.  This is the directory
	// that contains the site's _config.yml file.
	base string

	// config is the Jekyll configuration, parsed from the site's
	// _config.yml file.
	config map[string]interface{}
}

// NewSite creates a new Site at the given base path.
func NewSite(path string) (*Site, error) {
	s := &Site{base: path}
	if err := s.parseConfig(); err != nil {
		return nil, err
	}
	return s, nil
}

// parseConfig parses the site's _config.yml file and stores it in s.config.
func (s *Site) parseConfig() error {
	path := filepath.Join(s.base, configFile)
	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(byt, &s.config)
	if err != nil {
		return err
	}

	return nil
}

// loadPosts reads all of the posts for the site.
func (s *Site) loadPosts() ([]*Page, error) {
	var posts []*Page

	var source string
	if src, ok := s.config["source"]; ok {
		source, _ = src.(string)
	}
	postsPath := filepath.Join(s.base, source, postsDir)

	loadPost := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || err != nil {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		post, err := NewPage(f)
		if err != nil {
			return err
		}

		posts = append(posts, post)
		return nil
	}

	err := filepath.Walk(postsPath, loadPost)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
