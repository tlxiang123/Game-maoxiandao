# YOLO 接入说明

本文档说明如何把当前可用的 YOLO 运行库复制到其它 AutoGo 项目，以及每次训练后的模型文件应该放在哪里。

## 目录说明

当前整理后的文件在：

```text
E:\autogo\go项目\yolo源代码+杂七杂八
```

保留了两个目录：

```text
yolo运行库_复制到项目
yolo源码_可重新编译
```

`yolo运行库_复制到项目` 是给其它项目直接复制用的运行包。

`yolo源码_可重新编译` 是改好的 C++ YOLO 源码，里面已经使用 Vulkan 版 ncnn，并兼容当前 YOLOv8 NCNN 输出。

## 新项目需要复制的固定文件

从：

```text
E:\autogo\go项目\yolo源代码+杂七杂八\yolo运行库_复制到项目
```

复制到目标项目根目录，保持目录结构不变：

```text
AutoGo\yolo
resources\libs\arm64-v8a\libyolo.so
resources\libs\x86\libyolo.so
resources\libs\x86_64\libyolo.so
```

如果只打真机 arm64 包，至少要复制：

```text
resources\libs\arm64-v8a\libyolo.so
```

如果要跑模拟器，通常还需要：

```text
resources\libs\x86_64\libyolo.so
```

## 每个模型项目需要放的文件

每次重新训练和转换后，把模型文件放到目标项目根目录：

```text
best.ncnn.param
best.ncnn.bin
data.yaml
```

这三个文件必须和同一次训练对应：

```text
best.ncnn.param
best.ncnn.bin
```

必须成对使用，不能混搭。

`data.yaml` 的 `names:` 顺序必须和训练模型时的类别顺序一致，否则识别出来的标签名会错位。

## go:embed 需要包含模型文件

目标项目的 `main.go` 或资源嵌入文件里，需要把模型和标签文件嵌进去，例如：

```go
//go:embed resources 720.txt best.ncnn.param best.ncnn.bin data.yaml
var res embed.FS
```

如果目标项目没有 `720.txt`，就按它自己的资源写法调整，但必须包含：

```text
resources
best.ncnn.param
best.ncnn.bin
data.yaml
```

## 代码调用方式

Go 代码里使用：

```go
detector := yolo.New("v8", 4, paramPath, binPath, labels)
if detector == nil {
    return false
}
defer detector.Close()

results := detector.DetectFromImage(img)
```

当前 `libyolo.so` 是 Vulkan 版：设备支持 Vulkan 时会走 GPU；不支持时会自动回 CPU。

## 重新编译 C++ 库

源码在：

```text
E:\autogo\go项目\yolo源代码+杂七杂八\yolo源码_可重新编译
```

里面保留了：

```text
src\Main.cpp
CMakeLists.txt
ncnn-20240820-android-vulkan
nlohmann_json
```

这个源码已经做过以下修改：

```text
YOLOv8 动态类别数
兼容 out0 / output 输出节点
兼容 in0 / images 输入节点
支持已解码的 YOLOv8 NCNN 输出
开启 ncnn Vulkan 自动检测
```

正常情况下，新项目只需要复制运行库和替换模型文件，不需要重新编译 C++。
