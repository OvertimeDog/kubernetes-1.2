// Copyright 2015 Google Inc. All Rights Reserved.
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

package storage

import (
	"bytes"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

var ArgDbUsername = flag.String("storage_driver_user", "root", "database username")
var ArgDbPassword = flag.String("storage_driver_password", "root", "database password")
var ArgDbHost = flag.String("storage_driver_host", "localhost:8086", "database host:port")
var ArgDbName = flag.String("storage_driver_db", "cadvisor", "database name")
var ArgDbTable = flag.String("storage_driver_table", "stats", "table name")
var ArgDbIsSecure = flag.Bool("storage_driver_secure", false, "use secure connection with database")
var ArgKubeletIp = flag.String("storage_kubelet_address", "127.0.0.1", "kubelet api address")
var ArgKubeletPort = flag.Int("storage_kubelet_port", 10255, "kubelet api port")
var ArgDataId = flag.Uint64("storage_dataid", 0, "report data to the storage by this dataid, default is -1")
var ArgDbBufferDuration = flag.Duration("storage_driver_buffer_duration", 60*time.Second, "Writes in the storage driver will be buffered for this duration, and committed to the non memory backends as a single transaction")

type Uri struct {
	Key string
	Val url.URL
}

func (u *Uri) String() string {
	val := u.Val.String()
	if val == "" {
		return fmt.Sprintf("%s", u.Key)
	}
	return fmt.Sprintf("%s:%s", u.Key, val)
}

func (u *Uri) Set(value string) error {
	s := strings.SplitN(value, ":", 2)
	if s[0] == "" {
		return fmt.Errorf("missing uri key in '%s'", value)
	}
	u.Key = s[0]
	if len(s) > 1 && s[1] != "" {
		e := os.ExpandEnv(s[1])
		uri, err := url.Parse(e)
		if err != nil {
			return err
		}
		u.Val = *uri
	}
	return nil
}

type Uris []Uri

func (us *Uris) String() string {
	var b bytes.Buffer
	b.WriteString("[")
	for i, u := range *us {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(u.String())
	}
	b.WriteString("]")
	return b.String()
}

func (us *Uris) Set(value string) error {
	var u Uri
	if err := u.Set(value); err != nil {
		return err
	}
	*us = append(*us, u)
	return nil
}
