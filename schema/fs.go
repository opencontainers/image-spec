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
		size:    774,
		modtime: 1478057674,
		compressed: `
H4sIAAAJbogA/5SRvW7rMAyFdz+F4WS8ju7QKWsfoEPHooMqUzEDWFRJZgiKvHv1EzcxEBTuEsSH/M4R
ya+mbbsBxDFGRQrdvu1eIoRnCmoxALfpn8dD+xrBoUdnS9e/jG3FjTDZjIyqcW/MUSj0Vd0RH8zA1mv/
/8lUbVM5HGZEEkMpzc1pUrDabXCyBzCu5FdSzxEySx9HcFq1yMmBFUFSJY+TNMdgFYYf4Q4VZQzVruie
eLKaK0NCesUJulK71JbOnnQk/sVq2c1uRE2POzGsZUjWdl53cde9ZfDl8eClr+VdvsLGJAUD5mvJvMOF
FxOpl797XbmF14iixOdHY1hme76tO+1mug9dHTtHXLlLM/+WN3QMnyfkcvK3B5e4bXo5ffp4by7NdwAA
AP//XlvgsQYDAAA=
`,
	},

	"/content-descriptor.json": {
		local:   "content-descriptor.json",
		size:    836,
		modtime: 1478141504,
		compressed: `
H4sIAAAJbogA/5SSP2/iQBDFe3+KkaE88BXoCtpLnyLpohSLd2wPwrvO7CDkRHz37HhtMImUPw3Cb99v
5j173zKA3GIomToh7/It5Pcduv/eiSGHDPoPncDdaPIMDx2WVFFpBuKPjliGssHWKN6IdNui2AfvVkld
e64Ly6aS1d9NkbRF4shOSIiMj5vLaXMYsOQu7GV74qTvUEm/22MpSes48iyEIZ5osai1aMk8JnOSPveV
BmHw6VDwFajAWCHHNGhhtmPAl/FMOYtVWFFralxr2UV8Jkc6NRTXvQN1TnAe6PW7JGoBcrDrBcMv06Qc
xTwIOfm3uQ1hqcYg8xg/aTRSN5OOfAhf1DFwoCDaQY1QsW/h1FDZxE4UxibQmh52CNaf3MEbi/babfrK
htn0V5kE2/nemTMIk6sv1nhSeW6Nto1pKR/1czb9DlVyxpcjMeplfPp4ceaf7vYNxofn7Jy9BwAA//9L
DLQ9RAMAAA==
`,
	},

	"/defs-config.json": {
		local:   "defs-config.json",
		size:    2270,
		modtime: 1478057674,
		compressed: `
H4sIAAAJbogA/+RVzY7TMBC+5yksw7GwF8SB6y5HVKQKOCBUucl4O0vsMeMJEKF9d5xs6cZJG7q79MSh
ajLx9zMz9vhXoZSuIJaMQZC8fqP0FVj02L1FFQwLlk1tWAmpZQB/SV4MemCVnixeq1WAEi2Wpscv7gj3
DImw00jBsl++f08RaQN0grS5gVJ6aB8PTAGSLsTB6hT/EIGzyIAjCqO/1vtPt4t73DtwxO0Y+ZzBdshk
Nr68ieSfXQx8X6CX16/m+FY/TPh3nJehWW0Nj1J+EuXbn4EiVO+J5aGszoRVX9DlXWsO8/vvx7phmE2r
F8NPKODGNmb6l6SOiAq3gVLaY23ysOxy+pwJ5HKzFmds/sVqbndkedaFb+o6ZxryfDm8VVz13+b+kerG
TQ/Jifk/dtefO6tPxF+T7BWePt6K4f+OSzOR2PiUEVuhtWusJgU+/8HerZ/LPpMF37hJx3VtWuCoD1e8
GKlqhm8NMlQZz30Nxu6KIeOfmm8xSn67PLjoJYMRmBzquewtsTPS7+UEfSHo4PCQNo1s6VG35s7VejO5
OU9Ck3MwHdKnQMEFadd9J4/BN0Q1GH/0PBTd77b4HQAA//9fxiKD3ggAAA==
`,
	},

	"/defs-image.json": {
		local:   "defs-image.json",
		size:    2753,
		modtime: 1478141527,
		compressed: `
H4sIAAAJbogA/7RWy27bMBA8x19BKAFy8EM9FAVqBAGK5tJTCqSnBm6xoVYWU4lUSTqpE/jfu5RkmXol
UOOebC6Xw5nhLqnnCWNBhIZrkVuhZLBkwRXGQgo3MiwHbQXfpKCZVew6R/lZSQtComZfMlgju8mRi1hw
KJbPSrwagPDcFhTMMBLwbZtjHaKgiNyGibW5WYahIni+hzcLpdeh4QlmEAq3VXiAmO0BbAkYGKuFXB/i
OViLupDz4xbmT6tpePtu/pH+fZp/X0xX07OgSN2VK4JIrNFYn1rLFJsg43qbW7XWkCeCM2LGf5lNxsq1
TMXMJam7e+R2xoQshhURdn6RgEkulxcJ/oGILMsgrVZeno/XQ1p+Thfz1XRJo5hGHUUZSBET+lWlQ+k3
GN/F6jAudR/iGn9vhEa3z20Va1TB7BA04qkxrg7Di+Qp2FjpLKhCq4MxmshTjaLx5A2U29C5FrlOx/4M
NcaoyRKMWEtWAXFG827tqVfoXnHWqbu2xtepuDRXOndbi2YkHWJjFvdGydPQJyak/fC+n1Sn6oflVak+
zEmw0al5RRSwVJTN4ZJZrFXGHql9ElImTKWHZbBld8gi9ShTRe0RNRXuawy0hm1zSljMSg4njdS6gcqw
Kx5wSomzCFxw1+tIXWgtVeP6pUbpVdF3hr3dUsyA5gmJ5HajsbGG5pQJvMCqgTjQGF3M9uw/qw1foNo9
mcb0rivseLQI7I1kFg+0Q1nQxyNVg76VXIzgHD+uZQfUQXrdbiy3rDuSsRHC/I7sE/oAWoBsX1aj/XrJ
rP8mrAMFUVTcq5B+9ds0htTgZMgT78aa+L/7Fx+kVBaan12jS8AHqR9ZJfE6blxNvQ/GwAOUQX5T+HLT
csfz5bnvopSbNA06oleT/Wg32U3+BgAA//8jeF0swQoAAA==
`,
	},

	"/defs.json": {
		local:   "defs.json",
		size:    3193,
		modtime: 1470056192,
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
		modtime: 1470056192,
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
		modtime: 1470056192,
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
