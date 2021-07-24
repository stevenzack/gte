# GTE - Golang模版引擎

![logo](https://repository-images.githubusercontent.com/383689103/64c8877a-8516-4f53-8851-abe89cc2a7be)

用类似JSP的方式编写原生HTML/JS/CSS网站

我们支持对原生HTML网站的以下特性：

- 语法
    - [x][模板语法入门](grammar.md)
    - [x][导入html文件](import.md)
    - [x][请求上下文](context.md)
    - [x][内部接口请求](api.md)
    - [-]访问数据库
- 配置
    - [x][自定义路由](route.md)
    - [x][内置语言国际化](globalization.md)
- 打包
    - [x][支持gzip压缩支持](gzip.md)
    - [x][支持图片webp自动压缩](webp.md)

# 快速认识GTE

为了构建一个组件化的网站，我们将`.html`文件视为一个组件。组件和组件之间可以相互引用。

例如，我们在项目目录下新建一个`header.html`文件，和一个`footer.html`文件:

```html
<!-- header.html -->
<header>Header</header>
```

```html
<!-- footer.html -->
<footer>Footer</footer>
```

然后就可以在`index.html`文件中引用上述两个组件：
```html
<!-- index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Home</title>
</head>
<body>
    {{template "/header.html"}}
    <div>Content</div>
    {{template "/footer.html"}}
</body>
</html>
```

然后使用GTE命令，运行该项目：
```shell
gte serve
```

您将在浏览器看到运行结果：
```
Header
Content
Footer
```

# 接下来

- 语法
    - [模板语法入门](grammar.md)
    - [导入html文件](import.md)
    - [请求上下文](context.md)