package pipe

import (
	"errors"
	"reflect"
	"testing"
)

func TestPipe(t *testing.T) {

	type input struct {
		val int
	}

	permittedErr := errors.New("permitted error")

	cases := []struct {
		description string
		initFns     []func(*input)
		actions     []actionFunc[input]
		opts        [][]func(cfg *runCfg)
		want        input
	}{
		{
			description: "no actions",
			actions:     []actionFunc[input]{},
			want:        input{val: 0},
		},
		{
			description: "Initializer",
			initFns:     []func(*input){func(i *input) { i.val = 1 }},
			want:        input{val: 1},
		},
		{
			description: "All functions run till the end",
			actions: []actionFunc[input]{
				func(i *input) error { i.val = 1; return nil },
				func(i *input) error { i.val++; return nil },
				func(i *input) error { i.val *= 2; return nil },
			},
			want: input{val: 4},
		},
		{
			description: "An error stops the pipe",
			actions: []actionFunc[input]{
				func(i *input) error { i.val = 1; return nil },
				func(i *input) error { return errors.New("Oops") },
				func(i *input) error { i.val *= 2; return nil },
			},
			want: input{val: 1},
		},
		{
			description: "A permitted error does not stop the pipe",
			actions: []actionFunc[input]{
				func(i *input) error { i.val = 1; return nil },
				func(i *input) error { return permittedErr },
				func(i *input) error { i.val *= 2; return nil },
			},
			want: input{val: 2},
			opts: [][]func(opts *runCfg){
				1: {
					PermitErrors(permittedErr),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			p := New(c.initFns...)

			for idx, f := range c.actions {
				var opts []func(cfg *runCfg)
				if len(c.opts) > idx && c.opts[idx] != nil {
					opts = c.opts[idx]
				}
				p.Next(f, opts...)
			}

			out, err := p.Do()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if got, want := out, c.want; reflect.DeepEqual(got, want) == false {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
