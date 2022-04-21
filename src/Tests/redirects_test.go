package redirects_test

import (
	"fmt"
	"testing"

	redirects "github.com/fission-suite/go-redirects/src"
	"github.com/tj/assert"
)

var expected = redirects.Rule{
	From:     "",
	To:       "",
	Status:   301,
	Force:    false,
	Params:   nil,
	Country:  nil,
	Language: nil,
}

func TestParams_Has(t *testing.T) {
	p := redirects.Params{
		"foo": true,
		"bar": "baz",
	}

	assert.True(t, p.Has("foo"))
	assert.True(t, p.Has("bar"))
	assert.False(t, p.Has("baz"))
}

func TestParams_Get(t *testing.T) {
	p := redirects.Params{
		"foo": true,
		"bar": "baz",
	}

	assert.Equal(t, true, p.Get("foo"))
	assert.Equal(t, "baz", p.Get("bar"))
	assert.Equal(t, nil, p.Get("baz"))
}

func TestRule_IsProxy(t *testing.T) {
	t.Run("without host", func(t *testing.T) {
		r := redirects.Rule{
			From: "/blog",
			To:   "/blog/engineering",
		}

		assert.False(t, r.IsProxy())
	})

	t.Run("with host", func(t *testing.T) {
		r := redirects.Rule{
			From: "/blog",
			To:   "https://blog.apex.sh",
		}

		assert.True(t, r.IsProxy())
	})
}

func TestRule_IsRewrite(t *testing.T) {
	t.Run("with 3xx", func(t *testing.T) {
		r := redirects.Rule{
			From:   "/blog",
			To:     "/blog/engineering",
			Status: 302,
		}

		assert.False(t, r.IsRewrite())
	})

	t.Run("with 200", func(t *testing.T) {
		r := redirects.Rule{
			From:   "/blog",
			To:     "/blog/engineering",
			Status: 200,
		}

		assert.True(t, r.IsRewrite())
	})

	t.Run("with 0", func(t *testing.T) {
		r := redirects.Rule{
			From: "/blog",
			To:   "/blog/engineering",
		}

		assert.False(t, r.IsRewrite())
	})
}

func TestImplicit(t *testing.T) {
	var actual []redirects.Rule
	var err error
	t.Run("implicit redirect", func(t *testing.T) {
		actual, err = redirects.ParseString(`
		# Implicit 301 redirects
		/home              /
		/blog/my-post.php  /blog/my-post
	    /news              /blog
		`)
		fmt.Println(err)
	})
	expected.From = "/home"
	expected.To = "/"
	assert.Equal(t, expected, actual[0])
}
func TestForce(t *testing.T) {
	var actual []redirects.Rule
	var err error
	t.Run("with force option", func(t *testing.T) {
		actual, err = redirects.ParseString(`
		# Forcing
		/app/*   /app/index.html   200!
		`)
		fmt.Println(err)
	})
	expected.From = "/app/*"
	expected.To = "/app/index.html"
	expected.Status = 200
	expected.Force = true
	assert.Equal(t, expected, actual[0])
}

func TestParams(t *testing.T) {
	var actual []redirects.Rule
	var err error
	t.Run("with parameters", func(t *testing.T) {
		actual, err = redirects.ParseString(`
		# Parameters
		/articles id=:id tag=:tag /posts/:tag/:id
		`)
		fmt.Println(err)
	})
	params := make(redirects.Params)
	params["id"] = ":id"
	params["tag"] = ":tag"

	expected.From = "/articles"
	expected.To = "/posts/:tag/:id"
	expected.Status = 301
	expected.Force = false
	expected.Params = params
	assert.Equal(t, expected, actual[0])
}

func TestOptions(t *testing.T) {
	var actual []redirects.Rule
	var err error
	t.Run("with options", func(t *testing.T) {
		actual, err = redirects.ParseString(`
		# Country&Language
		/israel/* splat=:splat /israel/he/:splat 302! Country=au,nz Language=he
		`)
		fmt.Println(err)
	})
	params := make(redirects.Params)
	params["splat"] = ":splat"

	expected.From = "/israel/*"
	expected.To = "/israel/he/:splat"
	expected.Status = 302
	expected.Force = true
	expected.Params = params
	expected.Country = []string{"au", "nz"}
	expected.Language = []string{"he"}
	assert.Equal(t, expected, actual[0])
}
