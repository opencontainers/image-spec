// Copyright 2016 The Linux Foundation
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

package schema

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(f)
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/config-schema.json": {
		local:   "config-schema.json",
		size:    707,
		modtime: 1463700693,
		compressed: `
H4sIAAAJbogA/5SRPW7DMAyF5/oUhpOxjjp0ytoDdOgJVJmKGcCiQDJDUPju1U/c2kvhLobx+L73JOqr
adtuAHGMUZFCd2679wjhjYJaDMBt+vN4aT8iOPTobHE9Z+woboTJZmRUjWdjrkKhr+qJ+GIGtl77l1dT
tUPlcFgQSQylNre0ScGq2+BkL2Bc6a+k3iNklj6v4LRqkVMCK4KkSb5O0hyDVRh+hBUqyhhqXNE98WQ1
T4aE9IoTdGU2V0tnbzoS/xG1dbMbUdPhbgx7GZK9zscuVu4jgy+HBy99HZ/yKxxMUjBgfi1ZdrjJYiL1
8v+sB7fJGlGU+F7CnlZ3sMz2/rvrtJhp3bi7c8l/cHOzfOdmbr4DAAD//1EkCUvDAgAA
`,
	},

	"/defs-config.json": {
		local:   "defs-config.json",
		size:    1755,
		modtime: 1463700704,
		compressed: `
H4sIAAAJbogA/7xUzc7TMBA8N09hBY6BXhAHri1HVKQIOFZusm63xF5rvQEi1HfHSauS/qRKKV8PVeOx
Z2bXY/t3olRaQigYvSC59INK52DQYTsKymsWLOpKsxJSCw9uRk40OmAVvwyuVe6hQIOF7vjZXvCoEAVb
jwgW3fLjOCLSeGgNabWFQjpqh3smD9EXQm91xL8E4BOkpxGE0a3T49Qu+8v7BJa4OWe+ZjAtMxYb3m4D
uVfTXt1TdPL+3S29/Kf2/09z5ut8o/ms5YckP/7yFKD8TCz3qlrt825DF/toruu7H0NpaGbdpFl/CgXs
eRk38otWA6bCjafY9vO9Z7Z8vulXqmp797EYFeA34u9xyRzH36qk/3/QSplITHjkZpdozBLLiy5ffnsr
3QAP+o7rf4NBTh+YuzegYNACg8frUEeWTCYRNcRWS5d+JL0RtHA9YF3Lhv7pyTzUs1xdPJuj2GQtDN/Q
W1SwXppll8oQfUVUgXaDqSTtb5f8CQAA//8Cok052wYAAA==
`,
	},

	"/defs-image.json": {
		local:   "defs-image.json",
		size:    3100,
		modtime: 1462965360,
		compressed: `
H4sIAAAJbogA/+xWy27bMBC86ysIJUAOfqjXGkGAoEGBnlIgPdVQgY20kphKpLqkgzqB/r3U+x3Xidte
erK55A5nhmPSzxZjto/KI55qLoW9YfYNBlzwfKRYCqS5t4uBmJbsNkXxQQoNXCCxTwmEyO5S9HjAPSja
lyVeA2Dw8i1MMUGfw5d9ik3JFLmfbxhpnaqN40gD79Xwai0pdJQXYQIOz7dyWohlDaBLQFtp4iJs6ylo
jVTI+baF1ZO7cLbvVu/Nt+vV1/XCXZzbxdKs7LB9HqLSXWoDU3SEzKN9qmVIkEbcY4aZ913tElb2Mhmw
fJG8f0BPLxkXxbAiwi4uI1DR1eYywp/gG8sSiKvOq4vj9Rgt7mJjvgXXq4/FYCiooi/p9X53MEYES5lt
nfDHjhPm+Nuq1jv0ZVtU/Kk3rryvCm6rmQxBEz9UHQkzSZo7smJtzrk+HsIAychGnw0kFBDnZj7vPetk
uJO7Zmk21HOYSr4sT8X9XqP6TTq121xoDJGm9x9ld15J32oDY3U/6+wkIHhg1t2cIEMTWH8hS51KGoMO
JCX/8/Uv8jV1EAOg4/LUoEzqmNI4maZiBsiLuDYNO8Jej5mTyu4U3B7iTHDGmMPZV6t1XqA6fjV609lY
2OloGbA3klk/mh3KFJ+OVAP6VnIBQu74aS1rUWfpARHsx9MmAckUlwPC2nt+WugjEAcx/IUf7ddLZv0x
YSMo8P3iMoL4c/dnGkCs0JrzJDvwIoIQUkP/H+3REeiCNI+QFHgb9K6myVvWXLJq/aCkOHN6Lwekd4Uv
dwN3Or48T12UYhfH478BrlWPMiuzfgUAAP//VjUNyBwMAAA=
`,
	},

	"/defs.json": {
		local:   "defs.json",
		size:    3193,
		modtime: 1463700693,
		compressed: `
H4sIAAAJbogA/7RWTXPaMBC98ys8tEfa2PIX9NYp/cghAzOZnjo9uGYBtSCpstxpmuG/VzLGWPZiMKWH
JPau9r23T6tYzwPHGS4gSyUVinI2fOMMp7CkjJq3zMkzWDhqLXm+WvNc6UdwZgLYO85UQhlI51FASpc0
TYry0R6vAtB4hkIHKVPj6k2/qycBhk3HYQWyqCwSW127zbc698oj42M4+V2GPRIXwd2oQvaivtA+iSMM
3MRb8D7pC0+8IA7GfhRgHFWyRRQFfYkmhPh+TFw/GodBHEeu6yKMyCqLOr9idzAeEoYt3N57gwFHYei3
oXvvCwYdkEkwiWIyaeP33g4M3xsHQRQHgRv7sTsJQ4KZ70VzbnBlnZAzmC114EsZcKpUkX4pwWSHL+5q
B+6utLxauBvh1YduWL7Z1FaXT18RL26pUDt7U4WZkpSt+is8cOzrb6tpm4jHAnb/Gxsl/u07pOo4SSJR
Wj+bSy5AKgpZrUinXz97o50V6mphUP/bEjXbU/9nUSXWGVGf76d1IafHRh94q/DjtYVvpUyeZktdn2EW
JCZ9dIAq2Da6xqmMHrTDkm9v/ZWU+EbbPB/oBuaJWmMM9brD+vfs13kDG+ItgE+c/6gjiBNTImxRHWxV
C9hh1DatsstwMNVNNLDa7wAzPp0Z4pLPGHLTmSocRhnvpw+JEI1/Lac2YM0zZZ2WDsr6iWlalh5ufrcA
y+gfuBIFdeSB50xd4kbGc5leSN09kPryrChLysvzP8NxYd+bO6EuGfFy/Pp9MqoplfAzpxIW1hfU6nnU
MrVJbn8cB+ZnN/gbAAD//0JyEpx5DAAA
`,
	},

	"/image-manifest-schema.json": {
		local:   "image-manifest-schema.json",
		size:    1064,
		modtime: 1462965360,
		compressed: `
H4sIAAAJbogA/6RTvVLjMBDu/RQ7TspzdMVVaa+64oaCDA1DIeyVvZlYMlrBTCaTd0c/UZAJBSSlV/v9
Sj5UAHWH3FqaHBldr6G+m1D/NdpJ0mjh3yh7hP9Sk0J2cD9hS4paGbd/BfiS2wFHGaCDc9NaiC0b3aTp
ythedFYq1/z+I9JskXDUZQh7jPGqbVblCEvbgoIDMZ4cJKzbTxjQ5nmL7Wk2Wc9hHSH7kxDMzxLFg2dM
4dL4MvNmIAZFuOuAU0JkcANCFIcsDokP3hIhSAapgbTDHm10EcmvSybmZs9sOWuWifNjOq5H7Ehu0sbh
Rv0PrrP20qIKXB0qbuL6KlzuQvgBaQr1cYGbWfOaivrS17fY8s2YT0l3cu/tl3S5GGmt3BftOxzLvWuF
vfTMgNTauPju+faymx35xkvKn3VeIqvsNTqtLb68ksVg6/Grv+Di5czva163/3iqjtV7AAAA///++ypf
KAQAAA==
`,
	},

	"/manifest-list-schema.json": {
		local:   "manifest-list-schema.json",
		size:    1010,
		modtime: 1462965360,
		compressed: `
H4sIAAAJbogA/6ySMU/7MBDF93yKU9rxn/o/MHWFBQnEQMWCGExyaa5q7OAzSFXV747tS0qiMIDoUqkv
fu9+7+xjBpBXyKWjzpM1+Rryhw7NtTVek0EHt63eItxrQzWyhzsKP48dllRTqZPlX8xYctlgq6O/8b5b
K7VjawpRV9ZtVeV07Yv/V0q0hfioGiwcPDaMLofRnGxyWlHEUG2PUewDhgT4Q4cxwr7usOy1zoUg5wk5
fIkVgyY5TyFWaoo8b79piKEm3FfAUhMZfIOQCGBCABIKH5IKmkEbIONxiy6hpAl/6Kim2OfIofUwK+kn
+Zy3WJHeyInjJSC+As8AS4d1DKyw5iJ5VvHCFyoIZChuk0e+KV8fzmO+oZF2Th9Gu/PYjs/9eHQ/46a/
XdvvKFBMWLQx1qd3zJfa1jjyd/saO7OBNZHmDt/eyWHEev7uQc+ufrbr8P8lO2WfAQAA//+46c2u8gMA
AA==
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},
}
