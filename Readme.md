# Typora-Minio

该项目基于 Go 语言开发，利用 Minio 提供的 SDK 实现了在 Typora 中将图片上传到自建的 Minio 中

# 用法
1. 可下载构建好的文件
2. 下载 [配置文件](https://raw.githubusercontent.com/Chipmunkkk/typora-minio/master/minio.yaml) 并修改相关的配置
3. 在 Typora 的设置中修改图像=>上传服务设定,选择 `Custom Command`
4. 在命令一栏中输入`/path/to/typora_minio_xxx /path/to/minio.yaml`