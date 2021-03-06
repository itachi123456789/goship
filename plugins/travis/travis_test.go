package travis

import (
	"errors"
	"html/template"
	"reflect"
	"testing"

	"github.com/coreos/go-etcd/etcd"
	"github.com/gengo/goship/lib/config"
	"github.com/gengo/goship/plugins/plugin"
)

type tokenMockClient struct {
	Token string
}

func (*tokenMockClient) Set(s, c string, x uint64) (*etcd.Response, error) {
	return nil, nil
}

//Mock calls to ETCD here. Each etcd Response should return the structs you need.
func (*tokenMockClient) Get(s string, t bool, x bool) (*etcd.Response, error) {
	m := make(map[string]*etcd.Response)
	m["/projects/test_private/travis_token"] = &etcd.Response{
		Action: "Get",
		Node: &etcd.Node{
			Key: "/projects/test_private/travis_token", Value: "test_token",
		},
	}
	mockResponse, ok := m[s]
	if !ok {
		return nil, errors.New("Key doesn't exist!")
	}
	return mockResponse, nil
}

func TestRenderHeader(t *testing.T) {
	c := TravisColumn{
		Project:      "test_public",
		Token:        "",
		Organization: "test",
	}
	got, err := c.RenderHeader()
	if err != nil {
		t.Errorf(err.Error())
	}
	want := template.HTML(`<th style="min-width: 100px">Build Status</th>`)
	if want != got {
		t.Errorf("Want %#v, got %#v", want, got)
	}
}

func TestRenderDetailPublic(t *testing.T) {
	c := TravisColumn{
		Project:      "test_public",
		Token:        "",
		Organization: "test",
	}
	got, err := c.RenderDetail()
	if err != nil {
		t.Errorf(err.Error())
	}
	want := template.HTML(`<td><a target=_blank href=https://travis-ci.org/test/test_public><img src=https://travis-ci.org/test/test_public.svg?branch=master onerror='this.style.display = "none"'></img></a></td>`)
	if want != got {
		t.Errorf("Want %#v, got %#v", want, got)
	}
}

func TestRenderDetailPrivate(t *testing.T) {
	c := TravisColumn{
		Project:      "test_private",
		Token:        "test_token",
		Organization: "test",
	}
	got, err := c.RenderDetail()
	if err != nil {
		t.Errorf(err.Error())
	}
	want := template.HTML(`<td><a target=_blank href=https://magnum.travis-ci.com/test/test_private><img src=https://magnum.travis-ci.com/test/test_private.svg?token=test_token&branch=master onerror='this.style.display = "none"'></img></a></td>`)
	if want != got {
		t.Errorf("Want %#v, got %#v", want, got)
	}
}

func TestApply(t *testing.T) {
	p := &TravisPlugin{}
	proj := config.Project{
		Repo: config.Repo{
			RepoName:  "test_project",
			RepoOwner: "test",
		},
		TravisToken: "XXXXXX",
	}
	cols, err := p.Apply(proj)
	if err != nil {
		t.Fatalf("Error applying plugin %v", err)
	}
	want := []plugin.Column{
		TravisColumn{
			Organization: "test",
			Project:      "test_project",
			Token:        "XXXXXX",
		},
	}
	if got := cols; !reflect.DeepEqual(got, want) {
		t.Errorf("cols = %#v; want %#v", got, want)
	}
}
