package internal_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/pkgx/internal"
)

func TestWrapAndUnwrap(t *testing.T) {
	cases := [][2]string{
		{"net", "net"},
		{"fmt", "fmt"},
		{"encoding/json", "xwrap_encoding_s_json"},
		{"github.com/path/to/pkg.Type", "xwrap_github_d_com_s_path_s_to_s_pkg_d_Type"},
		{"github.com/path/to/pkg_test.Type", "xwrap_github_d_com_s_path_s_to_s_pkg_u_test_d_Type"},
	}
	w := internal.NewWrapper()

	w.Clear()
	for _, c := range cases {
		Expect(t, w.Wrap(c[0]), Equal(c[1]))
	}

	w.Clear()
	for _, c := range cases {
		Expect(t, w.Unwrap(c[1]), Equal(c[0]))
	}

	w.Clear()
	Expect(t, w.Wrap("xwrap_net"), Equal("xwrap_net"))
	Expect(t, w.Unwrap("xwrap_net"), Equal("net"))
	Expect(t, w.Unwrap("xwrap_net"), Equal("net"))
	Expect(t, w.Unwrap("net.Addr"), Equal("net.Addr"))
	Expect(t, w.Wrap("net"), Equal("net"))
	Expect(t, w.Unwrap("xwrap_io_d_Reader"), Equal("io.Reader"))
	Expect(t, w.Unwrap("io.Reader"), Equal("io.Reader"))
	Expect(t, w.Unwrap("xwrap_io_d_Reader"), Equal("io.Reader"))
}
