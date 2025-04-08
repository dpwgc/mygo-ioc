package ioc

import (
	"reflect"
)

type Container struct {
	beans       []*Bean
	beansByName map[string]*Bean
	beansByTag  map[string][]*Bean
	handles     []Handle
}

func NewContainer() *Container {
	return &Container{
		beansByName: make(map[string]*Bean),
		beansByTag:  make(map[string][]*Bean),
	}
}

func (c *Container) Use(middlewares ...Handle) *Container {
	c.handles = append(c.handles, middlewares...)
	return c
}

func (c *Container) Register(beans ...any) *Container {
	for _, bean := range beans {
		c.autowired(reflect.ValueOf(bean).Elem())
	}
	return c
}

func (c *Container) GetBeanByName(name string) *Bean {
	return c.beansByName[name]
}

func (c *Container) GetBeans() []*Bean {
	return c.beans
}

func (c *Container) GetBeansByTag(tag string) []*Bean {
	return c.beansByTag[tag]
}

func (c *Container) addBean(tag string, name string, rv reflect.Value, rt reflect.Type) {
	numMethod := rt.NumMethod()
	methods := make(map[string]reflect.Value, numMethod)
	for i := 0; i < numMethod; i++ {
		methods[rt.Method(i).Name] = rv.Method(i)
	}
	bean := &Bean{
		tag:     tag,
		name:    name,
		value:   rv.Interface(),
		handles: append([]Handle{}, c.handles...),
		methods: methods,
	}
	c.beans = append(c.beans, bean)
	c.beansByName[name] = bean
	c.beansByTag[tag] = append(c.beansByTag[tag], bean)
}

func (c *Container) autowired(val reflect.Value) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		switch field.Kind() {
		case reflect.Ptr:
			tag := fieldType.Tag.Get("autowired")
			if field.CanSet() && field.IsNil() && tag != "" {
				auto := reflect.New(fieldType.Type.Elem())
				field.Set(auto)
				c.addBean(tag, fieldType.Name, auto, fieldType.Type)
			}
			c.autowired(field.Elem())
		case reflect.Struct:
			c.autowired(field)
		}
	}
}
