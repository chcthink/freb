
#### 介绍
用于下载小说并以 EPub 格式存储,支持Apple books(ibooks)自定义字体

#### How to use
1. 在 [github release界面](https://github.com/chcthink/freb/releases)下载对应系统的可执行文件
2. 在 [69书吧搜索界面](https://69shuba.cx/modules/article/search.php)找到要下载的小说
2. 进入`https://69shuba.cx/book/12345.htm` 介绍页,`12345`
为该书本 ID
3. 在命令行输入以下命令下载小说
``` sh
# ID 为 69 书吧小说的 ID
./freb -i 46901 ID
```

#### Tips
 - 命令每次执行默认会从 github 下载静态文件,包括样式文件和配置文件暂存至本地,若想提高命令速度,可将代码库下的 `config.toml` 和`assets`目录下载至命令行同一目录
 - 可以通过修改 `assets`目录和 `config.toml` 文件来自定义 EPub
 - 排版样式参考使用“阡陌居-笙歌夜夜”