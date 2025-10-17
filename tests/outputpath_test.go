package tests

import (
	"testing"

	"github.com/codeshelldev/goplater/utils/fsutils"
)

type outputPathCase struct {
	source string
	output string
	expect string
}

func TestOutputPath(t *testing.T) {
	cases := []outputPathCase{
		{
			source: "source.txt",
			output: "source.txt",
			expect: "source.txt",
		},
		{
			source: "source.txt",
			output: ".",
			expect: "source.txt",
		},
		{
			source: "source.txt",
			output: "other.txt",
			expect: "other.txt",
		},
		{
			source: "source/source.txt",
			output: "output/other.txt",
			expect: "output/other.txt",
		},
		{
			source: "source/source.txt",
			output: "output/",
			expect: "output/source.txt",
		},
		{
			source: "source/source.txt",
			output: "output/sub/",
			expect: "output/sub/source.txt",
		},
		{
			source: "source/sub/",
			output: "output/",
			expect: "output/",
		},
		{
			source: "source/sub/folder/source.txt",
			output: "output/",
			expect: "output/source.txt",
		},
		{
			source: "fs/source/sub/folder/source.txt",
			output: "fs/output/",
			expect: "fs/output/source.txt",
		},
	}

	for _, outputPathCase := range cases {
		got := fsutils.ResolveOutput(outputPathCase.source, outputPathCase.output, false)

		if got != outputPathCase.expect {
			t.Error("\nERROR\nsource:", outputPathCase.source, "\noutput:", outputPathCase.output, "\nexpected:", outputPathCase.expect, "\ngot:", got)
		}
	}
}

func TestOutputPathPreserve(t *testing.T) {
	cases := []outputPathCase{
		{
			source: "source.txt",
			output: "source.txt",
			expect: "source.txt",
		},
		{
			source: "source.txt",
			output: ".",
			expect: "source.txt",
		},
		{
			source: "source.txt",
			output: "other.txt",
			expect: "other.txt",
		},
		{
			source: "source/source.txt",
			output: "output/other.txt",
			expect: "output/other.txt",
		},
		{
			source: "source/source.txt",
			output: "output/",
			expect: "output/source.txt",
		},
		{
			source: "source/source.txt",
			output: "output/sub/",
			expect: "output/sub/source.txt",
		},
		{
			source: "source/sub/",
			output: "output/",
			expect: "output/sub/",
		},
		{
			source: "source/sub/folder/source.txt",
			output: "output/",
			expect: "output/sub/folder/source.txt",
		},
		{
			source: "fs/source/sub/folder/source.txt",
			output: "fs/output/",
			expect: "fs/output/sub/folder/source.txt",
		},
		{
			source: "fs/source/sub/folder/",
			output: "fs/output/",
			expect: "fs/output/sub/folder/",
		},
	}

	for _, outputPathCase := range cases {
		got := fsutils.ResolveOutput(outputPathCase.source, outputPathCase.output, true)

		if got != outputPathCase.expect {
			t.Error("\nERROR\nsource:", outputPathCase.source, "\noutput:", outputPathCase.output, "\nexpected:", outputPathCase.expect, "\ngot:", got)
		}
	}
}