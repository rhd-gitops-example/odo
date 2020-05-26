// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
)

func TestAuthorization(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/user").
		MatchHeader("Authorization", "Bearer mF_9.B5f-4.1JqM").
		Reply(200)

	client := &http.Client{
		Transport: &Authorization{
			Scheme:      "Bearer",
			Credentials: "mF_9.B5f-4.1JqM",
		},
	}

	res, err := client.Get("https://api.github.com/user")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
}

func TestAuthorization_DontOverwriteHeader(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/user").
		MatchHeader("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==").
		Reply(200)

	client := &http.Client{
		Transport: &Authorization{
			Scheme:      "Bearer",
			Credentials: "mF_9.B5f-4.1JqM",
		},
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
}
