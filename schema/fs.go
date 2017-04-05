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
		size:    2771,
		modtime: 1491372287,
		compressed: `
H4sIAAAJbogA/+RWQW/bPAy9+1cYbo9t/R2+U67dbgMyINh2KIZAsemEnSVqFD3MGPLfB8vJZtmym3XI
aScDFB/f4xMl60eSplkJrmC0gmSyVZqtLZhHMqLQAKePZCrcpxsLBVZYKJ9118FuXXEArTrIQcSu8vzZ
kbnvow/E+7xkVcn9f//nfeymx2F5hrhVnpMFU5zZnIf12TlqtYe88Pw9UloLHZZ2z1BIH7NMFlgQXLZK
u3bSNCsYlED5KzCAOmE0fTkfr4i1km6lVAL3ghoyv3bsUzLVyIF4oVSYzcUBBQppGC7FkLs08+RFJHvg
iI9HXPHxDw44iMwwDlh9ztvvlhyU74nFjfG3DJU3ECr30I3ATV5ChQa7UXG5VnbjK697jfH65tucLMWs
2uxuuIQCeixjoZE0Pc6QCreW0MiYmwysu56eAoKQblHigswXpIZyR5IXVZimrsNKwzqfoxY86vKf7f0j
1Y2GyThf2P9rp/7aXX0i/oJm/wZfdc7fqR3U17ZkE9n4a1qyEbIb3BtVX2xJMvyer18mkip6WV96/ZZY
VVssJwZf/6475S91H9CCafRkx7NatcAuizuejFgzhq8Nsv8PP0U8GKtLhhXPnh/QCXEbMz00K2LU3PbM
b1D07fCyW0vviMlWxN4UcYpZ/Enjdtf+RQ3SGiZ/vj8oANpKu/UTMV9kR1SDMjPzGZ6y5MQwnZvwWfX7
2RSey6SbnWPyMwAA//9KY9sL0woAAA==
`,
	},

	"/content-descriptor.json": {
		local:   "content-descriptor.json",
		size:    1091,
		modtime: 1491372287,
		compressed: `
H4sIAAAJbogA/5yTwW7UMBCG73mKUVqpl27NoeIQVb3AnQPcEAevPY6nJLYZz6oE1HdHjrNsAoiFve2O
/m/mm2j8vQFoLWbDlIRiaDto3yUMb2IQTQEZyi8MAm+XUGR4n9CQI6Nn4ra0uM7G46gL7kVSp9RTjmFX
q3eRe2VZO9m9ule1dlU5skckd0rFhMEcJ+cZq2llf06vnEwJCxn3T2ik1hLHhCyEue2gLAbQjmhJf6jh
Wvp9X/EIc640heigFBgdMgaDFlYzZvya0RXOosu7k9hd2fhKWXQUqPTO6jR9Zl9qizbTt3M+JQIUYD8J
5v90+oMIBXl9v5Ww1GOWMxqGpySxZ508GTAezed8GKGyR63qclt0y9+kRZAD3Dx4nf1j9+Dxq7ZoaNTD
Qj7eXPI1F+PNFgce8l920DBQFS1BcBxHePZkPIinvJjDqCfYI9j4HIaoLdpL7GaTjZsOIcr8RjaK/3ry
NOoeV4ev1v0uEFzj1bNZXFvGLwdiLGIff30365vdnk4D8Kl5aX4EAAD//4LEuuxDBAAA
`,
	},

	"/defs-descriptor.json": {
		local:   "defs-descriptor.json",
		size:    1966,
		modtime: 1491372287,
		compressed: `
H4sIAAAJbogA/7SVT1PbMBDF73yKrWDKIX9Me+hMMwwznXLvgZ7KBGaR19HSWFJXMjRl/N07ihPbCYSh
NNzitfT2/Z680cMBgMopaGEf2Vk1AXVOBVtOTwE8SmRdzVEgOvjmyX51NiJbEjhfbXMCF540F6xxqTFs
RFsVNYHUB0CVlDN+X3hqSwCK89TVxOjDJMucJ6vXPcLYySwL2lCJGZc4oyxvu2ad2nCtFRttFaKwnXV1
jzGSLPGuLr+MfuDoz8no87T7+e7w6P3o6no8mD6cDD98/FRnL1t2pJYt6qaTynlGIfbptsKNhkDLwkc3
E/SGNWhD+meoSmj2gisgLXI3t6TjENguH1cAcHxqMJizyamh35iT5hLnq51nx/+YQ4N0PRiPpoPJJY6K
xDrYIqpkHp7hQZhzYzothEJcCfeGtYFoOKwooMQF3BDk7t7OHeaUP3aKIrjoyhyp7PfdzQSgCiclptRV
JaxW9XoDw88xpmV7+fBasUcUDW9XF/pVsVBqdNn5RdGGI+lYCfU5XFibn3YnJs6TRKatNDY0+m/2AZbt
cPj4FNpX9SbH3h258Bof4zuS0Hyqe/fTar/GV0GYsn2ToDrxJ51tTho8PW3PsqyHa5vrDoXRxm2ml0Sy
K483Mv3kfwRa6yJuXln/cyR9vXaej4SKpJdTEca3wdnDrHdVZiX6i6Xxi579ZLM+qA/+BgAA//+NODnR
rgcAAA==
`,
	},

	"/defs.json": {
		local:   "defs.json",
		size:    1670,
		modtime: 1491372287,
		compressed: `
H4sIAAAJbogA/7STza6bMBCF9zzFyO2S9oJtbGDb7hMpy6oLSiaJq2AjY6RWEe9e8RNChFuJKneRgGc8
3zmeMbcAgByxKa2qnTKa5EC+4klp1a8aaBs8grtY054vpnXgLgi7GvUXo12hNFo41FiqkyqLoTwceTOA
5NBLABClXTqvAIj7XWOvprTDM9qhckhUSquqrUgOn2KaPsLFrykcUzkEu3Amx2IrmlEpfPA+vsIzuhVP
Yy55ygT3aczJlZDgW4UyShmTNGIiTbiUIooij6Jn15N0+x/T8enQJFlxN8/GBxZJwtbozXPxoTnNeCYk
zdb8zePw8eOUcyE5jySTUZYk1Nf8WOxNz7VLQaNxdyI5fJsCMKeG9EeLfZZ8eFt8cG9Ty+eNXeivvp9G
t9frYvf09t3Ti1c6FPy1DhtnlT5vd3jXGOtf66kq6sOAHf99V8n8+Imle9ykunAOrd5bU6N1CptFEQD5
fIvD7in0ryMEy+fK1G6UfmdTE+tvpoL+1wV/AgAA//96IpqyhgYAAA==
`,
	},

	"/image-index-schema.json": {
		local:   "image-index-schema.json",
		size:    2107,
		modtime: 1491372287,
		compressed: `
H4sIAAAJbogA/6SVP2/bMBDFd3+KgxIgSxIWRdBBCLK0S6YODboUGWjyZF4q/ilJN3YLf/eCpCXLlgwk
6pac7r37Pfok/l0AVBKD8OQiWVPVUH11aD5bEzkZ9PCo+Qrh0UjcwDeHghoSPLdeJ+1lEAo1TzoVo6sZ
ewnW3JTqrfUrJj1v4s2HO1ZqF0VHspOEmjHr0IhuZMiy0s0ojWeUxhdh3DpMUrt8QRFLzXnr0EfCUNWQ
IgFURf8dfSixSnmc9klRgIawlRBKPAwQFUKeDHkyFDP4XdyAB+AGyERcoc8I2XlGJnaM2Vt1KUczNBnS
a13V8PFQ45uulku78qTS3FCDIYZh/M6Ze8+3A/aIetj37jjdsC/747UH6Olfbf/E4681eUyzfvTVBI+S
+FOSXQ/Lgf6cVCStMMSqLz0PzCcWY2R//GC8IGkXcnfKALbJy+GxQY9GoIRRpGxy6bFJaolNuJH9mdym
t+OCSWzIUJoQ2IFk4LAbh34LZ2oEMrDcRgyzWCcAycRPd+fh9uf/Fjzhty7aledOkQChUPwMaw3FocMt
jNcpRvrX8RjRG7i6Vzyoh/pe4YZLFKR5u1c+XM0//dPtOU239u3p7oyzcWipBEjt0Hir4VWRUBDT16Uk
As23sESQ9tW0lkuU86kz1Vlm1/LYWK/H3O97pQ8IrLecjdw7nMXmxtiYr5aJE59NPnSdDT80GfIvTv/a
HX2ApxPNuiemYvxPgASaMccf4GrqThpcJguA58Vu8S8AAP//NlOuVTsIAAA=
`,
	},

	"/image-layout-schema.json": {
		local:   "image-layout-schema.json",
		size:    439,
		modtime: 1491372287,
		compressed: `
H4sIAAAJbogA/2yPQUvEMBCF7/0VQ/Sg4DYVPOW6pwVhD4IX8VDTaTvLNonJVFik/12SaRXRU5g38+W9
91kBqA6TjRSYvFMG1DGg23vHLTmMcJjaAeGxvfiZ4cmOOLXqLlPXSQYDamQORutT8m4nau3joLvY9rxr
HrRoV8JRtyHJaO0DOruZpYLJtaZsrM/FWEi+BMysfzuhXbUQfcDIhEkZyG2yQyYl8TPGJLVk97fth1yA
74FHhOP+8LvyDbmy8JZ2EgZ6OuNtsS8fbrESR3LDj45unpSBl3UGUPd1UzdqnV/Lu1QAS2kS8X2miN03
8l+PKnNL9RUAAP//k31n5bcBAAA=
`,
	},

	"/image-manifest-schema.json": {
		local:   "image-manifest-schema.json",
		size:    921,
		modtime: 1491372287,
		compressed: `
H4sIAAAJbogA/5ySMW8iMRCF+/0VI0MJ+O501bZXUZxSJEoTpXB2x7uDWNsZmygo4r9HtnHAkCKifTvv
zTdv/dEAiB59x+QCWSNaEHcOzT9rgiKDDOtJDQj/lSGNPsC9w440dSpNL6J97rsRJxWtYwiulXLjrVlm
dWV5kD0rHZa//sqszbKP+mLxrZTWoenKVp9seVpSJJDTkSB7w95hdNuXDXZHzbF1yIHQixbiYQAiRzwi
+3xclq9vfhjJgybc9uDzheghjAhpOZTlkPPgLQeC8qAMkAk4ICeKFH7bZbKG/Uort16tmcjQtJtEC39O
mnovWpIO+YvorNE0nDcwZ9QxNqKhCcvSiOVV/H+ism/VHtmf2wuVYlb7imkdcIqjv099HJVi/ul2gENF
oYyxIb28CuXGus/TFpet9Kj9JdRM9qjJULJU9qawJlLB+Lojxoj19N07rP9JXXED8Nwcms8AAAD//7u3
Dj+ZAwAA
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},
}
