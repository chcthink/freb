#### 介绍

 - 下载小说并以 EPub 格式存储
 - 转化 TXT 至 EPub

> 限制爬取频率防止被 ban
> 默认频率(config.toml delay_time)为 每章节/0.8s
> 已测试在域名`69shuba.cx` 下 每章节/0.8s不会被 ban

#### 功能
 - 支持Apple books(ibooks)自定义字体
 - 删除多余字符

#### 预览
<p align="center">
  <img src="docs/assets/img/git-intro.png" width="30%" height="auto" />
  <img src="docs/assets/img/git-vol.png" width="30%" height="auto"  />
  <img src="docs/assets/img/git-chapter.png" width="30%"  height="auto" />
</p>
#### How to use

##### (一) 爬取小说并转化为 EPub
> 69 书吧地址 
> 1. [69 书吧](69shuba.cx)         书籍 ID 为数字
> 2. [69 阅读](www.69yuedu.net)    书籍 ID 为字母
> 某些小说(书籍 ID 为字母)会从 69 书吧 跳转到 69 阅读

1. 在 [github release界面](https://github.com/chcthink/freb/releases)下载对应系统的可执行文件
2. 在 69书吧搜索界面找到要下载的小说
3. 进入`https://www.69yuedu.net/article/abcdefg.html` 或者`https://69shuba.cx/book/123456.htm` 介绍页,`abcdefg`或`123456`
为该书本 ID
4. 自行定义 config.toml 配置,如 每章下载延时delay_time
5. 在命令行输入以下命令下载小说

> 当小说 id 为**数字**时,可以不用 -c 指定封面,工具会自动下载封面(若69书吧存在封面)
``` shell
# ID 为 69 书吧小说的 ID
./freb -i abcdefg ID
```

##### (二) txt 转 EPub

```shell
# -p 指定路径 -a 指定作者 -c 指定封面路径 默认为当前目录下的cover.png
# -o 输出路径
./freb -p xxx.txt -a xxx -c cover.jpg -o xxx.epub
```

#### Tips
 - 当书籍 ID 为数字且在搜索列表存在封面情况下无法获取封面,尝试重新执行命令
 - 命令每次执行默认会从 github 下载静态文件,包括样式文件和配置文件暂存至本地,若想提高命令速度,可将代码库下的 `config.toml` 和`assets`目录下载至命令同一目录
 - 可以通过修改 `assets`目录和 `config.toml` 文件来自定义 EPub
 - txt 读取整合 [kaf-cli](https://github.com/ystyle/kaf-cli)
 - 排版样式参考使用“阡陌居-笙歌夜夜”
 - 图片来源 pixabay 作者: BiancaVanDijk
 - 字体采用 [975 Maru](https://github.com/lxgw/975maru)
