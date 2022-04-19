package pipe

import "errors"

type actionFunc[T any] func(input *T) error

type action[T any] struct {
	fn  actionFunc[T]
	cfg runCfg
}

type Pipe[T any] struct {
	initFns []func(*T)
	actions []action[T]
	onErr   func(error) error
}

func New[T any](initFns ...func(*T)) *Pipe[T] {
	return &Pipe[T]{
		initFns: initFns,
		onErr: func(e error) error {
			/* do nothing */
			return nil
		},
	}
}

func (p *Pipe[T]) OnErr(handler func(error) error) *Pipe[T] {
	p.onErr = handler
	return p
}

type runCfg struct {
	permittedErrs []error
}

func (p *Pipe[T]) Next(fn actionFunc[T], opts ...func(opts *runCfg)) *Pipe[T] {

	cfg := runCfg{}

	for _, opt := range opts {
		opt(&cfg)
	}

	p.actions = append(p.actions, action[T]{fn: fn, cfg: cfg})
	return p
}

func (p *Pipe[T]) Do() (T, error) {
	var t T
	for _, initFn := range p.initFns {
		initFn(&t)
	}

outer:
	for _, action := range p.actions {
		if err := action.fn(&t); err != nil {
			for _, permittedErr := range action.cfg.permittedErrs {
				if errors.Is(err, permittedErr) {
					continue outer
				}
			}

			// stop the chain prematurely
			err = p.onErr(err)
			return t, err
		}
	}

	p.onErr(nil)
	return t, nil
}

func PermitErrors(errs ...error) func(opts *runCfg) {
	return func(opts *runCfg) {
		opts.permittedErrs = append(opts.permittedErrs, errs...)
	}
}
