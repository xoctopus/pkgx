module github.com/xoctopus/pkgx

go 1.24.1

require (
	github.com/onsi/gomega v1.36.3
	github.com/xoctopus/pkgx/testdata v0.0.0-00010101000000-000000000000
	github.com/xoctopus/x v0.0.34
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394
	golang.org/x/mod v0.24.0
	golang.org/x/tools v0.31.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/xoctopus/pkgx/testdata => ./testdata
