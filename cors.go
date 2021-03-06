// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cors provides handlers to enable CORS support.
package cors

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerMaxAge           = "Access-Control-Max-Age"

	headerOrigin         = "Origin"
	headerRequestMethod  = "Access-Control-Request-Method"
	headerRequestHeaders = "Access-Control-Request-Headers"
)

var (
	defaultAllowHeaders = []string{"Origin", "Accept", "Content-Type", "Authorization"}
	// Regex patterns are generated from AllowOrigins. These are used and generated internally.
	allowOriginPatterns = []string{}
)

// Options represents Access Control options.
type Options struct {
	// If set, all origins are allowed.
	AllowAllOrigins bool
	// A list of allowed origins. Wild cards and FQDNs are supported.
	AllowOrigins []string
	// If set, allows to share auth credentials such as cookies.
	AllowCredentials bool
	// A list of allowed HTTP methods.
	AllowMethods []string
	// A list of allowed HTTP headers.
	AllowHeaders []string
	// A list of exposed HTTP headers.
	ExposeHeaders []string
	// Max age of the CORS headers.
	MaxAge time.Duration
}

// Header converts options into CORS headers.
func (o *Options) Header(origin string) (headers map[string]string) {
	headers = make(map[string]string)
	// if origin is not allowed, don't extend the headers
	// with CORS headers.
	if !o.AllowAllOrigins && !o.IsOriginAllowed(origin) {
		return
	}

	// add allow origin, if empty add * -- this is ok b/c we already determined allowed above
	if (origin == "") {
		headers[headerAllowOrigin] = "*"
	} else {
		headers[headerAllowOrigin] = origin		
	}

	if (o.AllowCredentials) {
		// add allow credentials
		headers[headerAllowCredentials] = strconv.FormatBool(o.AllowCredentials)
	}

	// add allow methods
	if len(o.AllowMethods) > 0 {
		headers[headerAllowMethods] = strings.Join(o.AllowMethods, ",")
	}

	// add allow headers
	if len(o.AllowHeaders) > 0 {
		// TODO: Add default headers
		headers[headerAllowHeaders] = strings.Join(o.AllowHeaders, ",")
	} else {
		headers[headerAllowHeaders] = strings.Join(defaultAllowHeaders, ",")
	}

	// add exposed header
	if len(o.ExposeHeaders) > 0 {
		headers[headerExposeHeaders] = strings.Join(o.ExposeHeaders, ",")
	}
	// add a max age header
	if o.MaxAge > time.Duration(0) {
		headers[headerMaxAge] = strconv.FormatInt(int64(o.MaxAge/time.Second), 10)
	}
	return
}

// PreflightHeader converts options into CORS headers for a preflight response.
func (o *Options) PreflightHeader(origin, rMethod, rHeaders string) (headers map[string]string) {
	headers = make(map[string]string)
	if !o.AllowAllOrigins && !o.IsOriginAllowed(origin) {
		return
	}
	// verify if requested method is allowed
	// TODO: Too many for loops
	for _, method := range o.AllowMethods {
		if method == rMethod {
			headers[headerAllowMethods] = strings.Join(o.AllowMethods, ",")
			break
		}
	}

	allowHeaders := o.AllowHeaders
	if len(allowHeaders) == 0 {
		allowHeaders = defaultAllowHeaders
	}
	// verify if requested headers are allowed
	var allowed []string
	for _, rHeader := range strings.Split(rHeaders, ",") {
		rHeader = strings.TrimSpace(rHeader)
	lookupLoop:
		for _, allowedHeader := range allowHeaders {
			if strings.ToLower(rHeader) == strings.ToLower(allowedHeader) {
				allowed = append(allowed, rHeader)
				break lookupLoop
			}
		}
	}

	if (o.AllowCredentials) {
		headers[headerAllowCredentials] = strconv.FormatBool(o.AllowCredentials)
	}

	// add allow origin
	headers[headerAllowOrigin] = origin

	// add allowed headers
	if len(allowed) > 0 {
		headers[headerAllowHeaders] = strings.Join(allowed, ",")
	}

	// add exposed headers
	if len(o.ExposeHeaders) > 0 {
		headers[headerExposeHeaders] = strings.Join(o.ExposeHeaders, ",")
	}
	// add a max age header
	if o.MaxAge > time.Duration(0) {
		headers[headerMaxAge] = strconv.FormatInt(int64(o.MaxAge/time.Second), 10)
	}
	return
}

// IsOriginAllowed looks up if the origin matches one of the patterns
// generated from Options.AllowOrigins patterns.
func (o *Options) IsOriginAllowed(origin string) (allowed bool) {
	for _, pattern := range allowOriginPatterns {
		allowed, _ = regexp.MatchString(pattern, origin)
		if allowed {
			return
		}
	}
	return
}

// Allow enables CORS for requests those match the provided options.
func (o *Options) Allow(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	for _, origin := range o.AllowOrigins {
		pattern := regexp.QuoteMeta(origin)
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)
		allowOriginPatterns = append(allowOriginPatterns, "^"+pattern+"$")
	}

	var (
		origin           = req.Header.Get(headerOrigin)
		requestedMethod  = req.Header.Get(headerRequestMethod)
		requestedHeaders = req.Header.Get(headerRequestHeaders)
		// additional headers to be added
		// to the response.
		headers map[string]string
	)

	if req.Method == "OPTIONS" &&
		(requestedMethod != "" || requestedHeaders != "") {
		// TODO: if preflight, respond with exact headers if allowed
		headers = o.PreflightHeader(origin, requestedMethod, requestedHeaders)
		for key, value := range headers {
			res.Header().Set(key, value)
		}
		res.WriteHeader(http.StatusOK)
		return
	}
		
	headers = o.Header(origin)
	
	for key, value := range headers {
		
		res.Header().Set(key, value)
	}
	
	next(res, req)

}
