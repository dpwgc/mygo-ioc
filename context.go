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

func (c *Context) Package() string {
	return c.bean.pkg
}

func (c *Context) Struct() string {
	return c.bean.name
}

func (c *Context) Tag() string {
	return c.bean.tag
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
