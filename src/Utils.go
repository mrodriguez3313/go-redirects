package redirects

import (
	"fmt"
	"strconv"
	"strings"
)

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
func parseOptions(list *Token) (Country []string, Language []string, err error) {
	// parse for country and or language options, skips if empty
	for list != nil {
		// if we find something other than a key/pair past the `status code` place, error out
		if !isKeyPair(list.token) {
			err = fmt.Errorf("got: %s, was expecting format %s", list.token, format)
			break
		}
		k, v := splitOptions(list.token)
		if k == "Country" {
			Country = v
		} else if k == "Language" {
			Language = v
		} else {
			err = fmt.Errorf("got: %s, was expecting format %s", list.token, format)
		}
		list = list.next
	}
	return
}

// SplitOptions returns parsed key/value pairs.
func splitOptions(options string) (string, []string) {
	parts := strings.Split(options, "=")
	if len(parts) > 1 {
		return parts[0], strings.Split(parts[1], ",")
	} else {
		return parts[0], []string{}
	}
}
