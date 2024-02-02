# 不理解Alis源码的callStorageHook的逻辑

```text

package op

import (
	"regexp"
	"strings"

	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Obj
type ObjsUpdateHook = func(parent string, objs []model.Obj)

var (
	objsUpdateHooks = make([]ObjsUpdateHook, 0)
)

func RegisterObjsUpdateHook(hook ObjsUpdateHook) {
	objsUpdateHooks = append(objsUpdateHooks, hook)
}

func HandleObjsUpdateHook(parent string, objs []model.Obj) {
	for _, hook := range objsUpdateHooks {
		hook(parent, objs)
	}
}

// Setting
type SettingItemHook func(item *model.SettingItem) error

var settingItemHooks = map[string]SettingItemHook{
	conf.VideoTypes: func(item *model.SettingItem) error {
		conf.SlicesMap[conf.VideoTypes] = strings.Split(item.Value, ",")
		return nil
	},
	conf.AudioTypes: func(item *model.SettingItem) error {
		conf.SlicesMap[conf.AudioTypes] = strings.Split(item.Value, ",")
		return nil
	},
	conf.ImageTypes: func(item *model.SettingItem) error {
		conf.SlicesMap[conf.ImageTypes] = strings.Split(item.Value, ",")
		return nil
	},
	conf.TextTypes: func(item *model.SettingItem) error {
		conf.SlicesMap[conf.TextTypes] = strings.Split(item.Value, ",")
		return nil
	},
	conf.ProxyTypes: func(item *model.SettingItem) error {
		conf.SlicesMap[conf.ProxyTypes] = strings.Split(item.Value, ",")
		return nil
	},
	conf.ProxyIgnoreHeaders: func(item *model.SettingItem) error {
		conf.SlicesMap[conf.ProxyIgnoreHeaders] = strings.Split(item.Value, ",")
		return nil
	},
	conf.PrivacyRegs: func(item *model.SettingItem) error {
		regStrs := strings.Split(item.Value, "\n")
		regs := make([]*regexp.Regexp, 0, len(regStrs))
		for _, regStr := range regStrs {
			reg, err := regexp.Compile(regStr)
			if err != nil {
				return errors.WithStack(err)
			}
			regs = append(regs, reg)
		}
		conf.PrivacyReg = regs
		return nil
	},
	conf.FilenameCharMapping: func(item *model.SettingItem) error {
		err := utils.Json.UnmarshalFromString(item.Value, &conf.FilenameCharMap)
		if err != nil {
			return err
		}
		log.Debugf("filename char mapping: %+v", conf.FilenameCharMap)
		return nil
	},
	conf.IgnoreDirectLinkParams: func(item *model.SettingItem) error {
		conf.SlicesMap[conf.IgnoreDirectLinkParams] = strings.Split(item.Value, ",")
		return nil
	},
}

func RegisterSettingItemHook(key string, hook SettingItemHook) {
	settingItemHooks[key] = hook
}

func HandleSettingItemHook(item *model.SettingItem) (hasHook bool, err error) {
	if hook, ok := settingItemHooks[item.Key]; ok {
		return true, hook(item)
	}
	return false, nil
}

// Storage
type StorageHook func(typ string, storage driver.Driver)

var storageHooks = make([]StorageHook, 0)

func callStorageHooks(typ string, storage driver.Driver) {
	for _, hook := range storageHooks {
		hook(typ, storage)
	}
}

func RegisterStorageHook(hook StorageHook) {
	storageHooks = append(storageHooks, hook)
}
```

在提供的Go代码中，`callStorageHooks` 函数是用来调用一系列预先注册的 `StorageHook` 函数的。我将解释一下整个过程：

1. **StorageHook类型定义**：
   `StorageHook` 被定义为一个函数类型，接受两个参数：一个字符串 `typ` 和一个 `driver.Driver` 类型的 `storage`。这意味着任何符合这个签名的函数都可以被作为 `StorageHook` 使用。

2. **存储钩子的注册**：
   `storageHooks` 是一个 `StorageHook` 类型的切片，用来存储所有注册的钩子函数。`RegisterStorageHook` 函数允许你将一个新的 `StorageHook` 函数添加到这个切片中。

3. **调用存储钩子**：
   当 `callStorageHooks` 被调用时，它会遍历 `storageHooks` 切片中的所有钩子函数，并且依次调用它们，传入 `typ` 和 `storage` 作为参数。

4. **钩子函数的作用**：
   每个钩子函数可能会执行不同的操作，这取决于它们的具体实现。这可能包括初始化存储驱动、更新配置、记录日志、执行安全检查等。具体行为由注册到 `storageHooks` 中的函数决定。

5. **结果**：
   在调用 `callStorageHooks` 后，所有注册的 `StorageHook` 都会按照它们注册的顺序被执行。这是一个典型的观察者模式实现，允许代码在特定事件（例如存储初始化）发生时通知一系列观察者（即钩子函数）。

简而言之，调用 `callStorageHooks` 将触发一系列预先定义的操作，这些操作与存储系统的某些方面相关，例如初始化、配置更新或其他自定义逻辑。每个钩子都是独立定义的，因此它们的具体行为将根据你的应用需求和逻辑而有所不同。

# 不明白函数hook的执行逻辑在那里
```text
// Storage
type StorageHook func(typ string, storage driver.Driver)

var storageHooks = make([]StorageHook, 0)

func callStorageHooks(typ string, storage driver.Driver) {
	for _, hook := range storageHooks {
		hook(typ, storage)
	}
}

func RegisterStorageHook(hook StorageHook) {
	storageHooks = append(storageHooks, hook)
}
----
我就想问，type StorageHook func(typ string, storage driver.Driver) 声明了StorageHook是一个函数，那么这个函数的执行逻辑在那里？特别是callStorageHooks里面怎么用hook(typ, storage)执行逻辑
```
在您的代码中，`StorageHook` 被定义为一个函数类型。这种方式在Go中常用于声明回调函数或钩子（hooks）。`StorageHook` 的具体逻辑不在它的声明中定义，而是在使用这个类型的地方定义。

这里的关键是理解 `StorageHook` 类型的函数是如何被注册和调用的：

1. **StorageHook 函数类型定义**:
   ```go
   type StorageHook func(typ string, storage driver.Driver)
   ```
   这定义了 `StorageHook` 为一个接受一个字符串和一个 `driver.Driver` 类型参数的函数。

2. **注册 StorageHook**:
   `RegisterStorageHook` 函数允许其他代码注册一个符合 `StorageHook` 签名的函数。当你调用 `RegisterStorageHook`，你传递了一个具体的函数，这个函数包含了实际要执行的逻辑。

   例如：
   ```go
   RegisterStorageHook(myStorageHookFunction)
   ```
   其中 `myStorageHookFunction` 是你定义的一个符合 `StorageHook` 签名的函数。

3. **StorageHook 函数的存储**:
   所有通过 `RegisterStorageHook` 注册的函数都存储在 `storageHooks` 切片中。

4. **调用 StorageHook**:
   `callStorageHooks` 函数遍历 `storageHooks` 切片，并对每一个存储的 `StorageHook` 函数调用 `hook(typ, storage)`。这里的 `hook` 是切片中的一个 `StorageHook` 函数，当它被调用时，它执行的是在注册该函数时定义的具体逻辑。

简单来说，`StorageHook` 的具体执行逻辑是在注册这个钩子函数时定义的，而不是在 `StorageHook` 类型声明时定义的。`callStorageHooks` 函数仅仅是遍历所有注册的钩子，并按照它们注册时定义的逻辑执行它们。

# 还是不明白，请举例子
当然可以。让我们通过一个简单的例子来演示如何定义、注册和调用 `StorageHook` 类型的函数。

1. **定义 StorageHook 函数类型**:
   这一步已经在您提供的代码中完成了。`StorageHook` 是一个接受一个字符串和一个 `driver.Driver` 类型的参数的函数。

2. **实现一个具体的 StorageHook 函数**:
   我们将创建一个简单的 `StorageHook` 函数实例，它只是打印传入的参数。

   ```go
   func myStorageHook(typ string, storage driver.Driver) {
       fmt.Printf("StorageHook called with type: %s and storage: %v\n", typ, storage)
   }
   ```

   假设 `driver.Driver` 是一个接口，我们可以简单地使用空接口（`interface{}`）来模拟它，因为具体的实现对于理解这个例子来说并不重要。

3. **注册 StorageHook 函数**:
   接下来，我们将使用 `RegisterStorageHook` 函数将我们的 `myStorageHook` 函数注册为一个钩子。

   ```go
   RegisterStorageHook(myStorageHook)
   ```

4. **调用 StorageHooks**:
   最后，我们将调用 `callStorageHooks` 函数，该函数将遍历所有注册的钩子并执行它们。

下面是完整的示例代码：

```go
package main

import (
    "fmt"
)

// 假设的 driver.Driver 接口
type Driver interface{}

// StorageHook 函数类型
type StorageHook func(typ string, storage Driver)

// 存储所有注册的 StorageHook 函数
var storageHooks = make([]StorageHook, 0)

// 注册 StorageHook 函数
func RegisterStorageHook(hook StorageHook) {
    storageHooks = append(storageHooks, hook)
}

// 调用所有注册的 StorageHook 函数
func callStorageHooks(typ string, storage Driver) {
    for _, hook := range storageHooks {
        hook(typ, storage)
    }
}

// 实现一个 StorageHook 函数
func myStorageHook(typ string, storage Driver) {
    fmt.Printf("StorageHook called with type: %s and storage: %v\n", typ, storage)
}

func main() {
    // 注册 StorageHook
    RegisterStorageHook(myStorageHook)

    // 假设的 storage 对象
    var storage Driver // 这里可以是任何符合 Driver 接口的实例

    // 调用 StorageHooks
    callStorageHooks("exampleType", storage)
}
```

在这个例子中，当 `callStorageHooks` 被调用时，它会遍历所有注册的 `StorageHook` 函数（在这个例子中只有 `myStorageHook`），并且对每个钩子执行它们的逻辑，这里的逻辑是打印出传入的类型和存储对象。