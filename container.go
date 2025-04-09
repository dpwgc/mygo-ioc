package ioc

import (
	"reflect"
	"strings"
)

type Container struct {
	beans       []*Bean
	beansByName map[string]*Bean
	beansByTag  map[string][]*Bean
	beansByPkg  map[string][]*Bean
	handles     []Handle
}

func NewContainer() *Container {
	return &Container{
		beansByName: make(map[string]*Bean),
		beansByTag:  make(map[string][]*Bean),
		beansByPkg:  make(map[string][]*Bean),
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

func (c *Container) GetBeansByPackage(pkg string) []*Bean {
	return c.beansByPkg[pkg]
}

func (c *Container) addBean(pkg, stu, tag string, rv reflect.Value, rt reflect.Type) {
	key := pkg + "." + stu
	if c.beansByName[key] != nil {
		return
	}
	numMethod := rt.NumMethod()
	methods := make(map[string]reflect.Value, numMethod)
	for i := 0; i < numMethod; i++ {
		methods[rt.Method(i).Name] = rv.Method(i)
	}
	bean := &Bean{
		tag:            tag,
		name:           stu,
		pkg:            pkg,
		value:          rv.Interface(),
		handles:        append([]Handle{}, c.handles...),
		reflectValue:   rv,
		reflectType:    rt,
		reflectMethods: methods,
	}
	c.beans = append(c.beans, bean)
	c.beansByName[key] = bean
	c.beansByTag[tag] = append(c.beansByTag[tag], bean)
	c.beansByPkg[pkg] = append(c.beansByPkg[pkg], bean)
}

func (c *Container) getBean(pkg, stu string) *Bean {
	return c.GetBeanByName(pkg + "." + stu)
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
				path := fieldType.Type.String()
				items := strings.Split(strings.ReplaceAll(path, "*", ""), ".")
				bean := c.getBean(items[0], items[1])
				if bean != nil {
					field.Set(bean.reflectValue)
				} else {
					rv := reflect.New(fieldType.Type.Elem())
					field.Set(rv)
					c.addBean(items[0], items[1], tag, rv, fieldType.Type)
				}
			}
			c.autowired(field.Elem())
		case reflect.Struct:
			c.autowired(field)
		}
	}
}
