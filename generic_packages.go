//
// Copyright 2021, Sune Keller
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// GenericPackagesService handles communication with the packages related
// methods of the GitLab API.
//
// GitLab docs:
// https://docs.gitlab.com/ee/user/packages/generic_packages/index.html
type GenericPackagesService struct {
	client *Client
}

// PublishPackageFileOptions represents the available PublishPackageFile()
// options.
//
// GitLab docs:
// https://docs.gitlab.com/ee/user/packages/generic_packages/index.html#download-package-file
type PublishPackageFileOptions struct {
	Status *GenericPackageStatusValue `url:"status,omitempty" json:"status,omitempty"`
}

// PublishPackageFile uploads a file to a project's package registry.
//
// GitLab docs:
// https://docs.gitlab.com/ee/user/packages/generic_packages/index.html#download-package-file
func (s *GenericPackagesService) PublishPackageFile(pid interface{}, packageName, packageVersion, fileName string, content io.Reader, opt *PublishPackageFileOptions, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf(
		"projects/%s/packages/generic/%s/%s/%s",
		pathEscape(project),
		pathEscape(packageName),
		pathEscape(packageVersion),
		pathEscape(fileName),
	)

	// We need to create the request as a GET request to make sure the options
	// are set correctly. After the request is created we will overwrite both
	// the method and the body.
	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, err
	}

	// Overwrite the method and body.
	req.Method = http.MethodPut
	req.SetBody(content)

	return s.client.Do(req, nil)
}

// DownloadPackageFile allows you to download the package file.
//
// GitLab docs:
// https://docs.gitlab.com/ee/user/packages/generic_packages/index.html#download-package-file
func (s *GenericPackagesService) DownloadPackageFile(pid interface{}, packageName, packageVersion, fileName string, options ...RequestOptionFunc) ([]byte, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf(
		"projects/%s/packages/generic/%s/%s/%s",
		pathEscape(project),
		pathEscape(packageName),
		pathEscape(packageVersion),
		pathEscape(fileName),
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var f bytes.Buffer
	resp, err := s.client.Do(req, &f)
	if err != nil {
		return nil, resp, err
	}

	return f.Bytes(), resp, err
}
