/*
 * Copyright 2020 Netflix, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/netflix/weep/metadata"
	log "github.com/sirupsen/logrus"
)

func MetaDataServiceMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("ETag", strconv.FormatInt(rand.Int63n(10000000000), 10))
		w.Header().Set("Last-Modified", metadata.LastRenewal.UTC().Format("2006-01-02T15:04:05Z"))
		w.Header().Set("Server", "EC2ws")
		w.Header().Set("Content-Type", "text/plain")

		ua := r.Header.Get("User-Agent")
		metadataVersion := 1
		tokenTtl := r.Header.Get("X-Aws-Ec2-Metadata-Token-Ttl-Seconds")
		token := r.Header.Get("X-aws-ec2-metadata-token")
		// If either of these request headers exist, we can be reasonably confident that the request is for IMDSv2.
		// `X-Aws-Ec2-Metadata-Token-Ttl-Seconds` is used when requesting a token
		// `X-aws-ec2-metadata-token` is used to pass the token to the metadata service
		// Weep uses a static token, and does not perform any token validation.
		if token != "" || tokenTtl != "" {
			metadataVersion = 2
		}

		log.WithFields(log.Fields{
			"user-agent":       ua,
			"path":             r.URL.Path,
			"metadata_version": metadataVersion,
		}).Info()
		next.ServeHTTP(w, r)
	}
}
