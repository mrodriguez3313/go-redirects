// Package redirects provides Netlify style _redirects file format parsing.
package redirects

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// Params is a map of key/value pairs.
type Params map[string]interface{}

const format = "from [a=:save1 b=value] to [code][!]"

// Has returns true if the param is present.
func (p *Params) Has(key string) bool {
	if p == nil {
		return false
	}

	_, ok := (*p)[key]
	return ok
}

// Get returns the key value.
func (p *Params) Get(key string) interface{} {
	if p == nil {
		return nil
	}

	return (*p)[key]
}

// A Rule represents a single redirection or rewrite rule.
type Rule struct {
	// From is the path which is matched to perform the rule.
	From string

	// To is the destination which may be relative, or absolute
	// in order to proxy the request to another URL.
	To string

	// Status is one of the following:
	//
	// - 3xx a redirect
	// - 200 a rewrite
	// - defaults to 301 redirect
	//
	// When proxying this field is ignored.
	//
	Status int

	// Force is used to force a rewrite or redirect even
	// when a response (or static file) is present.
	Force bool

	// Params is an optional arbitrary map of key/value pairs.
	Params Params

	// Country is an optional arbitrary list of redirect options based on country ISO 3166-1 alpha-2 code
	// source: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Officially_assigned_code_elements
	Country []string

	// Language is an optional arbitrary list of redirect options based on lanugage ISO 639-1 codes
	// source: https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
	Language []string
}

// IsRewrite returns true if the rule represents a rewrite (status 200).
func (r *Rule) IsRewrite() bool {
	return r.Status == 200
}

// IsProxy returns true if it's a proxy rule (aka contains a hostname).
func (r *Rule) IsProxy() bool {
	u, err := url.Parse(r.To)
	if err != nil {
		return false
	}

	return u.Host != ""
}

// Must parse utility.
func Must(v []Rule, err error) []Rule {
	if err != nil {
		panic(err)
	}

	return v
}

// Parse the given reader.
func Parse(r io.Reader) (rules []Rule, err error) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		// empty
		if line == "" {
			continue
		}

		// comment
		if strings.HasPrefix(line, "#") {
			continue
		}

		// fields
		fields := strings.Fields(line)

		// missing dst
		if len(fields) <= 1 {
			return nil, fmt.Errorf("missing destination path: %q", line)
		}

		list := List{}
		for index, token := range fields {
			list.Insert(token, index)
		}
		list.Reverse()
		// PrintListFrom(link.head)

		// src and dst
		rule := Rule{
			Status: 301,
		}
		var parameters []string
		node := list.head

		// Assign `from` path
		if node != nil {
			if rule.From, err = isPath(node.token); err != nil {
				return nil, err
			}
			node = node.next
		} else {
			return nil, fmt.Errorf("Missing `from` field %s", format)
		}

		// assign parameters
		for node != nil {
			if !strings.ContainsAny(node.token, "=") {
				break
			}
			parameters = append(parameters, node.token)
			node = node.next
		}
		// if there were any paramters, add them to the rules
		if len(parameters) != 0 {
			rule.Params = parseParams(parameters)
		}

		// Assign `to` path
		if node != nil {
			if rule.To, err = isPath(node.token); err != nil {
				return nil, err
			}
			node = node.next
		} else {
			return nil, fmt.Errorf("Missing `to` field %s", format)
		}

		// assign status code
		if node != nil {
			if rule.Status, err = isInteger(node.token); err != nil {
				// not a number, could be [status code][!]
				if rule.Status, rule.Force, err = isStatusCode(node.token); err != nil {
					// not [status][!], could be key/pair
					// parse for country and/or language options, error out if anything else was found.
					rule.Status = 301
					if rule.Country, rule.Language, err = parseOptions(node); err != nil {
						return nil, err
					}
					rules = append(rules, rule)
					err = s.Err()
					return
				}
			}
			node = node.next
		}

		// assign country and language
		if rule.Country, rule.Language, err = parseOptions(node); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	err = s.Err()
	return
}

// ParseString parses the given string.
func ParseString(s string) ([]Rule, error) {
	return Parse(strings.NewReader(s))
}
