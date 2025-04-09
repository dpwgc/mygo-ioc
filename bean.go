package ioc

import "reflect"

type Bean struct {
	pkg            string
	tag            string
	name           string
	value          any
	handles        []Handle
	reflectValue   reflect.Value
	reflectType    reflect.Type
	reflectMethods map[string]reflect.Value
}

func (b *Bean) Use(middlewares ...Handle) *Bean {
	b.handles = append(b.handles, middlewares...)
	return b
}

func (b *Bean) Tag() string {
	return b.tag
}

func (b *Bean) Name() string {
	return b.name
}

func (b *Bean) Any() any {
	return b.value
}

func (b *Bean) Call(name string, args ...any) Args {
	ctx := &Context{
		method: name,
		In:     args,
		bean:   b,
	}
	var result Args
	ctx.bean.handles = append(ctx.bean.handles, func(ctx *Context) {
		var argValues []reflect.Value
		for i := 0; i < len(ctx.In); i++ {
			argValues = append(argValues, reflect.ValueOf(ctx.In[i]))
		}
		returnValues := b.reflectMethods[name].Call(argValues)
		for i := 0; i < len(returnValues); i++ {
			result = append(result, returnValues[i].Interface())
		}
		ctx.Out = result
	})

	for ctx.index < len(ctx.bean.handles) {
		ctx.bean.handles[ctx.index](ctx)
		ctx.index++
	}
	return result
}
