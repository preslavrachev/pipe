package pipe

type action[T any] func(input *T) error

type Pipe[T any] struct {
	onErr func(error) error
	funcs []action[T]
}

func NewPipe[T any]() *Pipe[T] {
	return &Pipe[T]{
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

func (p *Pipe[T]) Next(f action[T]) *Pipe[T] {
	p.funcs = append(p.funcs, f)
	return p
}

func (p *Pipe[T]) Do() (T, error) {
	var t T
	for _, f := range p.funcs {
		if err := f(&t); err != nil {
			// stop the chain prematurely
			err = p.onErr(err)
			return t, err
		}
	}

	p.onErr(nil)
	return t, nil
}
