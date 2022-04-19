package redirects_test

import (
	"testing"

	"github.com/fission-suite/go-redirects"
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

func TestRule_IsRequest(t *testing.T) {
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

func TestImplicit(t *testing.T) {
	var actual []redirects.Rule
	t.Run("implicit redirect", func(t *testing.T) {
		actual, _ = redirects.ParseString(`
		# Implicit 301 redirects
		/home              /
		`)
	})
	expected.From = "/home"
	expected.To = "/"
	assert.Equal(t, expected, actual[0])
}
