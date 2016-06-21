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
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
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
		size:    710,
		modtime: 1466466955,
		compressed: `
H4sIAAAJbogA/5SRPW7DMAyFd5/CcDLWUYdOWXuADj2BKlMxA1gUSGYICt+9+olbGygKdzGMx/e9J1Gf
Tdt2A4hjjIoUunPbvUUIrxTUYgBu05/HS/sewaFHZ4vrKWNHcSNMNiOjajwbcxUKfVVPxBczsPXaP7+Y
qh0qh8OCSGIotbmlTQpW3QYnewHjSn8l9R4hs/RxBadVi5wSWBEkTfJ1kuYYrMLwLaxQUcZQ44ruiSer
eTIkpFecoCuzuVo6e9OR+I+orZvdiJoOd2PYy5DsdT52sXIfGXw5PHjp6/iUX+FgkoIB82vJssNNFhOp
l/9nPbhN1oiixPffrmGZ7f1n3Wk307p0d+1S8eDmZvnOzdx8BQAA//964XeexgIAAA==
`,
	},

	"/content-descriptor.json": {
		local:   "content-descriptor.json",
		size:    616,
		modtime: 1466180793,
		compressed: `
H4sIAAAJbogA/4yRMVPDMAyF9/wKXdqR1gxMXWFngI1jcG0lVe9iG1kMhet/x4oTSGGgW/L8vvck+7MB
aD1mx5SEYmh30D4mDPcxiKWADPqFQeBhMkWGp4SOOnJ2JG40Yp3dAQer+EEk7Yw55hg2Vd1G7o1n28nm
9s5UbVU58jOSCxNLs5ub84hVt/Hf7ZWTU0Il4/6ITqqWuPAshLmc6GJFG9CTfa7mKv3dVw4Io09DIXag
AmOHXKZBD4uOEV+XM+U8dnlDg+1xq8uuyj8F0tRsfnpH6lzhNtPHf5OoBSjA/iSYr5hmvgkqz9QjX/Z5
6jHLsvGa4SeqJjVTWsv49k6M+mAvvy93ud5ldfl5bc7NVwAAAP//Zc2MR2gCAAA=
`,
	},

	"/defs-config.json": {
		local:   "defs-config.json",
		size:    2154,
		modtime: 1466467007,
		compressed: `
H4sIAAAJbogA/+RVzY7TMBC+5ymswLGwF8SB6y5HVKQIOCBUucl4O0vsscYTIEL77jjZquSnCWlLT3uo
mkz8/cyMPf6dKJUWEHJGL0gufafSOzDosHkLymsWzKtSsxJSaw/ulpxodMAqPhm8V5mHHA3musWvnggP
DJGw0YjBvF1+eI8RqT00grR9gFxaaBv3TB6iLoTO6hj/FIB7kQ5HEEZ3nx4+Pa7+4j6AJa6HyJcMpkFG
s+H1QyD34qbj+wadvH0zx5f91P7/cd76KttpHqR8EeX7X54CFB+J5VRWq33WFnT91Jrj/O7HVDc0s67T
VfcTCtihjZn+RakJUeHaU0x7qE0O1k1OX3sCfblZizM2/2G1b3dgedaFq8qyz9Tl+XZ8q9ji2eb+mcrK
jg/JwvzP3fXXzuoL8fcoe4fLx1vS/d9zpUwkJlwyYgs0ZoPFqMDXP9ilroEndZflv8Mg/Ul/cgFyBi0w
OmADH70CGGKrpd1XEfpK0MLxgakr2dFZN9je1WY7usUWoclaGA/MJVCwXupN25sp+JaoBO0me5M0v8fk
TwAAAP//dkZ6ZWoIAAA=
`,
	},

	"/defs-image.json": {
		local:   "defs-image.json",
		size:    2528,
		modtime: 1466438460,
		compressed: `
H4sIAAAJbogA/7RWy27bMBC8+ysIJUAOfqjXGkGAoEGBnlIgPTVQgY20kphKpLqkgzqB/r2k3k+nrt0b
ueQOZ4Zjym8LxpwAlU8801wKZ8ucOwy54HamWAakub9LgJiW7D5D8UkKDVwgsS8pRMgeMvR5yH0o2lcl
XgNg8OwRpphiwOHbPsOmZIo8sAfGWmdq67rSwPs1vNpIilzlx5iCy+1RbguxqgF0CegoTVxEbT0DrZEK
OT8eYf3qLd3HD+uPZnS7/r5ZestLp9ialx1OwCNUukttYIqOkfm0z7SMCLKY+8ww83+qXcrKXiZDZjfJ
p2f09YpxUUwrIuzqOgYV32yvY/wNgbEshaTqvLk6Xo/R4i23ZhTerj8Xk4GgFAQPDfhdJUPSCb6PsUaE
S9ltnfDXjhPacx6rWi8Eq7ao+GtvXt1Fp5IloENJqVOVvNYXMuRNRFF15M2kbe5ai71WR32FhCGSsQQD
NpBVQFyaddt70cl5J5vN1nyo8X0qdptNztNeo/pLOvUNcKExQpo+f5TveSXV1kmY5iIGQMflqUGZ1DGl
cTJNxQqQH3NtGnaEvR6zJpXTKXg9xJngjDGHq/+s1j1AdfzL7y3nY2Hno2XATiSzeTEnlCk+H6kG9FRy
IYJ1/LyWtaiz9IAI9uNlk4B0iss7woy0g0JfgDiI4S/8aL8OmfXfhI2gIAiKxwiSr92faQiJwsWcJ+24
HuW9LyIIITX0/5UcHYEuSPMRkgLvw97TNPnKmkdWbZ6VFBdu78sB2UPhy8PAnY4vb1MPpdgliTMS7S3q
Wb7IF38CAAD//wKthPngCQAA
`,
	},

	"/defs.json": {
		local:   "defs.json",
		size:    3193,
		modtime: 1466180793,
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
		size:    1032,
		modtime: 1466180793,
		compressed: `
H4sIAAAJbogA/6RSPU/zMBDe8ytOacc39TswdWViQAxULIjBJOfkqsYOPoNUVf3v+KMujsoAdMyTez7u
8R0qgLpDbi1Njoyu11A/TKhvjXaSNFq4G2WPcC81KWQHjxO2pKiVcfpfoC+5HXCUgTo4N62F2LLRTUJX
xvais1K55v+NSNgi8ajLFPYc413b7MqRlqYFhQRiPCVIXLefMLDN6xbbEzZZr2EdIfs/YTGPJYknr5iW
S/DlzpuBGBThrgNOGyKDGxCiOWRzSHrwkQRBMkgNpB32aGOKKP63zcQ87Fkt75ptIn5Mv+sRO5KbNHG4
0v9L6+y9tKiCVoeKmzi+Co+7EB4gTaE+LnizaN5TUV/mymohDWrX5EcwNqrO6Tu593FLei5CWiv3RdsO
x3Lup0beamYotTYu3jVfX2azI99oKfm7TktmlbPGpLXFt3eyGGI9f3flF5cxf495vf7jpTpWnwEAAP//
X3p8DwgEAAA=
`,
	},

	"/manifest-list-schema.json": {
		local:   "manifest-list-schema.json",
		size:    1010,
		modtime: 1466180793,
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
