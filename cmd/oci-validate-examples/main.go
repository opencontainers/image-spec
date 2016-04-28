package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
	"github.com/xeipuuv/gojsonschema"
)

var (
	verbose    bool
	skipped    bool
	schemaPath string
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "display examples no matter what")
	flag.BoolVar(&skipped, "skipped", false, "show skipped examples")
	flag.StringVar(&schemaPath, "schema", "./schema", "specify location of schema directory")
}

func main() {
	flag.Parse()

	examples, err := extractExamples(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	var fail bool
	for _, example := range examples {
		if example.Err != nil {
			printFields("error", example.Mediatype, example.Title, example.Err)
			fail = true
			continue
		}

		schema, err := schemaByMediatype(schemaPath, example.Mediatype)
		if err != nil {
			if err == errSchemaNotFound {
				if skipped {
					printFields("skip", example.Mediatype, example.Title)

					if verbose {
						fmt.Println(example.Body, "---")
					}
				}
				continue
			}
		}

		// BUG(stevvooe): Recursive validation is not working. Need to
		// investigate. Will use this code as information for bug.
		document := gojsonschema.NewStringLoader(example.Body)
		result, err := gojsonschema.Validate(schema, document)

		if err != nil {
			printFields("error", example.Mediatype, example.Title, err)
			fmt.Println(example.Body, "---")
			fail = true
			continue
		}

		if !result.Valid() {
			// TOOD(stevvooe): This is nearly useless without file, line no.
			printFields("invalid", example.Mediatype, example.Title)
			for _, desc := range result.Errors() {
				printFields("reason", example.Mediatype, example.Title, desc)
			}
			fmt.Println(example.Body, "---")
			fail = true
			continue
		}

		printFields("ok", example.Mediatype, example.Title)
		if verbose {
			fmt.Println(example.Body, "---")
		}
	}

	if fail {
		os.Exit(1)
	}
}

var (
	specsByMediaType = map[string]string{
		"application/vnd.oci.image.manifest.v1+json":      "image-manifest-schema.json",
		"application/vnd.oci.image.manifest.list.v1+json": "manifest-list-schema.json",
	}

	errSchemaNotFound = errors.New("schema: not found")
	errFormatInvalid  = errors.New("format: invalid")
)

func schemaByMediatype(root, mediatype string) (gojsonschema.JSONLoader, error) {
	name, ok := specsByMediaType[mediatype]
	if !ok {
		return nil, errSchemaNotFound
	}

	if !filepath.IsAbs(root) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		root = filepath.Join(wd, root)
	}

	// lookup path
	path := filepath.Join(root, name)
	return gojsonschema.NewReferenceLoader("file://" + path), nil
}

// renderer allows one to incercept fenced blocks in markdown documents.
type renderer struct {
	blackfriday.Renderer
	fn func(text []byte, lang string)
}

func (r *renderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	r.fn(text, lang)
	r.Renderer.BlockCode(out, text, lang)
}

type example struct {
	Lang      string // gets raw "lang" field
	Title     string
	Mediatype string
	Body      string
	Err       error

	// TODO(stevvooe): Figure out how to keep track of revision, file, line so
	// that we can trace back verification output.
}

// parseExample treats the field as a syntax,attribute tuple separated by a comma.
// Attributes are encoded as a url values.
//
// An example of this is `json,title=Foo%20Bar&mediatype=application/json. We
// get that the "lang" is json, the title is "Foo Bar" and the mediatype is
// "application/json".
//
// This preserves syntax highlighting and lets us tag examples with further
// metadata.
func parseExample(lang, body string) (e example) {
	e.Lang = lang
	e.Body = body

	parts := strings.SplitN(lang, ",", 2)
	if len(parts) < 2 {
		e.Err = errFormatInvalid
		return
	}

	m, err := url.ParseQuery(parts[1])
	if err != nil {
		e.Err = err
		return
	}

	e.Mediatype = m.Get("mediatype")
	e.Title = m.Get("title")
	return
}

func extractExamples(rd io.Reader) ([]example, error) {
	p, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	var examples []example
	renderer := &renderer{
		Renderer: blackfriday.HtmlRenderer(0, "test test", ""),
		fn: func(text []byte, lang string) {
			examples = append(examples, parseExample(lang, string(text)))
		},
	}

	// just pass over the markdown and ignore the rendered result. We just want
	// the side-effect of calling back for each code block.
	// TODO(stevvooe): Consider just parsing these with a scanner. It will be
	// faster and we can retain file, line no.
	blackfriday.MarkdownOptions(p, renderer, blackfriday.Options{
		Extensions: blackfriday.EXTENSION_FENCED_CODE,
	})

	return examples, nil
}

// printFields prints each value tab separated.
func printFields(vs ...interface{}) {
	var ss []string
	for _, f := range vs {
		ss = append(ss, fmt.Sprint(f))
	}
	fmt.Println(strings.Join(ss, "\t"))
}
