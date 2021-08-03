package validator

import "net/url"

func IsValidURL(urlString string) bool {
	u, err := url.Parse(urlString)
	return err == nil && u.Host != "" && u.Scheme != ""
}
