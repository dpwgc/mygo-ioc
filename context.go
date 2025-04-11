package ioc

type Context struct {
	method string
	index  int
	In     Args
	Out    Args
	Cache  any
	bean   *Bean
}

type Handle func(ctx *Context)

func (c *Context) Name() string {
	return c.bean.name
}

func (c *Context) Method() string {
	return c.method
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.bean.handles) {
		c.bean.handles[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = len(c.bean.handles) + 1
}
