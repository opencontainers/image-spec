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

	"/.content-descriptor.json.un~": {
		local:   ".content-descriptor.json.un~",
		size:    1548,
		modtime: 1485467092,
		compressed: `
H4sIAAAJbogA/wrLzJ0fmueS/5SBSaMjRGHmQ/EFfl9eT6850djLIrZo9xKlhtQ9qr06a//Zz2RgYFBl
YGAIUAADpcwUJSsFpYySkoJiK339/ILUvOT8vJLEzLzUomK9/KJ0/eLkjNTcRP3M3MT0VP3cxLzMtNTi
Et2czOIS/cS8vPySxJLM/LxiJR0GBgYZBgYGUwYGBkYGBgZmBghgRsIMEV25F1lAskwM8RegCpgYEAAk
JQXj/P///z9YZBTQCER05R5jgYTwVwlo0MszMDAooohIw0RMG8GRhhy7DNAIBKljZRiNNDoAeKQxQqJI
EhpFqtCoYQVl61poXDGhaWZGyqWjcUUpgIcklO0IZZfB2BFduReg5R04rkDxJIsUV1QvhE0bn68CBAAA
///GLb+iDAYAAA==
`,
	},

	"/.defs-image.json.un~": {
		local:   ".defs-image.json.un~",
		size:    29232,
		modtime: 1485787724,
		compressed: `
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
`,
	},

	"/.manifest-list-schema.json.un~": {
		local:   ".manifest-list-schema.json.un~",
		size:    23207,
		modtime: 1485547022,
		compressed: `
H4sIAAAJbogA/+zcS28bVRTA8VM7fSdNk6ZxmjSt46RJ8/CDR3iEV4Hw6AJ2FLNCVTNBrogdxQYBRUis
kNo1S2CBVCQk+AKwYMcWsesexJIVsEAo6Lj3xie+M5NjnwOr85cqTSfOnFF/Gie2b+dabeuL1+rrjV8h
c/W7ya9e+eejv8d+vpn98l79yo+fvfPNvVcrf/3y7Sd/fv7p73UAmAWAZ/OuwtxOtFlYyxc2os1msbZ1
/a2odLPZqM+WN6LNWr3WqjXqzfLW9XptM2q21qPmjZ3adquxUwCAHACUAOAQAFTgfhXyB6p37p4awK8O
wJs/uQdkoBN+6Yz/y+7u7m57j9Vfp8j2GbJ9zW9Ub29PD9z/F/5jwj1oEgCm9+0Z83tWP26j4TdkyeEy
7jGHwdDkcdDybbRDHaKz7irG3aPtq9hfoM3CWv6Wg8u0L7xOWXfJjoDByePAvU6vtnEAmCBw8qdg55z1
F6NrwJwVEztfCZz7Y0bWI+QEDhuzYmLmZwLmvpRR9SiZf8SUFeMoV9OUnw6U+0FG1GNk/FFDVkyM/FSA
3Icxmh4n04+ZsWJi4ycD496JkfQEGX7ciBUTEz8REPcsjKInyewTJqyYWHgtEO4VGEEHyeiTBqyYGPjx
ALhHX/QcIpMHzVcxju8bab6Pha+Ne+Id7DqJIeNVjMO7kcb7aMDbky5qDnedD+ougunK4+i+l6Z7NdC9
0ai3onqruLFH2ZM3ntFpcibD5q0Yx/v9NO+X2d487mH/RO06bdyKiblfYnOztFF3lJzHiGkrJtZ+ka3N
wR7pOo1Rw1ZMjP0CG5thPerWA9AzMmutxNbrbOuDqf3aA9+YUSsmpn6ef1kfJD3mDu87a9KKiaWfY0sf
BO1hfeMGrZgYOlwklHhJpzv7A/ty5qyY2DlcJJTknM6MpOfICUwYs2Ji5nCRUOIvY2nKE27hru+cKSsm
Vg4XCSUppyEj6hQZP2nIiomRw0VCie9+Jxuj6XkyfcqMFRMbh4uEkoyTiaf8f6dwnTdixcTE4SKhJOJE
YRS9QGZPm7BiYuFwkVDyG2DxwAh6kYy+YMCKiYHDRUKJr5HjfdEzTyZfdL6XwXzlcXw/6G2RUOJHU7G8
yDlDBueNVzExb7hIKIE3Xhc1C2TujOkqJtZ9hKsbizvjjuQrGK5iYtxVLm6cLVrOkamzZqsYx/ZWmu3D
XNsYWjzGJTJ0zmgVE9M+xKUNZVFynsy8ZLKKiWUf5P+07YJFyAUyct5gFRPDPsD+SdvlOu8NXQvmqpjY
tcJ17WJd8G8/uS4bq2Ji1jKXdb8q6i2ReYumqphYtcRV3YeKiMtk3JKhKsZB/TANtchFpaZouEKmLZup
YmLTFfZr1w4pEhbJsBUjVUxMuswl7YiuuNtT+oomqhhfNBMvusT+HcmDImCZjCrRtcsGKiyXsE1Bv4+7
RAt7e3LBnnFyMPzWRS66Ny+Re8uC8/f3nTVzaTnyFInbq277Xb9dvb39Q9s8+x/fxLLcdWYVu5PD/1n1
zt0hd6dm5vr0VOfo7Wgr6vpk4Lev/w0AAP//DEZu9qdaAAA=
`,
	},

	"/config-schema.json": {
		local:   "config-schema.json",
		size:    774,
		modtime: 1485787737,
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
		size:    956,
		modtime: 1485797995,
		compressed: `
H4sIAAAJbogA/5STPW/bMBCGd/2KAxMgSxx1CDoIQZZ279BuRQeaPIqXmh89npGqhf97QVGK7bZI4MWw
X7wP7zmZ+t0BKIvFMGWhFNUA6lPG+CFF0RSRoX7DKPBxKSWGzxkNOTJ6Jm7rEdfFeAy64l4kD33/VFLc
tPQu8dhb1k427+77ll01juyKlKHvU8Zo1sllxlq7ty/TGydTxkqm7RMaaVnmlJGFsKgB6mIAKqAl/aWV
W/TvvuIR5l49FJKDGjA6ZIwGLZzMmPFrRlc5i65sKOgR7+qyV71FR5HqsaU/Dp6xQ6NVoV9vqdQKUITt
JFgu1PmPCEV5f38uYWnEIm9oGJ6ypJF19mTAeDTfyz5AY1et5nJbdevPrEWQI9w8eF384/Dg8ae2aCjo
3UI+3lz4IBfZswX2vCuv6GvYUXOsRXCcAjx7Mh7EU1mkIegJtgg2Pcdd0hbtUWy9XZpZT8eYBMPp3JNm
EaY4vlQBlEscdH3Ias+klvzQrZ/zKorxx54Y60vw9e8Le3pjzv+4DuBbd+j+BAAA///PVSBUvAMAAA==
`,
	},

	"/content-descriptor.json~": {
		local:   "content-descriptor.json~",
		size:    1122,
		modtime: 1485467079,
		compressed: `
H4sIAAAJbogA/5STz24TMRCH73mKkVupl6bmUHGIql7gzgFuiINjj+Mp8R/GE5WA8u7I691mFxAlt93R
7/N8M7J/rgCUw2qZilBOagPqQ8H0LicxlJChfWESeD+GMsPHgpY8WTMQt+2I62oDRtPwIFI2Wj/VnNa9
epd5px0bL+s397rXrjpHbkLqRutcMNmpcx2wntbupXvn5FiwkXn7hFZ6rXAuyEJY1QbaYAAqoiPzqYd7
6c95JSAMuXYoZA+twOiRMVl0MOsx4NeMvnEOfV1TNDu8a8NeaYeeErVjqz43HrBTp1WlH6+ptAhQgu1R
sF6o8xcRSvL2finhaIdVXtGwfCySd2xKIAs2oP1aDxE6O2l1l9um236LEUFOcPMQTA2Pm4eA341DS9Hs
R/Lx5sJFjrKLAQ68r//QN7Cn7tiC4DlHeA5kA0igOkpDNEfYIrj8nPbZOHRnsel2GWZzPJdJMM77zpJV
mNLuJQqgfOZo2pLVgUmN9dNiDJNSluERLab53zcxbExHk8hjlXWbWc+PvGzPc3I1uQ6mivHbgRib1uff
n9X8Xi+v1wrgy+q0+hUAAP//qPU6nmIEAAA=
`,
	},

	"/defs-config.json": {
		local:   "defs-config.json",
		size:    2236,
		modtime: 1485787737,
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
		modtime: 1485797995,
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

	"/defs-image.json~": {
		local:   "defs-image.json~",
		size:    3185,
		modtime: 1485549194,
		compressed: `
H4sIAAAJbogA/8RWTW/bOBC9+1fMKsHm4A9lF4sAawQBiubSUw7pqYETTKiRNalEqiSd1DH03wtKsvUd
143R3swh+ea9NzO0NiMALyAjNKeWlfTm4F1TyJLdykCK2rJYxajBKrhJSX5U0iJL0vApwSXBbUqCQxaY
X58UeDsAbw4uBYCXUMD4eZ3SLgTgceASRtamZu77KiUptvBmpvTSNyKiBH12qfwKYrIFsAWgZ6xmuazi
KVpLOpdzf/dh+gWnr+fT/xfVz79OTv+e3j/MxovN+eSffy8y/+eOnXp5iqzI5AW8JGPrklpm2ohA6HVq
1VJjGrEAEZH4alYJFHdBheAOqccnEnYCLPNlKQDOLiM00dX8MqLvGJDgBOPy5tXZgT4Ukh7Gs+liPL/D
aei0jluKEpQckrHXpQ6l31GwLlaHcaG7imv6tmJNLs9dGWt0z6QKGn5trMtilIFFZYJWKWnLZGpSBlpy
qIb5Wcd5Wy9NIWmSggJoScghTjWF7u6JX5uGWgfvzmZtQfu5uGOuTx7XlsyBfAIKzezJKNkkxtJe/NdP
qtPif67N33C2WfumgpWOzR7+CDEXJN1hCLVK4CViEYGN2JSsIcE1PBIE6kXGCgMKmsy2LY1a47q5xZaS
NofhsS13Q6UTdM57K81ebS/rlZnGaN2VttTDZnaH0iutr7V6JzbfQS0itiTsSlNbnDJ1RYsG4sDAdjHb
u7+s1n+DardOje2sK+x4tJR5L5nZM2lTdPnxSO1A30suJHSOH9eyCnWQXndEYXhM9wqrT2Sf0GfUjLL9
hh7s11tm/SZhex8hlFJZbH77bTkcVM86UONxUJJuws5b0zFkzx9egultLvh2v/GDZstVHHuDBi1G7WjW
+NLqt+oINvVYtDngW2DQmpotm74/h6YdW9GFDW6VjbLRjwAAAP//FgEL5nEMAAA=
`,
	},

	"/defs.json": {
		local:   "defs.json",
		size:    3193,
		modtime: 1485787737,
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
		modtime: 1485787737,
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
		modtime: 1485787737,
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
		modtime: 1485797995,
		compressed: `
H4sIAAAJbogA/6SSQU/jMBCF7/kVI3ePm3q12lOuywUJxAHEBXEwyaSZqrGNZ4qoUP87sh2XROWA1OvL
vDffm/ijAlAdchvICzmrGlB3Hu1/Z8WQxQDXo9kg3BpLPbLADbHAvceWempNsvyOGb+4HXA00T+I+Ebr
LTtbZ3XtwkZ3wfRS//mns7bKPuqKhRutnUfbltWcbHlaU8TQ44RR74glB8jBY4xwL1tsJ80H5zEIIasG
YkUAlXMeMXCumeXz9g8DMfSEuw4410QGGRASASwIIIfCW04Fw2AskBXcYEgoacMFHfUS+xRZWp/tGsnS
uB9VA3+/NPNetCQd8xdVNvH8HCXZhGAOsw6C43wu/vKAfZzssOc6sa/jP1/pDnuyFA/KpzZX05VdUFPC
ccFirHWS3tOC5pLTzSNPPX4EPXdWhTWRqoCvewoYsZ6+e1hnp60Anqtj9RkAAP//b2/SMmkDAAA=
`,
	},

	"/manifest-list-schema.json~": {
		local:   "manifest-list-schema.json~",
		size:    1011,
		modtime: 1485533214,
		compressed: `
H4sIAAAJbogA/7SSz07zMBDE73mKlfsdv9QIccoVLkggDiAuiINJNs1W9R+8LqJCfXdkOy6JigRSxXXi
mf3Nbj4qANEht55cIGtEA+LOobm0Jigy6OFaqxXCrTLUIwe4IQ5w77ClnlqVLP9jxj9uB9Qq+ocQXCPl
mq2ps7q0fiU7r/pQn13IrC2yj7pi4UZK69C0ZTQnW34tKWJIPWLUG+KQA8LOYYywL2tsR81569AHQhYN
xIoAIuc8oudcM8vH7R8GYugJNx1wrokMYUBIBDAjgBwKbzkVFIMyQCbgCn1CSRNO6Cjn2IfI0vpoliZD
eqtFA+dfmnovWpL2+Ysok3i6jpKsvFe7SYeAevountxjH1922HOd2Jfx5gvZYU+G4kL50OZq3LL1YkzY
z1hwgxrNn6KMI34iUcbYkP7sGcwpR5xGHmr8innqrAprIhUeX7fkMWI9ffeLHx25Aniu9tVnAAAA//+b
GIGn8wMAAA==
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},
}
