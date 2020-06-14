package configutil

import "context"

// Float64Source is a type that can return a value.
type Float64Source interface {
	// Float should return a float64 if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Float64(context.Context) (*float64, error)
}

var (
	_ Float64Source = (*Float64)(nil)
)

// Float64 implements value provider.
type Float64 float64

// Float64 returns the value for a constant.
func (f Float64) Float64(_ context.Context) (*float64, error) {
	value := float64(f)
	return &value, nil
}
