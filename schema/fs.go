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
		size:    3044,
		modtime: 1463695223,
		compressed: `
H4sIAAAJbogA/6xWT3PaPhC98yk8/H5H2tiybENvnaZ/csiEmUxPnR5cs4BakFRZ7jTN8N0rGWMse3Ew
5QBYu9r3dt+usJ5HnjdeQJ4pJjUTfPzGG9/CknFmV7lX5LDw9FqJYrUWhTaP4D1I4O8E1ynjoLxHCRlb
siwtwyd7vBrA4FkKY2RcT+uVWesnCZbN2GEFqowsHVsTuy22xvcqINOjOf1dmQOSlMbdpEYO4qHQIUli
DNzaO/AhGQpPAprQaRhTjKN2dohiOpRoRkgYJsQP42lEkyT2fR9hRHY51MUF3cF4SBR1cAf3BgOOoyjs
Qg/uCwZNyYzO4oTMuviD24HhB1NK44RSPwkTfxZFBBM/iOfC4qomoeDwsDSGL5XBq12l+38F1jv+76Zx
4G4qyeuNuwkefaiGF5tNY3f19BXR4poZGmWvmmGuFeOr4RkeOPbx181pm8rHEnb/jY2S+PYdMn2cJJlq
kz+fKyFBaQZ5I8i4Xz8Hk51j6ith1Pw9JPX57raZyOkOmbPlBH68NPCtUunTw9LE55gEqXUfFWAatq2q
cSqbD1phxbcX/UJKXFOX5wPbwDzVa4yhGXfY/57/elnAVvIOwCchfjQR5IkpkW5SPWx1CdjcG5lW+Xk4
WNZtNHDK7wGzOr0wxBWfFeSqM1UqjDLe3d6nUrZO8akGrEWundPSQ9k8MW3JssMl6xpgOfsDF6KgityL
gutz1MhFobIzqfsH0txTNeNpdU/9Zzgh3StqL9Q5I16N37B/53pKFfwsmIKF87Jyap50RG2Tu++hkf3s
Rn8DAAD//2VgiEzkCwAA
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
