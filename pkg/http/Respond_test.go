// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

func TestRespondAcceptCSV(t *testing.T) {
	reg := NewDefaultRegistry()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/foo/bar", nil)
	r.Header.Set("Accept", "text/csv")
	data := map[string]interface{}{"foo": "bar"}
	status := http.StatusOK
	err := Respond(w, r, reg, data, status, "")
	if !assert.NoError(t, err) {
		return
	}
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "foo\nbar\n", string(body))
}

func TestRespondAcceptJSON(t *testing.T) {
	reg := NewDefaultRegistry()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/foo/bar", nil)
	r.Header.Set("Accept", "application/json")
	data := map[string]interface{}{"foo": "bar"}
	status := http.StatusOK
	err := Respond(w, r, reg, data, status, "")
	if !assert.NoError(t, err) {
		return
	}
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "{\"foo\":\"bar\"}", string(body))
}

func TestRespondAcceptYAML(t *testing.T) {
	reg := NewDefaultRegistry()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/foo/bar", nil)
	r.Header.Set("Accept", "application/json;q=0.8, text/yaml;q=0.9")
	data := map[string]interface{}{"foo": "bar"}
	status := http.StatusOK
	err := Respond(w, r, reg, data, status, "")
	if !assert.NoError(t, err) {
		return
	}
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "foo: bar\n", string(body))
}

func TestRespondExtensionCSV(t *testing.T) {
	reg := NewDefaultRegistry()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/foo/bar.csv", nil)
	data := map[string]interface{}{"foo": "bar"}
	status := http.StatusOK
	err := Respond(w, r, reg, data, status, "")
	if !assert.NoError(t, err) {
		return
	}
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "foo\nbar\n", string(body))
}

func TestRespondExtensionJSON(t *testing.T) {
	reg := NewDefaultRegistry()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/foo/bar.json", nil)
	data := map[string]interface{}{"foo": "bar"}
	status := http.StatusOK
	err := Respond(w, r, reg, data, status, "")
	if !assert.NoError(t, err) {
		return
	}
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "{\"foo\":\"bar\"}", string(body))
}

func TestRespondExtensionYAML(t *testing.T) {
	reg := NewDefaultRegistry()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/foo/bar.yml", nil)
	data := map[string]interface{}{"foo": "bar"}
	status := http.StatusOK
	err := Respond(w, r, reg, data, status, "")
	if !assert.NoError(t, err) {
		return
	}
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "foo: bar\n", string(body))
}
