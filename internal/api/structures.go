package api

import (
	"errors"
	"fmt"
)

type DeployRequest struct {
	Slug      string `json:"slug"`
	Owner     string `json:"owner"`
	Domain    string `json:"domain"`
	DBUri     string `json:"db_uri"`
	SecretKey string `json:"secret"`
}

func (d DeployRequest) validate() error {
	if d.Slug == "" {
		return errors.New("slug is required")
	}
	if d.Domain == "" {
		return errors.New("domain is required")
	}
	if d.Owner == "" {
		return errors.New("owner is required")
	}
	if d.DBUri == "" {
		return errors.New("db_uri is required")
	}
	if d.SecretKey == "" {
		return errors.New("secret is required")
	}
	return nil
}

func (d DeployRequest) ToHelmValues() map[string]interface{} {
	return map[string]interface{}{
		"nameOverride":     "",
		"fullnameOverride": "",
		"domain":           d.Domain,
		"ingress": map[string]interface{}{
			"enabled": true,
			"tls": map[string]interface{}{
				"enabled":    true,
				"secretName": fmt.Sprintf("tls-%s", d.Slug),
			},
		},
		"env": []map[string]interface{}{
			{"name": "OWNER", "value": d.Owner},
			{"name": "DEBUG", "value": "False"},
			{"name": "SECRET_KEY", "value": d.SecretKey},
			{"name": "DATABASE_URI", "value": d.DBUri},
		},
	}
}
