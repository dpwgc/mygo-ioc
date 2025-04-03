package ioc

import (
	"reflect"
	"sync"
)

var beans []*Bean
var beansByName = make(map[string]*Bean)
var beansByTag = make(map[string][]*Bean)
var beanSafe = sync.Mutex{}

type Bean struct {
	tag   string
	name  string
	value any
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

func RegisterBeans(beans ...any) {
	for _, bean := range beans {
		RegisterBean(bean)
	}
}

func RegisterBean(bean any) {
	val := reflect.ValueOf(bean).Elem()
	autowired(val)
}

func GetBeanByType[T any](typ T) T {
	return GetBeanByName(reflect.TypeOf(typ).Elem().Name()).(T)
}

func GetBeanByName(name string) any {
	beanSafe.Lock()
	defer beanSafe.Unlock()
	return beansByName[name].value
}

func GetBeans() []*Bean {
	beanSafe.Lock()
	defer beanSafe.Unlock()
	return beans
}

func GetBeansByTag(tag string) []*Bean {
	beanSafe.Lock()
	defer beanSafe.Unlock()
	return beansByTag[tag]
}

func addBean(tag string, name string, v any) {
	beanSafe.Lock()
	defer beanSafe.Unlock()
	bean := &Bean{
		tag:   tag,
		name:  name,
		value: v,
	}
	beans = append(beans, bean)
	beansByName[name] = bean
	beansByTag[tag] = append(beansByTag[tag], bean)
}

func autowired(val reflect.Value) {
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
				addBean(tag, fieldType.Name, auto.Interface())
			}
			autowired(field.Elem())
		case reflect.Struct:
			autowired(field)
		}
	}
}
