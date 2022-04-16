package pipe

import "testing"

func TestPipe(t *testing.T) {

	type input struct {
		val int
	}

	cases := []struct {
		actions []action[input]
		want    input
	}{
		{
			actions: []action[input]{
				func(i *input) error { i.val = 1; return nil },
				func(i *input) error { i.val++; return nil },
				func(i *input) error { i.val *= 2; return nil },
			},
			want: input{val: 4},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			p := NewPipe[input]()
			for _, f := range c.actions {
				p.Next(f)
			}

			out, err := p.Do()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if got, want := out, c.want; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
