# 导入html文件

为了构建组件化的网站，我们支持在`.html`文件中引入另外一个`.html`文件。

使用方法：
```shell
{{template "/component/head.html" .}}
```

上述代码中，用`{{...}}`包裹起来的部分属于[模板语法](grammar.md)。里面包含三个部分：

- `template`：函数，用于导入其他模板。
- `"/component/head.html"`：模板名称，填入你想要导入的`.html`文件的绝对路径（相对于项目根目录的绝对路径）
- `.`：传递给目标模板的参数，可以传入任何数据。这里的`.`代表当前请求的[上下文](context.md)，如果你想传入当前上下文的子字段，可以传比如`.Request.RequestURI`

