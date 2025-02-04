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

package evaluator

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/cloudbeat/config"
	csppolicies "github.com/elastic/csp-security-policies/bundle"
	"github.com/elastic/elastic-agent-libs/logp"
)

var (
	address = "127.0.0.1:18080"

	ServerAddress = fmt.Sprintf("http://%s", address)
)

func StartServer(ctx context.Context, cfg config.Config) (*http.Server, error) {
	log := logp.NewLogger("cloudbeat_bundle_server")

	bundle := loadBundle(cfg, log)
	h := csppolicies.NewServer()
	if err := csppolicies.HostBundle("bundle.tar.gz", bundle, ctx); err != nil {
		return nil, fmt.Errorf("failed to host bundle: %w", err)
	}

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Errorf("Bundle server closed: %v", err)
		}
	}()

	return srv, nil
}

// loadBundle loads the relevant bundle based on the config.
func loadBundle(cfg config.Config, log *logp.Logger) csppolicies.Bundle {
	var bundle csppolicies.Bundle

	if cfg.Type == config.InputTypeEks {
		bundle = csppolicies.CISEksBundle()
		log.Info("Loaded CIS EKS bundle successfully")
	} else {
		bundle = csppolicies.CISKubernetesBundle()
		log.Info("Loaded CIS Vanilla bundle successfully")
	}

	return bundle
}
