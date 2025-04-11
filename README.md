# MyGo-IOC

## Go IOC 容器实现

***

### 示例代码

#### 完整版 -> `./demo`

```go
package main

import (
	"fmt"
	ioc "mygo-ioc"
)

// ===== Base Demo 程序 =====

func main() {

	// 新建一个容器
	container := ioc.NewContainer()

	// 注册两个 Repository 的接口实现
	container.RegisterImplement("FirstRepository", new(FirstRepository))
	container.RegisterImplement("SecondRepository", new(SecondRepository))

	// 将 DemoService 注册到 IOC 容器
	var demoService DemoService
	container.RegisterBeans(&demoService)

	fmt.Println("\n1.", "demoService.FirstInsert:")

	// 直接调用 demoService 的 FirstInsert 方法
	demoService.FirstInsert("hello world")

	fmt.Println("\n2.", "GetBeans:")

	// 遍历获取所有 Bean
	for _, bean := range container.GetBeans() {
		fmt.Println(bean.Name())
	}

	fmt.Println("\n5.", "GetBeanByImplementName:")

	// 根据接口实现名称获取指定 Bean
	beanByImplName := container.GetBeanByName("SecondRepository").Any().(*SecondRepository)
	// 调用 SecondRepository 的 Insert 方法
	beanByImplName.Insert("hello world")

	fmt.Println("\n6.", "GetBeanByStructName:")

	// 根据结构体名称获取指定 Bean
	beanByStructName := container.GetBeanByName("main.FirstRepository").Any().(*FirstRepository)
	// 调用 main.FirstRepository 的 Insert 方法
	beanByStructName.Insert("hello world")
}

// ===== Service 服务 =====

type DemoService struct {
	// 使用 autowired 注解标记需要自动注入的字段
	// 使用 qualifier 注解指定接口实现
	FirstRepository  Repository `qualifier:"FirstRepository"`
	SecondRepository Repository `qualifier:"SecondRepository"`
	// 不使用接口，直接注入实例指针
	FirstRepositoryImpl *FirstRepository
	// 禁用自动注入
	NoAutowired *FirstRepository `autowired:"false"`
}

func (s *DemoService) FirstInsert(text string) bool {
	return s.FirstRepository.Insert(text)
}

func (s *DemoService) SecondInsert(text string) bool {
	return s.SecondRepository.Insert(text)
}

func (s *DemoService) FirstImplInsert(text string) bool {
	return s.FirstRepositoryImpl.Insert(text)
}

// ===== Repository 服务 =====

type Repository interface {
	Insert(text string) bool
}

type FirstRepository struct {
}

func (fr *FirstRepository) Insert(text string) bool {
	fmt.Println("first insert:", text)
	return true
}

type SecondRepository struct {
}

func (sr *SecondRepository) Insert(text string) bool {
	fmt.Println("second insert:", text)
	return true
}
```