package ioc

import (
	"fmt"
	"reflect"
	"strings"
)

type Container struct {
	beans       []*Bean
	implsByName map[string]reflect.Value
	beansByName map[string]*Bean
	beansByTag  map[string][]*Bean
	handles     []Handle
}

func NewContainer() *Container {
	return &Container{
		implsByName: make(map[string]reflect.Value),
		beansByName: make(map[string]*Bean),
		beansByTag:  make(map[string][]*Bean),
	}
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

func (c *Container) GetBeansByTag(tag string) []*Bean {
	return c.beansByTag[tag]
}

func (c *Container) addBean(name, tag string, rv reflect.Value, rt reflect.Type) {
	if c.beansByName[name] != nil {
		return
	}
	numMethod := rt.NumMethod()
	methods := make(map[string]reflect.Value, numMethod)
	for i := 0; i < numMethod; i++ {
		methods[rt.Method(i).Name] = rv.Method(i)
	}
	bean := &Bean{
		tag:            tag,
		name:           name,
		value:          rv.Interface(),
		handles:        append([]Handle{}, c.handles...),
		reflectValue:   rv,
		reflectType:    rt,
		reflectMethods: methods,
	}
	c.beans = append(c.beans, bean)
	c.beansByName[name] = bean
	c.beansByTag[tag] = append(c.beansByTag[tag], bean)
}

func (c *Container) getBean(pkg, stu string) *Bean {
	return c.GetBeanByName(pkg + "." + stu)
}

func (c *Container) autowired(val reflect.Value) {
	typ := val.Type()
	fmt.Println("au", typ)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		switch field.Kind() {
		case reflect.Ptr:
		case reflect.Interface:
			autowired := fieldType.Tag.Get("autowired")
			if field.CanSet() && field.IsNil() && autowired != "" {
				path := fieldType.Type.String()
				items := strings.Split(strings.ReplaceAll(path, "*", ""), ".")
				bean := c.getBean(items[0], items[1])
				if bean != nil {
					field.Set(bean.reflectValue)
				} else {
					qualifier := fieldType.Tag.Get("qualifier")
					if qualifier != "" {
						rv := c.implsByName[qualifier]
						field.Set(rv)
						c.addBean(qualifier, autowired, rv, fieldType.Type)
					} else {
						rv := reflect.New(fieldType.Type.Elem())
						field.Set(rv)
						c.addBean(items[0]+"."+items[1], autowired, rv, fieldType.Type)
					}
				}
			}
			if field.Kind() == reflect.Interface {
				c.autowired(field.Elem().Elem())
			} else {
				c.autowired(field.Elem())
			}
		case reflect.Struct:
			c.autowired(field)
		}
	}
}
