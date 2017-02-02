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
		modtime: 1483690724,
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
		size:    939,
		modtime: 1486025448,
		compressed: `
H4sIAAAJbogA/5STTW8TMRCG7/kVI7dSL03NoeKwqnqBOwe4IQ6OPV5PiT8YT1QWlP+OvN5tEoGoeomS
V+/jeWbj/b0BUA6rZSpCOakB1KeC6UNOYighQ/uGSeDjUsoMnwta8mTNTNy2I66rDRhNw4NIGbR+qjlt
e3qXedSOjZftu3vds6vOkVuROmidCya7Tq4z1tvavUzvnEwFG5l3T2ilZ4VzQRbCqgZoiwGoiI7Ml17u
0d/7SkCYe+1QyB5awOiRMVl0cDZjxq8ZfeMc+rqlaEa8a8teaYeeErVjqz4NnrFjp1WlX6+ptApQgt0k
WN+o8w8RSvL+/lLC0YhVXtGwPBXJI5sSyIINaL/XQ4TOrlrd5bbptp/FiCAnuHkIpobH4SHgT+PQUjT7
hXy8eeODXGQvFjjwvv5H38CeumMrgucc4TmQDSCB6iIN0UywQ3D5Oe2zcehOYuvtMsxmOsUkGM/nnjWr
MKXxpQqgfOZo2kNWBya15MfN+jmvohh/HIixvQRfzy7I5f+0Afi2OW7+BAAA//+dNQw+qwMAAA==
`,
	},

	"/defs-config.json": {
		local:   "defs-config.json",
		size:    2236,
		modtime: 1484805416,
		compressed: `
H4sIAAAJbogA/+RVzY7TMBC+5ylGhmNh73vtckMqUgUcEKrcZNydJfaY8QQRob47SrZ089OEVZaeOFRV
Yn8/883Y+ZUBmAJTLhSVOJhbMHfoKFDzlCBaUcqr0goowyZiWHNQSwEF1hwcHWAbMSdHuW3xq0fCM4O5
hUYDwOTt9vMzgNE6YiPI+wfMtYW276NwRFHC1NkNYD4mlN6bDkdSoXAw56Xj6gn37mfkhMUHFk1D/GtB
1+ALdOntQ+Lw6qbj/sbbuG2ZN48eL/OHH1O2rIitzaq7RIp+aGOmEIDjhKhKHZmCDrU54Kap6UtPoC83
a3HG5l+s9u0OLM+6CFVZ9pm6PF8vRrD2xX9b+ycuK4+jcX5m/Uun/tpVfWb5RuFwR4vO+Xu7x/LakWwv
NP6fR5J1/0/sRpjVpZdcoAU5t6NiFNH1b6vT/rmW9mQxVH7UM1PaGiWZ+czOKQh+r0iw6PE8ZTB0l3UZ
/2R+T0lZ6peEngtaxdFNNVe9Y/FW22m0im+UPF6eeFvpPS86KydXu329CM3e4/jL8xwo+qj1ru3kFHzP
XKINZuo8ZM3vmP0OAAD//96B+Ju8CAAA
`,
	},

	"/defs-image.json": {
		local:   "defs-image.json",
		size:    2916,
		modtime: 1485912532,
		compressed: `
H4sIAAAJbogA/8SWz0/cOhDH7/tXzAvocdgf4T09Ib0VQqrKpScO9FS0oMGZbIYmdmo70GWV/71yks3v
hW5B7Qk8tr/+fsYz3mwnAF5ARmhOLSvpLcG7pJAlu5GBFLVlkcWowSq4Skl+VNIiS9LwKcE1wXVKgkMW
WGyflXq1gLcEdwSAl1DA+HmTUh0C8DhwB0bWpmbp+yolKXbyZqH02jciogR9dkf5jcRsJ2BLQc9YzXLd
xFO0lnSBc3vzYf4F58+n8/9Xzb9/HR3/Pb+9W0xX29PZP/+e5f7PLTv2iiPy8iQv4DUZ20bqJdNGBEJv
UqvWGtOIBYiIxFeTJVDuBRWCW6TuH0jYGbAshhUAnJxHaKKL5XlE3zEgwQnG1c6LkwPzUCLdTRfz1XR5
g/PQsU57RAlKDsnYy4pD6Tdc2FBr4LjkbuKavmWsyZ1zU8U61TNrgoafO+PqMlqRNEYbKp14VWjVJEar
lLRlMi28PWW6716LtY5jd4eaQtIkBQXQwyokjjWFbu+R3+qQVlXXa/M+5Ote3DJXO/cbS+ZAPwGFZvFg
lOwaY2nP/hs3NSj7P1f6L2S2cjlKkOnYvOIfIebSpFsMoVYJPEUsIrARm8o1JLiBe4JAPclYYUBB19mu
zFFr3HSn2FLS97C/latZV83oMu9lmr3WXD6KWTdAD/WwPq5VRtHGSmu0i4sZ1CJiS8JmmvpwyrSJVh3F
PQ071OzP/jKt/4LV4T11pvMh2PvZUuatZhaPpE1Z5e9nqhZ9q7mQ0GX8fVPWqO61N2xR2N+mr4K1O3IM
9BE1o+y/oQfn66Vk/Saw8Udo0v67+7hAKZXF7pfhwbfaFql/z5Wkq7Dz2mwP+KlLML0uUK97wK10b8fe
PpnFsTeAXk12o3yST34EAAD//zd996RkCwAA
`,
	},

	"/defs.json": {
		local:   "defs.json",
		size:    3193,
		modtime: 1483690724,
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

	"/image-layout-schema.json": {
		local:   "image-layout-schema.json",
		size:    414,
		modtime: 1484805416,
		compressed: `
H4sIAAAJbogA/2yPwUrEMBCG732KIXq0TQVPue5pQdiD4EU8xHa2zWKTOJkKi/TdJZlWD91TmD/z8c3/
UwGoHlNHLrILXhlQp4j+EDxb55HgONkB4dlew8zw0o04WfWQqfskgwE1Mkej9SUFX0vaBBp0T/bMdfuk
JbsTzvUbkozWIaLvNlkqmGxrl8X1ZxELydeImQ0fF+zWLFKISOwwKQO5TTZkUi5+RUpSS/72bb9lA8IZ
eEQ4HY6wMxdusycm54f/HP08KQNv6wygHpu2adU6v5d3qQCWcjDh1+wI+z/k1rlV5pbqNwAA//8bwMuB
ngEAAA==
`,
	},

	"/image-manifest-schema.json": {
		local:   "image-manifest-schema.json",
		size:    921,
		modtime: 1485912532,
		compressed: `
H4sIAAAJbogA/5ySMW/bMBCFd/2KA+3RMtuik9ZOHooODbIEGRjpKJ1hkQyPDmIE/u8BSTM24wyB16d7
77574lsDIAbk3pMLZI3oQPxzaP5YExQZ9LCZ1YjwVxnSyAH+O+xJU6/S9Cral9xPOKtonUJwnZRbtqbN
6tr6UQ5e6dD++C2ztsg+GoqFOymtQ9OXrZxseVpSJJDziSB7w8FhdNunLfYnzXnr0AdCFh3EwwBEjrhH
z/m4LF/ffDcRgybcDcD5QmQIE0JaDmU55Dx4yYGgGJQBMgFH9Ikihd92maxhP9LKrVdrZjI072fRwa+z
pl6LlqRj/iJ6azSNlw0sPeoYG9HQhLY0Yv06/j9R2XfqgJ4v7YVKea8OFdMm4BxHf577OCnF/N3tAMeK
QhljQ3p5Fcptdbc74iAvI1efqxlQc5tcCWohB9RkKE1XzqawJlLh8XlPHiPWw1fvsP4ndcUNwGNzbN4D
AAD//1HKEXaZAwAA
`,
	},

	"/manifest-list-schema.json": {
		local:   "manifest-list-schema.json",
		size:    873,
		modtime: 1485912532,
		compressed: `
H4sIAAAJbogA/6SSQU/jMBCF7/kVI3ePm3q12lOuywUJxAHEBXEwyaSZqrGNZ4qoUP87sh2XROWA1OvL
vDffm/ijAlAdchvICzmrGlB3Hu1/Z8WQxQDXo9kg3BpLPbLADbHAvceWempNsvyOGb+4HXA00T+I+Ebr
LTtbZ3XtwkZ3wfRS//mns7bKPuqKhRutnUfbltWcbHlaU8TQ44RR74glB8jBY4xwL1tsJ80H5zEIIasG
YkUAlXMeMXCumeXz9g8DMfSEuw4410QGGRASASwIIIfCW04Fw2AskBXcYEgoacMFHfUS+xRZWp/tGsnS
uB9VA3+/NPNetCQd8xdVNvH8HCXZhGAOsw6C43wu/vKAfZzssOc6sa/jP1/pDnuyFA/KpzZX05VdUFPC
ccFirHWS3tOC5pLTzSNPPX4EPXdWhTWRqoCvewoYsZ6+e1hnp60Anqtj9RkAAP//b2/SMmkDAAA=
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},
}
