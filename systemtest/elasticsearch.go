// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package systemtest

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/elastic/apm-tools/pkg/espoll"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"

	"github.com/elastic/apm-server/systemtest/apmservertest"
)

const (
	adminElasticsearchUser  = "admin"
	adminElasticsearchPass  = "changeme"
	maxElasticsearchBackoff = 10 * time.Second
)

var (
	// Elasticsearch is an Elasticsearch client for use in tests.
	Elasticsearch *espoll.Client
)

func initElasticSearch() {
	cfg := newElasticsearchConfig()
	cfg.Username = adminElasticsearchUser
	cfg.Password = adminElasticsearchPass
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	Elasticsearch = &espoll.Client{Client: client}
}

func newElasticsearchConfig() elasticsearch.Config {
	var addresses []string
	for _, host := range apmservertest.DefaultConfig().Output.Elasticsearch.Hosts {
		u := url.URL{Scheme: "http", Host: host}
		addresses = append(addresses, u.String())
	}
	return elasticsearch.Config{
		Addresses:  addresses,
		MaxRetries: 5,
		RetryBackoff: func(attempt int) time.Duration {
			backoff := (500 * time.Millisecond) * (1 << (attempt - 1))
			if backoff > maxElasticsearchBackoff {
				backoff = maxElasticsearchBackoff
			}
			return backoff
		},
	}
}

// CleanupElasticsearch deletes all data streams created by APM Server.
func CleanupElasticsearch(t testing.TB) {
	err := cleanupElasticsearch()
	require.NoError(t, err)
}

func cleanupElasticsearch() error {
	_, err := Elasticsearch.Do(context.Background(), &esapi.IndicesDeleteDataStreamRequest{
		Name: []string{
			"traces-apm*",
			"metrics-apm*",
			"logs-apm*",
		},
		ExpandWildcards: "all",
	}, nil)
	return err
}

// ChangeUserPassword changes the password for a given user.
func ChangeUserPassword(t testing.TB, username, password string) {
	req := esapi.SecurityChangePasswordRequest{
		Username: username,
		Body:     esutil.NewJSONReader(map[string]interface{}{"password": password}),
	}
	if _, err := Elasticsearch.Do(context.Background(), req, nil); err != nil {
		t.Fatal(err)
	}
}

func CreateAPIKey(t testing.TB, name string, privileges []string) string {
	req := esapi.SecurityCreateAPIKeyRequest{
		Body: esutil.NewJSONReader(map[string]any{
			"name": name,
			"role_descriptors": map[string]any{
				"apm": map[string]any{
					"applications": []map[string]any{
						{
							"application": "apm",
							"privileges":  privileges,
							"resources":   []string{"*"},
						},
					},
				},
			},
			"metadata": map[string]any{"application": "apm"},
		}),
	}

	res, err := Elasticsearch.Do(context.Background(), req, nil)
	if err != nil {
		t.Fatal(err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	m := make(map[string]any)
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}

	return m["encoded"].(string)
}

// InvalidateAPIKeys invalidates all API Keys for the apm-server user.
func InvalidateAPIKeys(t testing.TB) {
	req := esapi.SecurityInvalidateAPIKeyRequest{
		Body: esutil.NewJSONReader(map[string]interface{}{
			"username": apmservertest.DefaultConfig().Output.Elasticsearch.Username,
		}),
	}
	if _, err := Elasticsearch.Do(context.Background(), req, nil); err != nil {
		t.Fatal(err)
	}
}

// InvalidateAPIKeyByName invalidates the API Key with the given name.
func InvalidateAPIKeyByName(t testing.TB, name string) {
	req := esapi.SecurityInvalidateAPIKeyRequest{
		Body: esutil.NewJSONReader(map[string]interface{}{"name": name}),
	}
	if _, err := Elasticsearch.Do(context.Background(), req, nil); err != nil {
		t.Fatal(err)
	}
}
