module github.com/xoctopus/pkgx

go 1.25.3

require (
	github.com/xoctopus/pkgx/testdata v0.0.0-20250331091630-3af90d68c457
	github.com/xoctopus/x v0.2.1-0.20251113131642-d51532724ff6
	golang.org/x/exp v0.0.0-20251113190631-e25ba8c21ef6
	golang.org/x/mod v0.30.0
	golang.org/x/tools v0.39.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xoctopus/errx v0.0.0-20251110065924-ca3348b82575 // indirect
	golang.org/x/sync v0.18.0 // indirect
)

replace github.com/xoctopus/pkgx/testdata => ./testdata
