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
		modtime: 1462375266,
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
		modtime: 1460709579,
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

	"/doc.go": {
		local:   "doc.go",
		size:    711,
		modtime: 1463675171,
		compressed: `
H4sIAAAJbogA/2SSMW/bMBCFd/+KB08t4Epphg7tpDoJKjSQC8tpkJGmTvKhEsmSVGT/+x5lB4hRwwt5
795998Q8x9q6k+fuEHF78/kLdgfCI5vxiAc7mkZFtmaR5/KXa00mUAO5J48oysIpPTfMlRV+kw/SgNvs
Bh+SYHkpLT9+SxYnO2JQJxgbMQYSDw5ouSfQUZOLYANtB9ezMpowcTzMcy4uWfJ4uXjYfVQiV9Lg5NS+
F0LFC3T6HWJ0X/N8mqZMzcCZ9V3en6UhfyzX91V9/0mgL01PpqcQ4OnvyF4W3p+gnEBptRfUXk2wHqrz
JLVoE/TkObLpVgi2jZPylGwaDtHzfoxXmb0hyubvBZKaMlgWNcp6ie9FXdarZPJc7n5snnZ4LrbbotqV
9zU2W6w31V25KzeVnB5QVC/4WVZ3K5AkJnPo6HzaQDA5pUnNHF1NdIXQ2jNScKS5ZS2rmW5UHaGzr+SN
bARHfuCQvmoQwCbZ9DxwnF9G+H+vbJEkv5T+k3yCpD0oNNSyoTDrNutSoFJxoIYV4slRWF0pOb6Nw6vq
+fwK0Y5Gz4Vs4a7sF/8CAAD//2ICNX/HAgAA
`,
	},

	"/gen.go": {
		local:   "gen.go",
		size:    839,
		modtime: 1463676428,
		compressed: `
H4sIAAAJbogA/2SSwW7bPBCE736KhU7JD1vKn0MPKXJQk7gVGsiApTTILZS0ollTJEtSlvX2XTIKkKCG
LxR3Zr8dbpbBnTazFfzg4frq/y9QHxAehRrPsNWj6pgXWq2yjP70uUXlsAP6jhY8VeaGtVEQb9bwC60j
AVynV3ARCpLlKrn8GixmPcLAZlDaw+iQPISDXkgEPLdoPAgFrR6MFEy1CJPwh9hncUmDx8vioRvPqJyR
wNCp/1gIzC/Q4Xfw3txk2TRNKYvAqbY8k2+lLnss7h7K6mFD0IvoSUl0Diz+GYWlgZsZmCGoljWEKtkE
2gLjFunO6wA9WeGF4mtwuvcTsxhsOuG8Fc3oP2X2jkiTfyyg1JiCJK+gqBL4lldFtQ4mz0X9Y/dUw3O+
3+dlXTxUsNvD3a68L+piV9JpC3n5Aj+L8n4NSIlRHzwbGyYgTBHSxC5GVyF+Quj1G5Iz2IpetDSa4iPj
CFyf0CqaCAzaQbjwqo4Au2AjxSB83Az371zpakUZH4OJo6gHtgqS76jQMo/BA3BoGuw6Gjq8TLql969m
53GIQEzKRRk3wwX56AIKuhYugsTRa3KadGxS2pZs+C0ap1VG95dpbFeHvXIHPcoOeGwt5QwNha5O+kiN
42a9Duz4Trnp3Sspub7hC2pstzFWnMJhY478dsHaCK60xdsk/Y/rBNLV3wAAAP//Xhj9JUcDAAA=
`,
	},

	"/image-manifest-schema.json": {
		local:   "image-manifest-schema.json",
		size:    1064,
		modtime: 1462375266,
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
		modtime: 1462375266,
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

	"/schema.go": {
		local:   "schema.go",
		size:    1719,
		modtime: 1463675171,
		compressed: `
H4sIAAAJbogA/6RUUW/bNhB+jn7FQU8JZlNtH/aQIQ+eG2PaUhuI3BbFMKy0dJK5SaRGUnbUof99d5Ts
OnFQDK3gF5N333333cdLEpibtreq2np49eLlj7DeItwp3T3AwnS6kF4ZHSUJ/eg4R+2wADpHC54iZ63M
Q0K4mcA7tI4S4JV4AZccEI9X8dVPDNGbDhrZgzYeOoeEoRyUqkbAhxxbD0pDbpq2VlLnCHvlt6HOiCIY
48OIYTZeUrikhJb+laeBIP1Imr+t9+11kuz3eyEDYWFsldRDqEvu0vntMrudEukx6a2u0Tmw+E+nLDW8
6UG2RCqXG6Jayz0YC7KySHfeMOm9VV7pagLOlH4vLTJMoZy3atP5R5odKFLnpwGkmtQQzzJIsxh+nmVp
NmGQ9+n6l9XbNbyf3d/Pluv0NoPVPcxXy9fpOl0t6d8CZssP8Fu6fD0BJMWoDj60ljsgmorVxCJIlyE+
olCagZJrMVelyqk1XXWyQqjMDq2mjqBF2yjHU3VEsGCYWjXKB2e4875EFJHGfzOII6kbGUVEwVgPsUaf
8CziiFHeYKEk+L5FdySymqdEmHPppJHeRTkV8XAZXYTwNUW/kVqVSIfn3ztZK7IsgfF3Ax/HqTHVZKcL
YXIlAr5oRhSxe/nDX87oj89UuFNnVb6lQq2eL5NyWIZWEeangHAs02keGzao2RpfLeNOAYQ1xpdOeGlF
9Um1Xy83N7pU1feUywPC/2tubpqN0uEpfHO1AeFLvaso2knL9iBDlQ4cTYbsxFbCZoNFQQXYcWJBOybr
nccmhG7ljr19sNyv2Wo5ujVsI8dPmqWEOIlFdEHIN/AnunyRXZaydkh1GYbfjaNl1LoAMyI0J8am7XCK
S1hDzg1n/X500x+8B3T1b3Rx7sLr4Ig4CDI9mGo6oAqWIZ48l8bmvYb4mMAufJr1mQWkPr6oQzvPdzY8
ddJg2mBjbD9saDcEKJ3XXXFQ73FzBLV+cgi0C6E2NF0ahfQhKShb0GbNqfVeRGWn8xMOl1dPZwYkzECM
hhx9jv4LAAD//6aSebu3BgAA
`,
	},

	"/spec_test.go": {
		local:   "spec_test.go",
		size:    4209,
		modtime: 1463675171,
		compressed: `
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
ppaZyd9DMuLC/K8AAAD//0k/ScFxEAAA
`,
	},

	"/validator.go": {
		local:   "validator.go",
		size:    1894,
		modtime: 1463675171,
		compressed: `
H4sIAAAJbogA/4xU0W7bNhR9tr7iTkABCdOkrA97yJAHL00wb4EDRGmDoigGWr6SuUikRlJ2jKH/vkNK
ju12Dw0C2xLvPfecw0MWBV3rfm9ks3H09uLnX+hxw3Qn1fBCt3pQa+GkVlFR4B+vK1aW14T3bMihct6L
KjSElYw+sLFooLf5BSW+IJ6W4vRXD7HXA3ViT0o7GiwDQ1qqZcvELxX3jqSiSnd9K4WqmHbSbcKcCSX3
GB8nDL1yAuUCDT2e6tNCEm4i7f82zvWXRbHb7XIRCOfaNEU7ltribnF9syxvfgLpqem9atlaMvzPIA0E
r/YkepCqxApUW7EjbUg0hrHmtCe9M9JJ1WRkde12wrCHWUvrjFwN7syzA0UoPy2Aa0JRPC9pUcb027xc
lJkHeVo8/n7//pGe5g8P8+Xj4qak+we6vl++Wzwu7pd4uqX58iP9uVi+y4jhGObwS2+8AtCU3k1eB+tK
5jMKtR4p2Z4rWcsK0lQziIap0Vs2CoqoZ9NJ63fVguDaw7Syky4kw36rK48iePzsQSys7kQUgYI2jpJo
Ftedi/El9fhZSD042cYRnhpQH1Y5tr/on5uCjdHGxucLLyz7YdgWjf7bYnrAj6M08qw+iFYir5C0M6IH
Wep4LQW5fQ8q8Bli5JqVg1I2vgNygj3c4a2l7Qjg8ysaRMs6YPxR3i8nIXkUoI5zRtDT4ei98byRSeXT
CRZtG9wZ5eCncLRBmFhhz9dDYHUcfD7iFQ2DhsrRv9EMz5Y+fQ5o0ZcoqgdVUfJNR0rhK0kPwtFq2A1G
EXYgL3u8dHUSv9nGyEzuUVMPd1TCB1ZsA/9GblnhRAi/3Qd7QnaCN4fz563voexofT5x3B6NS1+HJNZU
JHX+EHDT0SVPdjXUmX+iyysaMxJq5m3rW9JoJuuw/MMVKdn6joO+0ej8CUQS/M4oHlQ4tjionj7uDiVr
ti4GCiTPbOuHnCYqX/LugWs2jEvoTntqt7ijyr113CWxv69wnRTxj/7c2E/bzxnV8G/W/S9SGXZghEnG
7UggL0VqwdkOrXtVetZ69KjNqGu/U3M9iZ625Y29pKP+w5YePchoO9kA8JHNOBjRORmAeWMV+/yBaSee
OZlymNFFRi2rZOoP0bNJCoUzf8P8ldGabeXbDO4Xpq/qwqAAfEXhZKy9BpuFqIYin1QbjzDpxHdi9lXw
PZRP8yUFCF+K4nCoBvV61hHQw9k9hPNs9TsCenKajhQvz2FGd79E/wUAAP//3Fs/YGYHAAA=
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},
}
