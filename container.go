package ioc

import (
	"fmt"
	"reflect"
	"strings"
)

type Container struct {
	defaultAutowired string
	beans            []*Bean
	implsByName      map[string]reflect.Value
	beansByName      map[string]*Bean
	handles          []Handle
}

func NewContainer() *Container {
	return &Container{
		defaultAutowired: "true",
		implsByName:      make(map[string]reflect.Value),
		beansByName:      make(map[string]*Bean),
	}
}

func (c *Container) DefaultAutowired(v bool) *Container {
	c.defaultAutowired = fmt.Sprintf("%v", v)
	return c
}

func (c *Container) Use(middlewares ...Handle) *Container {
	c.handles = append(c.handles, middlewares...)
	return c
}

func (c *Container) RegisterImplement(name string, impl any) *Container {
	c.implsByName[name] = reflect.ValueOf(impl)
	return c
}

func (c *Container) RegisterBeans(beans ...any) *Container {
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

func (c *Container) addBean(name string, rv reflect.Value, rt reflect.Type) {
	if c.beansByName[name] != nil {
		return
	}
	numMethod := rt.NumMethod()
	methods := make(map[string]reflect.Value, numMethod)
	for i := 0; i < numMethod; i++ {
		methods[rt.Method(i).Name] = rv.Method(i)
	}
	bean := &Bean{
		name:           name,
		value:          rv.Interface(),
		handles:        append([]Handle{}, c.handles...),
		reflectValue:   rv,
		reflectType:    rt,
		reflectMethods: methods,
	}
	c.beans = append(c.beans, bean)
	c.beansByName[name] = bean
}

const (
	True  = "true"
	False = "false"
)

func (c *Container) autowired(val reflect.Value) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		kind := field.Kind()
		if kind == reflect.Ptr || kind == reflect.Interface {
			autowired := fieldType.Tag.Get("autowired")
			if autowired != True && autowired != False {
				autowired = c.defaultAutowired
			}
			if field.CanSet() && field.IsNil() && autowired != False {
				qualifier := fieldType.Tag.Get("qualifier")
				key := qualifier
				if key == "" {
					path := fieldType.Type.String()
					items := strings.Split(strings.ReplaceAll(path, "*", ""), ".")
					key = items[0] + "." + items[1]
				}
				bean := c.GetBeanByName(key)
				if bean != nil {
					field.Set(bean.reflectValue)
				} else {
					var rv reflect.Value
					if qualifier != "" {
						rv = c.implsByName[qualifier]
					} else {
						rv = reflect.New(fieldType.Type.Elem())
					}
					field.Set(rv)
					c.addBean(key, rv, fieldType.Type)
				}
				if kind == reflect.Interface {
					c.autowired(field.Elem().Elem())
				} else {
					c.autowired(field.Elem())
				}
			}
		} else if kind == reflect.Struct {
			c.autowired(field)
		}
	}
}
