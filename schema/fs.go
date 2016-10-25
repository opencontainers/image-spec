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
		modtime: 1477035771,
		compressed: `
H4sIAAAJbogA/5SRsW4DIQyG93sKRDL2QodOt+YBOnSsOlAwOUc6TI0zRFXevQJyTU6KquvCYP/f/2P7
u1NKe8iOMQlS1IPSrwninqJYjMBqTzHgQb0lcBjQ2ap6Ktg2uxEmW5BRJA3GHDPFvlV3xAfj2Qbpn19M
q20ah35G8mAMJYhuTssVa2qDkz2AcTW/kXJOUFj6PIKTVktMCVgQsh5UGUcp7RisgP8t3KFZGGOzq/VA
PFkpHW8FesEJdO1dmkTbk4zEf1gt1exGFHByYljLUF6rvO7iTr1lCPXzEHLf2rtyhY3xEDBiuVaed7jw
YiIJ+f9eV27hNWIW4vOjMSyzPd/WjQLTfejq2Dniyl26+a1/0AxfJ+R68vcHl7htejl9p9RHd+l+AgAA
//9eW+CxBgMAAA==
`,
	},

	"/content-descriptor.json": {
		local:   "content-descriptor.json",
		size:    637,
		modtime: 1476325139,
		compressed: `
H4sIAAAJbogA/5SRP0/DMBDF93yKk9uR1gwVQ1bYGWBDDKn93F6l2sY+hoL63dHFLaQwAFt0er/3J37v
iIxHdYWzcIqmJ3OfEW9TlIEjCukXotDdSZQKPWQ4DuyGkbhSi3l1W+wHxbciubd2V1NctOsylY31ZQiy
uF7Zdps1jv0Zqb21KSO6c3Idsaa2/jO9cXLIUDKtd3DSbrmkjCKManrSYURmD8/DYxO308+9sgWNOjWl
FEgPBQEF0cHTJGPE5wVBOY9QF7wfNljq2JlH4MjqWu1X7kgdG2wqv/3WRCXEkdYHQf1nm9bDTotwlJvV
ZQnPG1SZ1vjLohPVnLqTmyl4eeUCfcWn7398uvkyuiN67o7dRwAAAP//zGqYSn0CAAA=
`,
	},

	"/defs-config.json": {
		local:   "defs-config.json",
		size:    2270,
		modtime: 1476325139,
		compressed: `
H4sIAAAJbogA/+RVzY7TMBC+5ylGhmNhL4jDXrscUZEq4IBQ5SbjdpbYY8YTIEJ9d5Rs6eanDd1detpD
VcX29zMznvHvDMAUmHKhqMTBXIO5QUeBmq8E0YpSXpVWQBkWEcOcg1oKKDDn4GgDy4g5Ocpti5/dER4Y
zDU0GgAmb48fvgGM1hEbQV7fYq4ttF2PwhFFCVPnNID5mFB6Kx2OpEJhYw5bu9k97j16lnqIfCnoGmSB
Lr2+TRxeXHV8X1HQt2+m+JY/bfx/nPNYLbdWBiE/ifLdr8gJiw8s+lBWb+OyTejirjTH+cOPU9WwIrY2
s+4WKfqhjYn6AexOiKrUkSnoUJsDLpqYvvQE+nKTFids/sNq3+7A8qSLUJVln6nL8/X4VfHFs439E5eV
HzfJmfE/9tZfOqrPLN8obG7o/PGWdf/3XEaY1aWnjNiCnFtRMUrw5Rt7f34q+p4shsqPKm5KW6Mkczzj
2UDVCH6vSLDo8dznYOgu6zL+zfmWkvZflwcnPRe0iqOmnoresXir7V22iq+UPB4f0rbSLT/q1dy7Wq1H
L+dZaPYex0P6HCj6qPWqreQp+Jq5RBtO9kPW/HbZnwAAAP//X8Yig94IAAA=
`,
	},

	"/defs-image.json": {
		local:   "defs-image.json",
		size:    2550,
		modtime: 1476325139,
		compressed: `
H4sIAAAJbogA/7RWTW/bMAy951cQaoEeEsc9DAMaFAWG9bJTB3SnBdnAynSszpY0SimWFv7vgx3Hn3GL
rNnNosin96hn2i8TABGRk6ysV0aLBYhbipVWxcqBRfZKblJk8AbuLOnPRntUmhi+ZLgmuLckVawkluWz
HV4NIBZQHAEgMooUfttaqkMAQkXFgYn31i3C0FjScg/v5obXoZMJZRiq4qiwgZjtAfwOUDjPSq+buEXv
iUs5P5YYPK+m4fIyuMLg+VPwfT5dTc9FmZrvKkSk1uR8m1qvKT4hkLy13qwZbaIkyITkL7fJYFcLJoYi
yTw8kvQzULpcVkTg4jpBl9wsrhP6gxFJlWFaVd5cHK/nMrj6OZ0Hq+liiUF8GVwNFGWoVUzO31Y6DL+j
8UOsAeOd7ibO9HujmIpzllWs44JZE3TqubOuLqMVsSn62HAmqtCqaQwbS+wVuZa8EbuN3WuZW+jY3yFT
TExaUgQ9WSXEOVNc1J61jN4yZ52a9zW+TaVIK6zzsPXkjqQTUezmj87os7BNTGn/8cNhUgPXj8urUg/C
1LfTAzrOZDVKR9yIvcYtVu4gy0R5kn7D1KkBEMaJVmDVQRxx0xCzv/vPasNXqA7nQWc7Hwo7HS3j3ktm
/kTsdtY+Haka9L3kYsKi46dtWYM6Sg+ZcTvcVp6yQ1zeEAaQvyr0CVmh7r/hR/frtWb9N2EDKIyichhh
+rX9msaYOpqM9aR53j/lnc8kam08dv9VjrZAG6T+MhlNd3FnNB2csiNTO0N7X/blvtedVl9eDg1KvUlT
MRC9muxX+SSf/A0AAP//XF/NGfYJAAA=
`,
	},

	"/defs.json": {
		local:   "defs.json",
		size:    3193,
		modtime: 1476325139,
		compressed: `
H4sIAAAJbogA/7SWQZPSMBTH73yKTPSIbpumLXBzRNc97MDMjifHQy0PiEIS09Rx3eG7OymlNG0oFPGw
C03yfv/3/nmheRkghBeQpYpJzQTHE4SnsGScmacM5RkskF4rka/WItdIrwHNJPD3guuEcVDoSULKlixN
ivDhnlcB8AQZCYQw43pUPSGE9bMEo8a4hhWoIrKY2DLOtvkWT9Abn4yOw8nvctgncTG4G1ZkP+qLDkgc
ueBmvIUPSF888WlMR0FEXRrVZEsoon2FxoQEQUy8IBqFNI4jz/Mcio5VlnR+xe64dEgYtri998YFjsIw
aKN774sLTcmYjqOYjNv83tvh4vsjSqOYUi8OYm8chsRlvh/NheGquqDgMFviCfpSDqBqqph+rcDM4ld3
tQN3V1peLdwN3dGHani+2dRWl9++Ory4ZYYRvW2GmVaMr/pneNDYx982p20inwrs/r+rlcS375DqYyfJ
RGtQfK6EBKUZZLUghPDbF3+4s4a6ShjUP1tJzfbS/zmpknUmqc8P03oip9smINgKvL828J1SyfNsef8w
zVwWJGb66ADTsG1U7ZYy+TgrLPX21l8p6d5oW+cj28A80WuXQj3usP4D/3XewEbyFuCTED/qBHmiS6Sd
VIdaVYLrMOJErbLLOK6smzSwyu+AGZ/ONHGpZwy5aU8VDjsVH6aPiZSNn5ZTG7AWmbZOS4dk/cQ0LUsP
N79bwDL2B66kOB15FDnXl7iRiVylF0p3N+QCMs14Ul6e/xknpH1v7kRd0uJl+/V7ZVRdquBnzhQsrDeo
VfOwZWpT3H45DszfbvA3AAD//0JyEpx5DAAA
`,
	},

	"/image-manifest-schema.json": {
		local:   "image-manifest-schema.json",
		size:    1032,
		modtime: 1476325139,
		compressed: `
H4sIAAAJbogA/6RSvW7jMAze/RSEk/Ec3XCT15tuOHRo0KXooNqUzSCWVFItEAR590JSlNpIh7ZZKX6/
4rECqHuUjskHcrZuob7zaP86GzRZZPg36QHhv7ZkUALce+zIUKfT9q8IX0s34qQjdAzBt0rtxNkmTzeO
B9WzNqH5/Ufl2SrjqC8QaZVyHm1XVCXB8rai6EBNZwcZGw4eI9o977A7zzw7jxwIpW4hBgOoM8UDsuRw
eXydeTuSgCHc9yA5IQqEESGJQxGHzAdvmRC0gLZANuCAnFwk8p8lU0uzF7aStcik+Sk/1xP2pLd543ij
/gfXRXvNaCJXj0aatL6Jn7tSPRqyFOuTGW5hrXPW0DD3VdiiG7ShKZ/gOLEu4Xt9QJY5vBShmfVh1nbA
ab73VSGA00JQW+tCumu5vcxmTxLUnPJ7nc6RVfGanNaML6/EGG09fnblV5ex/I9lvRXAU3Wq3gMAAP//
X3p8DwgEAAA=
`,
	},

	"/manifest-list-schema.json": {
		local:   "manifest-list-schema.json",
		size:    1010,
		modtime: 1476325139,
		compressed: `
H4sIAAAJbogA/6ySMU/7MBDF93yKU9rxn/o/MGWFBQnEQMWCGExybq5qbOM7kKqq3x3ZTkqiMoDoevF7
93vvcigAyha5CeSFnC1rKB882mtnRZPFALe93iDca0sGWeCOWODRY0OGGp0k/6LHkpsOex31nYivldqy
s1WerlzYqDZoI9X/K5Vni6yjdpRwrZTzaJtxNSdZfq0oYqh+wKh2xJINZO8xWrjXLTbDzAfnMQghlzXE
iABl9nnCwDlmHp+nX3fEYAh3LXCOiQzSISQCmBFANoWP7AqaQVsgK7jBkFDShj9kVHPsk+WYetyV5sf8
ueyxJb3OLw6XgPgyPAEsA5po2KLhKmlW8eAL1aIhS7FNnujmfIM5T/nGRDoEvZ90J9hP3/149bDjZriu
GzoCOM5YtLVO0n/Ml2pravm7vqbKYmRNpGXAt3cKGLGev/uhz05/1nUB8FIci88AAAD//7jpza7yAwAA
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},
}
