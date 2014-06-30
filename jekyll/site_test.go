package jekyll

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func newTestSite(config string) (string, func(), error) {
	dir, err := ioutil.TempDir(os.TempDir(), "gum")
	if err != nil {
		return "", nil, err
	}

	// write config file
	if err := ioutil.WriteFile(path.Join(dir, configFile), []byte(config), 0644); err != nil {
		return "", nil, err
	}

	cleanup := func() { os.RemoveAll(dir) }
	return dir, cleanup, nil
}

func TestSite_PermalinkTemplate(t *testing.T) {
	tests := []struct {
		config map[string]interface{}
		want   string
	}{
		{
			config: nil,
			want:   "/:year/:month/:day/:title.html",
		},
		{
			config: map[string]interface{}{"permalink": ""},
			want:   "/:year/:month/:day/:title.html",
		},
		{
			config: map[string]interface{}{"permalink": "date"},
			want:   "/:year/:month/:day/:title.html",
		},
		{
			config: map[string]interface{}{"permalink": "pretty"},
			want:   "/:year/:month/:day/:title/",
		},
		{
			config: map[string]interface{}{"permalink": "none"},
			want:   "/:title.html",
		},
		{
			config: map[string]interface{}{"permalink": "/foo"},
			want:   "/foo",
		},
	}

	for _, tt := range tests {
		s := &Site{config: tt.config}
		if got := s.PermalinkTemplate(); got != tt.want {
			t.Errorf("PermalinkTemplate with config %q got: %v, want: %v", tt.config, got, tt.want)
		}
	}
}

func TestSite_ParseConfig(t *testing.T) {
	dir, cleanup, err := newTestSite("source: src")
	if err != nil {
		t.Fatalf("error creating test site: %v", err)
	}
	defer cleanup()

	s := &Site{base: dir}
	s.parseConfig()

	if source, ok := s.config["source"]; !ok {
		t.Fatal("s.config does not contain key: source")
	} else if got, ok := source.(string); !ok {
		t.Fatal("s.config[source] is not a string")
	} else if want := "src"; got != want {
		t.Fatalf("s.config[source] got: %v, want: %v", got, want)
	}
}

func TestSite_LoadPosts(t *testing.T) {
	dir, cleanup, err := newTestSite("source: src")
	if err != nil {
		t.Fatalf("error creating test site: %v", err)
	}
	defer cleanup()

	// write post
	os.MkdirAll(path.Join(dir, "src", "_posts"), 0755)
	postPath := path.Join(dir, "src", "_posts", "2014-05-28-test.md")
	if err := ioutil.WriteFile(postPath, []byte("---\ntitle: test\n---\n"), 0644); err != nil {
		t.Fatalf("error creating test post: %v", err)
	}

	s, err := NewSite(dir)
	if err != nil {
		t.Fatalf("NewSite(%q) returned error: %v", dir, err)
	}

	posts, err := s.loadPosts()
	if err != nil {
		t.Fatalf("Error loading posts: %v", err)
	}

	if got, want := len(posts), 1; got != want {
		t.Fatalf("len(posts): %d, want: %d", got, want)
	}

	want := time.Date(2014, 5, 28, 0, 0, 0, 0, time.UTC)
	if got := posts[0].Time(); !got.Equal(want) {
		t.Fatalf("post[0].Time(): %v, want: %v", got, want)
	}
}
