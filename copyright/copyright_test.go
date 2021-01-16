/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT
license that can be found in the LICENSE file.

*/

package copyright

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Copyright_mergeFileSections(t *testing.T) {
	its := assert.New(t)

	merged := Copyright{}.mergeFileSections([]byte("foo"), []byte("bar"), []byte("baz"))
	its.Equal("foobarbaz", string(merged))
}

func Test_Copyright_fileHasCopyrightHeader(t *testing.T) {
	its := assert.New(t)

	var goodCorpus = []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(2021)
	its.Nil(err)

	goodCorpusWithNotice := Copyright{}.mergeFileSections([]byte(notice), goodCorpus)
	its.Contains(string(goodCorpusWithNotice), "Copyright (c) 2021")
	its.True((Copyright{}).fileHasCopyrightHeader(goodCorpusWithNotice, []byte(notice)))
}

func Test_Copyright_fileHasCopyrightHeader_invalid(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	var invalidCorpus = []byte(`foo
bar
baz
`)
	expectedNotice, err := generateGoNotice(2021)
	its.Nil(err)

	its.False(c.fileHasCopyrightHeader(invalidCorpus, []byte(expectedNotice)), "we haven't added the notice")
}

func Test_Copyright_fileHasCopyrightHeader_differentYear(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	var goodCorpus = []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(2020)
	its.Nil(err)

	goodCorpusWithNotice := c.mergeFileSections(notice, goodCorpus)
	its.Contains(string(goodCorpusWithNotice), "Copyright (c) 2020")

	newNotice, err := generateGoNotice(2021)
	its.Nil(err)

	its.True(c.fileHasCopyrightHeader(goodCorpusWithNotice, []byte(newNotice)))
}

func Test_Copyright_fileHasCopyrightHeader_leadingWhitespace(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	var goodCorpus = []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(2021)
	its.Nil(err)

	goodCorpusWithNotice := c.mergeFileSections([]byte("\n\n"), notice, goodCorpus)
	its.HasPrefix(string(goodCorpusWithNotice), "\n\n")
	its.Contains(string(goodCorpusWithNotice), "Copyright (c) 2021")

	its.True(c.fileHasCopyrightHeader(goodCorpusWithNotice, []byte(notice)))
}

func Test_Copyright_goBuildTagMatch(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	buildTag := []byte(`// +build foo

`)
	corpus := []byte(`foo
bar
baz
`)

	file := (Copyright{}).mergeFileSections(buildTag, corpus)

	its.False(goBuildTagMatch.Match(corpus))
	its.True(goBuildTagMatch.Match(c.mergeFileSections(buildTag)))

	found := goBuildTagMatch.FindAll(file, -1)
	its.NotEmpty(found)
	its.True(goBuildTagMatch.Match(file))
}

func Test_Copyright_goInjectNotice(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	file := []byte(`foo
bar
baz
`)

	notice, err := generateGoNotice(2021)
	its.Nil(err)

	output := c.goInjectNotice("foo.go", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.HasSuffix(string(output), string(file))
}

func Test_Copyright_goInjectNotice_buildTags(t *testing.T) {
	its := assert.New(t)
	c := Copyright{}

	buildTag := []byte(`// +build foo`)
	corpus := []byte(`foo
bar
baz
`)

	file := c.mergeFileSections(buildTag, []byte("\n\n"), corpus)

	notice, err := generateGoNotice(2021)
	its.Nil(err)

	output := c.goInjectNotice("foo.go", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.HasPrefix(string(output), string(buildTag)+"\n")
	its.HasSuffix(string(output), string(corpus))

	outputRepeat := c.goInjectNotice("foo.go", output, notice)
	its.Empty(outputRepeat, "inject notice functions should return an empty slice if the header already exists")
}

func Test_Copyright_injectNotice_typescript(t *testing.T) {
	its := assert.New(t)

	c := Copyright{}

	file := []byte(`foo
bar
baz
`)

	notice, err := generateTypescriptNotice(2021)
	its.Nil(err)

	output := c.injectNotice("foo.ts", file, notice)
	its.Contains(string(output), "Copyright (c) 2021")
	its.HasSuffix(string(output), string(file))

	outputRepeat := c.injectNotice("foo.ts", output, notice)
	its.Empty(outputRepeat, "inject notice functions should return an empty slice if the header already exists")
}

func generateGoNotice(year int) ([]byte, error) {
	noticeBody, err := (Copyright{
		Config: Config{
			Year: year,
		},
	}).compileNoticeBodyTemplate(DefaultNoticeBodyTemplate)
	if err != nil {
		return nil, err
	}

	compiled, err := (Copyright{}).compileNoticeTemplate(goNoticeTemplate, noticeBody)
	if err != nil {
		return nil, err
	}
	return []byte(compiled), nil
}

func generateTypescriptNotice(year int) ([]byte, error) {
	noticeBody, err := (Copyright{
		Config: Config{
			Year: year,
		},
	}).compileNoticeBodyTemplate(DefaultNoticeBodyTemplate)
	if err != nil {
		return nil, err
	}

	compiled, err := (Copyright{}).compileNoticeTemplate(tsNoticeTemplate, noticeBody)
	if err != nil {
		return nil, err
	}
	return []byte(compiled), nil
}
