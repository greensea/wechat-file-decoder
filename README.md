# wechat-file-decoder
一个微信缓存文件解码（解密）工具

## 使用说明

假设你的微信缓存目录是 FileStorage/Image/2020-12

首先我们需要执行下面的命令，猜测`编码异或值`

```shell
wechat-file-decoder FileStorage/Image/2020-12
```

你会看到类似这样的输出 

```shell
# ...
====扫描结果====
可能的异或值 | 出现数量 | 占比
ec               259    100.0
==============

猜测出的异或值是: ec
```

其中的 ec 就是`编码异或值`，每个微信账户的`编码异或值`都不同。

接着使用`编码异或值`解码文件，建议使用脚本进行批处理，在这里只演示转换单个文件

```
wechat-file-decoder ec FileStorage/Image/2020-12/a786850b9cba6440cbc84e1f9c25cfa4.dat /home/me/解码后的文件.jpg
```

## 安装

可以直接到 [release](https://github.com/greensea/wechat-file-decoder/releases/tag/v0.0.1) 页面下载。也可以直接使用 go 安装

```shell
go install github.com/greensea/wechat-file-decoder
```


## 参考资料

https://blog.csdn.net/a386115360/article/details/103215560
