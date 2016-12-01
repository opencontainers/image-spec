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
		modtime: 1479372487,
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
		modtime: 1479372487,
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
		size:    2483,
		modtime: 1480556605,
		compressed: `
H4sIAAAJbogA/+RWzY7TMBC+5ykiw7GwF8SB6y43UJEi4IBQ5SaT7Syxx4wnQIT23XGypRsnrel2tycO
VRvH389845/+zvJcVeBLRidIVr3J1RXUaLF/8rnTLFi2jeZcKF86sJdkRaMFzsOvGq/zwkGJNZZ6wC/u
CHcMgbDXCIPlMH33HEakc9AL0voGShmgw7hjchB0wY9mh/GPHjgaGXF4YbTXavfqdnGPew+GuJsinzPU
PTKY9S9vPNlnFyPfF2jl9asUX/FTu6fjvHRtsdE8KflRlG9/OfJQfSCWh7Ia7Yoh0OVda/bz2x+HuqGZ
dacW41coYKY2Ev0LUgdEhTtHoeypNllY9jV9iQRiuaTFhM1/WI3tTiwnXdi2aWKmMc/X/UvFVP9t7Z+o
ac18kxxZ/6mr/txVfSb+FmSv8KTj7Z1eQ3PuSIo9jX/ySLLx95ZdMZHU/jH3RoV1vcJqFtH5T6vt/FRL
I1mwrZn1TDW6A/YqndkuBYbvLTJUEc99BlN32Zjxb+Yb9BJfmQ8OvWTQArOTKlV9TWy0DKsxQF8IGti/
4nUrGzppr2xdrdazvwNHockYmN88x0DBOOlWQycPwddEDWirDu2HrP/cZn8CAAD//0zDaxqzCQAA
`,
	},

	"/defs-image.json": {
		local:   "defs-image.json",
		size:    2736,
		modtime: 1480556516,
		compressed: `
H4sIAAAJbogA/7yWy27bOhCG934KQgmQhS86i4MCNYIARbPpKot01cAtJtTImlQiVZJO6gR+95K6UrcE
boyuEg7Jn/83nBH9MmMsiFBzRbkhKYI1C64xJkFupFkOyhDfpaCYkewmR/FZCgMkULEvGWyR3ebIKSYO
xfZFqdcIWD13hA1mGBF83efYhGyQIndgYkyu12EorTyv5fVKqm2oeYIZhOSOCluJRS1gSsFAG0Vi28Zz
MAZVgfP9DpbPm3l499/yo/3v0/Lbar6ZnwfF0kO5I4hoi9r41npJMQkyrva5kVsFeUKcWWf8p95lrNzL
ZMzcInn/gNwsGIliWBlhF5cJ6ORqfZngb4hsyjJIq51XF8fzWJYf89VyM1/bUWxHA6IMBMVW/brikOod
iR9qDRyX3G1c4a8dKXTn3FWxThUs2qCm5864ugwvkqdgYqmyoApt2sQoa97WKGoPb6Lcpu61WOs46jtU
GKOyKcGI9bAKiXM77/aeeYXuFWez9NBnfNuKW+ZK535vUB9px7rRqwctxVnoGyNhPvw/bmpQ9dN41dJR
mZ1K9RtswFIqe8QtZrGSGXuyXZRYQNIVFstgz+6RRfJJpNJ2SdQFrUsNlIJ9d4oMZn0P0+1UzbqKAodv
CSjw5g6jmE0R9lCP66VGZRRt7H5HO6mYAcUTS87NTmEfTmqfaNNRnGiaoWZ/9q9pw1esDu+pM30Ygp3O
lhV7p5nVoz2hrPLTmWpE32suRnAZP23KWtVJe8MWLY8cbdM3wfyOHAN9BEUg+h+yo/P1WrL+Edj4R2jm
/60feBBCGuj+yjr6Vn2R5k2VAm/iztdm9H2YeG8yyG8L1NsesJful7Fvn9ilaTCA3szq0WF2mP0JAAD/
/0efdTKwCgAA
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
		size:    1139,
		modtime: 1480556516,
		compressed: `
H4sIAAAJbogA/6RSO4/UMBDu8ytGvuu4xByiupbqCkTBiQZRmHiSzCl+YPtWrFb73/Ej3mx2KRBb5ou/
18wcGgAm0feObCCj2ROwLxb1J6ODII0OnpUYET4LTQP6AF8t9jRQL/Lrh0S/9/2ESiTqFIJ94vzVG90W
tDNu5NKJIbTvP/KC3RUeyUrxkWOia19dfaaV15xSAq6WBIUb9hYT2/x8xX7BrIsaLhD6+CcVi1iR+BYV
S7kCX3d+mcjDQDhL8KUheggTQjaHag5FD3ZFEIQHoYF0wBFdTpHF/68Z34Y9qdWuVzaKNKk3Ff99WDHx
u2IZOpY/TKEk8VKUDjfmXLWuMvrgSI8rjjpn+b58R0RYOy/nw3dadqanLut3Vb/bPb5LF8QW0o9NkZhw
oPG8xb3DIXmn7KhDW1drXLfKVPos9rHcOb1GF86J/Wa4zwFVevq4LnZBDmuff3KP/psUQmsT8gz87fto
Z4pLOZd8uByNxMG3Zcop1B2PQDye/HrDbGrWnJQ5/PVGDuVphexvN8ouD6IuaTvzJq3y2PwJAAD//xOl
pJlzBAAA
`,
	},

	"/manifest-list-schema.json": {
		local:   "manifest-list-schema.json",
		size:    1101,
		modtime: 1479372487,
		compressed: `
H4sIAAAJbogA/6ySP48TMRDF+3yKkY8OsgZEdS00SCAKTjSIwqxnd+cU/8Hjizid8t3xn/Vml1AgXZpI
eet57/fsedoBCI3cB/KRnBW3IL54tO+djYosBvho1IjwWVkakCN8ovTz1WNPA/WqjLzKHi+4n9CoPD/F
6G+lvGdn91XtXBilDmqI+9fvZNVu6hzpNsJpxqXovkVzGaunJWUMaWaM/SFhVIP46DFbuJ/32M+aD8ko
REJOX3LFpFWfb8m21qzyZfu7iRgGwoMGrjWRIU4IhQA2BFBN4VhdQTEoC2QjjhgKSkl4Rke5xV4sW+uL
LEOWzINJ396eNfW7aUU61S/CoCZ1V52ergF7NrwA5RjIjmcdbQH6Pv9PivL+MK+UPFrduZ66EtK1kC6H
dMc3L/NqiXnyx7bSfJTXlRqDCkE9rp4lolmfy1sccMgnNQ68r+E560YmIV1sZuOl84d5cVxoLKcNi7LW
xdKHr3XBa8ulx39Bryd3jbWQioC/HiigXp5D/GvpxN+Pu7rrXX6G0+5PAAAA//+jTbrrTQQAAA==
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},
}
