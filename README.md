# MyGo-IOC

## Go IOC 容器实现

***

### 示例代码

```go
package main

import (
	"fmt"
	ioc "mygo-ioc"
)

type DemoService struct {
	DemoRepository *DemoRepository `autowired:"repository"`
}

type DemoRepository struct {
	DemoMapper *DemoMapper `autowired:"mapper"`
}

type DemoMapper struct{}

func main() {

	// 将 DemoService 注册到容器
	var demoService DemoService
	ioc.RegisterBeans(&demoService)

	// 调用 print 方法
	demoService.print("hello world")

	// 遍历获取所有 Bean
	for _, bean := range ioc.GetBeans() {
		fmt.Println("GetBeans:", bean.Tag(), bean.Name())
	}

	// 遍历获取所有被标记为 mapper 的 Bean
	for _, bean := range ioc.GetBeansByTag("mapper") {
		fmt.Println("GetBeansByTag:", bean.Tag(), bean.Name())
		// 调用 print 方法
		bean.Any().(*DemoMapper).print("hello world")
	}

	// 根据结构体名称获取指定 Bean
	beanByName := ioc.GetBeanByName("DemoRepository").(*DemoRepository)
	// 调用 print 方法
	beanByName.print("hello world")

	// 根据结构体类型获取指定 Bean
	beanByType := ioc.GetBeanByType(&DemoRepository{})
	// 调用 print 方法
	beanByType.print("hello world")
}

// print 方法

func (s *DemoService) print(text string) {
	s.DemoRepository.print(text)
}

func (r *DemoRepository) print(text string) {
	r.DemoMapper.print(text)
}

func (m *DemoMapper) print(text string) {
	fmt.Println("print:", text)
}
```