package main

import (
	"fmt"
	ioc "mygo-ioc"
	"mygo-ioc/demo/mapper"
	"mygo-ioc/demo/repository"
	"mygo-ioc/demo/service"
)

func main() {

	// 新建一个容器
	container := ioc.NewContainer()

	// 将 DemoService 注册到容器
	var demoService service.DemoService
	container.Register(&demoService)

	fmt.Println("\n1.", "demoService.Print:")

	// 直接调用 demoService 的 Print 方法
	demoService.Print("hello world")

	fmt.Println("\n2.", "GetBeans:")

	// 遍历获取所有 Bean
	for _, bean := range container.GetBeans() {
		fmt.Println(bean.Tag(), bean.Name())
	}

	fmt.Println("\n3.", "GetBeansByPackage:")

	// 遍历获取所有属于 service 包的 Bean
	for _, bean := range container.GetBeansByPackage("service") {
		fmt.Println(bean.Tag(), bean.Name())
	}

	fmt.Println("\n4.", "GetBeansByTag:")

	// 遍历获取所有被标记为 mapper 的 Bean
	for _, bean := range container.GetBeansByTag("mapper") {
		fmt.Println(bean.Tag(), bean.Name())
		// 调用 print 方法
		bean.Any().(*mapper.DemoMapper).Print("hello world")
	}

	fmt.Println("\n5.", "GetBeanByName:")

	// 根据结构体名称获取指定 Bean
	beanByName := container.GetBeanByName("repository.DemoRepository").Any().(*repository.DemoRepository)
	// 调用 print 方法
	beanByName.Print("hello world")

	fmt.Println("\n6.", "use AOP")

	// 新建一个带 AOP 的容器
	containerWithAOP := ioc.NewContainer().Use(func(ctx *ioc.Context) {
		fmt.Println(ctx.Package(), ctx.Struct(), "aop start")
		ctx.Next()
		fmt.Println(ctx.Package(), ctx.Struct(), "aop end")
	})

	// 将 DemoService 注册到 AOP 容器
	containerWithAOP.Register(&service.DemoService{})

	// 通过 Call 方法调用 DemoRepository 的 Print 方法
	containerWithAOP.GetBeanByName("repository.DemoRepository").Call("Print", "hello world").Zero()
}
