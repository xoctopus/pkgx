module github.com/xoctopus/pkgx

go 1.25.1

require (
	github.com/xoctopus/pkgx/testdata v0.0.0-20250331091630-3af90d68c457
	github.com/xoctopus/x v0.2.0
	golang.org/x/exp v0.0.0-20251017212417-90e834f514db
	golang.org/x/mod v0.29.0
	golang.org/x/tools v0.38.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sync v0.17.0 // indirect
)

replace github.com/xoctopus/pkgx/testdata => ./testdata
