# 上下文

每一次对`.html`文件的HTTP请求，都会产生一个`上下文`对象，您可以在代码中通过`.`符号来访问。

上下文对象里面包含以下信息:

- 1.当前服务器环境配置：`.Config`
- 2.当前请求信息：`.Request`
- 3.用于响应请求的一些操作方法：`.Response`

# .Config

`.Config`对象包含了所有服务器的环境和配置信息。

字段列表：

- `.Config.Host`:`string`类型，当前服务器监听的Host
- `.Config.Port`:`int`类型，当前服务器监听的端口号
- `.Config.ApiServer`:`string`类型，API接口请求的服务器地址，详见[发起内部请求](api.md)
- `.Config.Env`:`string`类型，当前所运行的自定义环境名称，详见[自定义环境](env.md)
- `.Strs`:`map[string]map[string]string`类型，所有语言的资源包。例如：`.Strs.zh-HK.HELLO_WORLD_`可获取中文语言包下面`key`为`HELLO_WORLD_`的翻译值。详见[国际化](globalization.md)

# .Request

`.Request`对象包含了当前请求的信息。

字段列表：

- `.Request.Method`:`string`类型，当前请求的方法。例如`GET`
- `.Request.Proto`:`string`类型，当前请求的协议。例如`HTTP/1.0`

方法列表：

# .Response

`.Response`对象包含了响应体的一些操作方法

方法列表

- `.SetStatusCode`:
    
    `func(int)`函数类型，用于设置响应体的HTTP状态码（如果不设置的话，HTTP状态码默认是200）


    场景案例：
    
    我们需要一个文章页面，根据路由`/a/{article-id}`来获取文章内容并返回，如果没有获取到则返回404状态码。

    1.首先我们定义一下`gte.config.json`配置文件：
    ---
    ```json
    {
        "host": "localhost",
        "port": 8080,
        "routes": [{
                "path": "/a/:id",
                "to": "/article.html"
            }
        ],
        "blackList": [
            "/article.html"
        ]
    }
    ```
    这里定义了一个路由：`/a/:id`，所有符合这个模式的请求都会由`/article.html`文件来处理。例如：一个`/a/my-first-article`的请求是符合该路由模式的，进入`/article.html`之后，我们会把路由后半段的`my-firsrt-article`作为一个名为`id`的路由变量传递过去，您可以通过`.Request.GetParam "id"`的方式获取`id`的值。

    这里的`blackList`是路由黑名单，它禁止了用户直接通过`/article.html`路由来访问`article.html`文件。


    2.接下来我们开始编写`article.html`文件：
    ---
    ```html
    <!-- article.html -->
    {{$res := httpGetJson ( print "http://localhost:8080/api/articles/" (.Request.GetParam "id") ) }}
    {{if eq $res.StatusCode 200}}
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>{{$res.Data.Title}} - Article</title>
    </head>
    <body>
        <h1>{{$res.Data.Title}}</h1>
        <div>
            {{$res.Data.Content}}
        </div>
    </body>
    </html>
    {{else}}
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Article not found</title>
    </head>
    <body>
        404 article not found: {{$res.Error}}
    </body>
    </html>
    {{end}}
    ```
    这里我们通过`httpGetJson`函数来发起一个内部接口请求（您需要准备一个后端接口服务，用于接收该请求），传入的参数是一个URL，而这个URL是由`print "http://localhost:8080/api/articles/" (.Request.GetParam "id")`拼接而成的。

    `print`是一个拼接字符串的函数，类型为`JsonResponse func({any}...)`，可传入任何类型，任何数量的参数。
    因为`print`的第二个参数是通过`(.Request.GetParam "id")`获取当前请求路由中的后半段，得到的值为`"my-first-article"`，所以拼接的最终结果为：`"http://localhost:8080/api/articles/my-first-article"`

    ---

    发起请求之后，我们通过定义一个`$res`变量来接收`httpGetJson`函数的响应数据。（`:=`符号用于声明一个变量并赋值）。接下来我们判断响应体`if eq $res.StatusCode 200`如果接口响应的状态码为200，说明请求文章内容数据成功，则显示具体的文章内容数据，

    文件内容通过`$res.Data`来获取，这是一个JSON对象，里面包含了文章标题`$res.Data.Title`和文章内容`$res.Data.Content`。

    如果响应的HTTP状态码不为200，则显示`{{else}}`后面的内容

    ---

    