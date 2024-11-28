
#### 介绍
用于下载小说并以 EPub 格式存储,支持Apple books(ibooks)自定义字体

#### How to use
##### (一) 爬取小说并转化为 EPub
1. 在 [github release界面](https://github.com/chcthink/freb/releases)下载对应系统的可执行文件
2. 在 [69书吧搜索界面](https://www.69yuedu.net/modules/article/search.php)找到要下载的小说
2. 进入`https://www.69yuedu.net/article/abcdefg.html` 介绍页,`abcdefg`
为该书本 ID
3. 在命令行输入以下命令下载小说

``` shell
# ID 为 69 书吧小说的 ID
./freb -i abcdefg ID
```

##### (二) txt 转 EPub

```shell
# -p 指定路径 -a 指定作者 -c 指定封面路径 默认为当前目录下的cover.png
# -o 输出路径
./freb -p xxx.txt -a xxx -o xxx.epub
```

#### Tips
 - 命令每次执行默认会从 github 下载静态文件,包括样式文件和配置文件暂存至本地,若想提高命令速度,可将代码库下的 `config.toml` 和`assets`目录下载至命令行同一目录
 - 可以通过修改 `assets`目录和 `config.toml` 文件来自定义 EPub
 - txt 读取整合 [kaf-cli](https://github.com/ystyle/kaf-cli)
 - 排版样式参考使用“阡陌居-笙歌夜夜”