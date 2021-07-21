# GTE - Golang模版引擎

用类似JSP的方式编写原生HTML/JS/CSS网站

# 快速认识GTE

为了构建一个组件化的网站，我们将`.html`文件视为一个组件。组件和组件之间可以相互引用。

例如，我们在项目目录下新建一个`header.html`文件，和一个`footer.html`文件:

```html
<!-- header -->
<header>Header</header>
```

```html
<!-- footer -->
<footer>Footer</footer>
```

然后将可以在`index.html`文件中引用上述两个组件：
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
