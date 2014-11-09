package swan

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	PyExtractorsDir = "test_data/python-goose/extractors/"
)

type Result struct {
	URL      string
	Expected Expected
}

type Expected struct {
	MetaDescription string `json:"meta_description"`
	MetaKeywords    string `json:"meta_keywords"`
	MetaLang        string `json:"meta_lang"`
	MetaFavicon     string `json:"meta_favicon"`
	CleanedText     string `json:"cleaned_text"`
	Tags            []string
	Title           string
}

func testExtract(t *testing.T, name string, html string, e Expected) {
	a, err := FromHTML(html, Config{})
	if err != nil {
		t.Fatal(err)
	}

	if e.MetaDescription != a.Meta.Description {
		t.Fatalf(
			"%s: MetaDescription does not match:\n"+
				"	Got: %s\n"+
				"	Expected: %s",
			name, a.Meta.Description, e.MetaDescription)
	}
}

func TestPyExtractors(t *testing.T) {
	fs, err := ioutil.ReadDir(PyExtractorsDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".json") {
			r := Result{}
			name := strings.Replace(f.Name(), ".json", "", -1)

			jsonf, err := os.Open(PyExtractorsDir + f.Name())
			if err != nil {
				t.Fatal(err)
			}

			if err = json.NewDecoder(jsonf).Decode(&r); err != nil {
				t.Fatal(err)
			}

			h := strings.Replace(f.Name(), ".json", ".html", -1)
			htmlf, err := os.Open(PyExtractorsDir + h)
			if err != nil {
				t.Fatal(err)
			}

			html, err := ioutil.ReadAll(htmlf)
			if err != nil {
				t.Fatal(err)
			}

			testExtract(t, name, string(html), r.Expected)

			break
		}
	}
}
