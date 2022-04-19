A simple and somewhat idiomatic way of constructing easy-to-glance-at pipes for sequential processes. Inspired by Elixir's [pipe operator](https://elixirschool.com/en/lessons/basics/pipe_operator). Refer to my original [blog post](https://preslav.me/2021/09/04/generic-golang-pipelines/) addressing the idea.

`pipe` uses Go 1.18's generic type parameters, and has no external dependencies.

Still, you probably don't need it at all.

# Motivation

Error handling, mostly. I am all for Go's way of explicit error handling, but I find it a bit too verbose at times. Especially, in convoluted functions where error handling obscures the main line of business. 

`pipe`'s idea was to make it easy to glance at the main logic, by not taking any shortcusts at proper and idiomatic error handling. How is that possible? By cleverly splitting the logic that could fit into a chain-able pipeline: 

```go
func (s *Service) DoSomeComplexThing() error {
	_, err := pipe.New[sendOrderParams]().
		Next(s.loadCustomer).
		Next(s.loadProduct).
		Next(s.sendOrder).
		Do()

	return err
}
```

Feel free to check the entire example in the `examples` directory.

# Should you use it?

Probably not. I do, but that is no way an endorsement or motivation that you use it. In fact, the code is so short that I'd encourage anyone to rather copy it than use this package as a dependency. Ideas and suggestions are welcome, but I don't expect to be able to maintain it beyond the needs of my projects, so please, feel free to fork and experiment on your own.

