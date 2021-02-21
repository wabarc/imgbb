// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package imgbb // import "github.com/wabarc/imgbb"

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wabarc/helper"
)

const (
	IMGBB_URI = "https://imgbb.com/json"
	IMGBB_API = "https://api.imgbb.com/1/upload"
)

type ImgBB struct {
	Key string

	Client *http.Client
}

func NewImgBB(client *http.Client, key string) *ImgBB {
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &ImgBB{
		Key:    key,
		Client: client,
	}
}

func (i *ImgBB) Upload(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if size := helper.FileSize(path); size > 33554432 {
		return "", fmt.Errorf("File too large, size: %d", size)
	}

	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		field := "source"
		if i.Key != "" {
			field = "image"
		}
		m.WriteField("key", i.Key)
		m.WriteField("type", "file")
		m.WriteField("action", "upload")
		part, err := m.CreateFormFile(field, filepath.Base(file.Name()))
		if err != nil {
			return
		}

		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}()

	endpoint := IMGBB_URI
	if i.Key != "" {
		endpoint = IMGBB_API
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, r)
	req.Header.Add("Content-Type", m.FormDataContentType())
	req.Header.Add("Host", "imgbb.com")
	req.Header.Add("Origin", "https://imgbb.com")
	req.Header.Add("Referer", "https://imgbb.com/")

	resp, err := i.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return i.Parse(resp)
}

func (i *ImgBB) Parse(resp *http.Response) (string, error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		return "", err
	}
	if i.Key == "" {
		if err, no := dat["error"].(interface{}); no || err != nil {
			return "", fmt.Errorf("%v", dat["status_txt"])
		}

		if success, has := dat["success"].(interface{}); !has || success == nil {
			return "", fmt.Errorf("%s", "Upload Failed")
		}

		img, has := dat["image"].(map[string]interface{})
		if !has || img == nil {
			return "", fmt.Errorf("%v", "Unrecognized")
		}

		if url, has := img["url"].(string); has || len(url) == 0 {
			return url, nil
		}
	} else {
		if status, has := dat["status"].(int); !has || status != 200 {
			return "", fmt.Errorf("%s", "Upload Failed")
		}

		if success, has := dat["success"].(bool); !has || success == false {
			return "", fmt.Errorf("%s", "Upload Failed")
		}

		img, has := dat["data"].(map[string]interface{})
		if !has || img == nil {
			return "", fmt.Errorf("%v", "Unrecognized")
		}

		if url, has := img["url"].(string); has || len(url) == 0 {
			return url, nil
		}
	}

	return "", fmt.Errorf("%s", "Failed")
}
