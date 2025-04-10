package main

import (
	"fmt"
	ioc "mygo-ioc"
	"mygo-ioc/demo/base/repository"
	"mygo-ioc/demo/base/service"
)

func main() {

	// 新建一个容器
	container := ioc.NewContainer()

	container.RegisterImplement("FirstRepository", new(repository.FirstRepository))
	container.RegisterImplement("SecondRepository", new(repository.SecondRepository))

	// 将 DemoService 注册到容器
	var demoService service.DemoService
	container.RegisterBeans(&demoService)

	fmt.Println("\n1.", "demoService.Add:")

	// 直接调用 demoService 的 FirstInsert 方法
	demoService.FirstInsert("hello world")

	fmt.Println("\n2.", "GetBeans:")

	// 遍历获取所有 Bean
	for _, bean := range container.GetBeans() {
		fmt.Println(bean.Tag(), bean.Name())
	}

	fmt.Println("\n4.", "GetBeansByTag:")

	// 遍历获取所有被标记为 repository 的 Bean
	for _, bean := range container.GetBeansByTag("repository") {
		fmt.Println(bean.Tag(), bean.Name())
	}

	fmt.Println("\n5.", "GetBeanByName:")

	// 根据结构体名称获取指定 Bean
	beanByName := container.GetBeanByName("SecondRepository").Any().(*repository.SecondRepository)
	// 调用 SecondRepository 的 Insert 方法
	beanByName.Insert("hello world")
}
