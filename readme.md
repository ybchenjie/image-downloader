### 图片下载器

#### 语言

- golang

#### 环境
- mac
- windows

#### 特性

- 并发下载图片
- 递归创建图片途径
- 自动识别图片格式
- 零依赖

#### 用法

下载项目
```
git clone medgit@git.medlinker.com:chenjie/image-downloader.git
```
构建
```
sh deploy.sh
```
编写数据源文件
```
参考 data.tpl.json 数据结构，将数据源文件放在根目录，命名为data.json
```
执行
```
./mac
```

#### 支持的数据格式
```json
[
  {
    "name": "图片名字",
    "url": "http://pub-med-logo.medlinker.com/icon/guanlijihuajilu.png",
    "path": "book/pic1"
  },
  {
    "name": "213222",
    "url": "http://pub-med-logo.medlinker.com/icon/guanlijihuajilu.png",
    "path": "less/looo/lll"
  }
]
```

#### 数据格式说明
- data为数组结构，每个item是下载的图片项
- name 表示图片名称，系统会自动加后缀，注意保证同一路径下名字唯一
- url 表示图片地址
- path 表示路径，支持递归，根目录为脚本执行目录