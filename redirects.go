// Package redirects provides Netlify style _redirects file format parsing.
package redirects

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// Params is a map of key/value pairs.
type Params map[string]interface{}

const format = "`from [a=:save1 b=value] to [code][!] [Country=x,y,z] [Language=x,y,z]`"

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

		// src and dst
		rule := Rule{
			From:   fields[0],
			Status: 301,
		}

		var parameters []string
		var i int
		var implyStatusCode bool

		// grab all parameters
		for i = 1; strings.ContainsAny(fields[i], "="); i++ {
			parameters = append(parameters, fields[i])
		}

		// if there were any paramters, add them to the rules
		if len(parameters) != 0 {
			rule.Params = parseParams(parameters)
		}

		// if `to` field is empty, else validate `to` path
		if field := fields[i:]; field == nil {
			return nil, fmt.Errorf("Missing `to` field %s", format)
		} else {
			if rule.To, err = isPath(field[0]); err != nil {
				return nil, err
			}
		}

		options := fields[i+1:]
		// if there is anything after `to` field
		if len(options) > 0 {
			// check for status code
			if rule.Status, err = isInteger(options[0]); err != nil {
				// not a number, could be [status code][!]
				if rule.Status, rule.Force, err = isStatusCode(options[0]); err != nil {
					// not [status][!], could be key/pair
					if !isKeyPair(options[0]) {
						return nil, fmt.Errorf("got: %s, was expecting format %s", options[0], format)
					}
					rule.Status = 301
					implyStatusCode = true
				}
			}
			if !implyStatusCode {
				options = options[1:]
			}

			// parse for country and/or language options, error out if anything else was found.
			if rule.Country, rule.Language, err = parseOptions(options); err != nil {
				return nil, err
			}
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

// isPath checks for mal-formed path, then returns it
func isPath(token string) (string, error) {
	if _, err := isInteger(token); err == nil {
		return "", fmt.Errorf("Numbers not allowed. Got: %s, was expecting format %s", token, format)
	}
	if strings.ContainsAny(token, "=") {
		return "", fmt.Errorf("`=` not allowed. Got: %s, was expecting format %s", token, format)
	}
	if strings.HasSuffix(token, "!") {
		return "", fmt.Errorf("`!` not allowed. Got: %s, was expecting format %s", token, format)
	}
	if (!strings.HasPrefix(token, "/") || !strings.HasPrefix(token, "http://") || !strings.HasPrefix(token, "https://")) == false {
		return "", fmt.Errorf("Got: %s, path must start with `/`, `http://`, or `https://`", token)
	}
	return token, nil
}

// isInteger returns results from strconv.Atoi(string)
func isInteger(token string) (int, error) {
	return strconv.Atoi(token)
}

// isStatusCode returns results from parseStatus(string)
func isStatusCode(token string) (int, bool, error) {
	return parseStatus(token)
}

// isKeyPair returns true if `=` is found in token
func isKeyPair(token string) bool {
	return strings.ContainsAny(token, "=")
}

// parseParams returns parsed param key/value pairs.
func parseParams(pairs []string) Params {
	m := make(Params)

	for _, p := range pairs {
		parts := strings.Split(p, "=")
		if len(parts) > 1 {
			m[parts[0]] = parts[1]
		} else {
			m[parts[0]] = true
		}
	}

	return m
}

// parseStatus returns the status code and force when "!" suffix is present.
func parseStatus(s string) (code int, force bool, err error) {
	if strings.HasSuffix(s, "!") {
		force = true
		s = strings.Replace(s, "!", "", -1)
	}

	code, err = strconv.Atoi(s)
	return
}

// parseOptions extracts optional Country and Language fields.
func parseOptions(options []string) (Country []string, Language []string, err error) {
	// parse for country and or language options, skips if empty
	for _, token := range options {
		// if we find something other than a key/pair past the `status code` place, error out
		if !strings.ContainsAny(token, "=") {
			err = fmt.Errorf("got: %s, was expecting format %s", token, format)
			break
		}
		k, v := splitOptions(token)
		if k == "Country" {
			Country = v
		} else if k == "Language" {
			Language = v
		} else {
			err = fmt.Errorf("got: %s, was expecting format %s", token, format)
		}
	}
	return
}

// parseParams returns parsed key/value pairs.
func splitOptions(options string) (string, []string) {
	parts := strings.Split(options, "=")
	if len(parts) > 1 {
		return parts[0], strings.Split(parts[1], ",")
	} else {
		return parts[0], []string{}
	}
}
