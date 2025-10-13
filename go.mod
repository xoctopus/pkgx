module github.com/xoctopus/pkgx

go 1.25.1

require (
	github.com/onsi/gomega v1.38.2
	github.com/xoctopus/pkgx/testdata v0.0.0-20250331091630-3af90d68c457
	github.com/xoctopus/x v0.1.3-0.20251007035101-13b306000929
	golang.org/x/exp v0.0.0-20251009144603-d2f985daa21b
	golang.org/x/mod v0.29.0
	golang.org/x/tools v0.38.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)

replace github.com/xoctopus/pkgx/testdata => ./testdata
