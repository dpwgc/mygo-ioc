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

	// 新建一个容器
	container := ioc.NewContainer()

	// 将 DemoService 注册到容器
	var demoService DemoService
	container.Register(&demoService)

	// 直接调用 demoService 的 Print 方法
	demoService.Print("hello world")

	// 遍历获取所有 Bean
	for _, bean := range container.GetBeans() {
		fmt.Println("GetBeans:", bean.Tag(), bean.Name())
	}

	// 遍历获取所有被标记为 mapper 的 Bean
	for _, bean := range container.GetBeansByTag("mapper") {
		fmt.Println("GetBeansByTag:", bean.Tag(), bean.Name())
		// 调用 print 方法
		bean.Any().(*DemoMapper).Print("hello world")
	}

	// 根据结构体名称获取指定 Bean
	beanByName := container.GetBeanByName("DemoRepository").Any().(*DemoRepository)
	// 调用 print 方法
	beanByName.Print("hello world")

	// 新建一个带 AOP 的容器
	containerWithAOP := ioc.NewContainer().Use(func(ctx *ioc.Context) {
		fmt.Println(ctx.BeanName(), "aop start")
		ctx.Next()
		fmt.Println(ctx.BeanName(), "aop end")
	})

	// 将 DemoService 注册到 AOP 容器
	containerWithAOP.Register(&DemoService{})

	// 通过 Call 方法调用 DemoRepository 的 Print 方法
	containerWithAOP.GetBeanByName("DemoRepository").Call("Print", "hello world").Zero()
}

// Print 方法

func (s *DemoService) Print(text string) {
	s.DemoRepository.Print(text)
}

func (r *DemoRepository) Print(text string) {
	r.DemoMapper.Print(text)
}

func (m *DemoMapper) Print(text string) {
	fmt.Println("print:", text)
}
```