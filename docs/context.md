# Context

Every HTTP request to the `.html` file will generate a `context` object, which you can access by using the `.` symbol in the code.

The context object contains the following information:

- 1. Current server environment configuration: `.Config`
- 2. Current request information: `.Request`
- 3. Some operation methods used to respond to requests: `.Response`

# .Config

The `.Config` object contains the environment and configuration information of all servers.

Field list:

- `.Config.Host`: `string` type, the Host monitored by the current server
- `.Config.Port`: `int` type, the port number currently monitored by the server
- `.Config.ApiServer`: `string` type, the server address requested by the API interface, see [Initiate an internal request](api.md) for details
- `.Config.Env`: `string` type, the name of the custom environment currently running, see [custom environment](env.md) for details
- `.Strs`: `map[string]map[string]string` type, resource packs for all languages. For example: `.Strs.zh-HK.HELLO_WORLD_` can get the translation value of `HELLO_WORLD_` under the `key` of the Chinese language pack. See [Internationalization](globalization.md) for details

# .Request

The `.Request` object contains information about the current request.

Field list:

- `.Request.Method`: `string` type, the method currently requested. E.g. `GET`
- `.Request.Proto`: `string` type, the protocol of the current request. For example, `HTTP/1.0`

Method list:

# .Response

The `.Response` object contains some operation methods of the response body

Method list

- `.SetStatusCode`:
    
    `func(int)` function type, used to set the HTTP status code of the response body (if not set, the HTTP status code defaults to 200)


    Scenario case:
    
    We need an article page to get the content of the article and return it according to the route `/a/{article-id}`, and return a 404 status code if it is not obtained.

    1.First, we define the `gte.config.json` configuration file:
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
    A route is defined here: `/a/:id`, and all requests conforming to this pattern will be processed by the `/article.html` file. For example: a request for `/a/my-first-article` is in line with the routing mode. After entering `/article.html`, we will regard the second half of the route `my-firsrt-article` as a name The route variable of `id` is passed, and you can get the value of `id` by way of `.Request.GetParam "id"`.

    The `blackList` here is a route blacklist, which prohibits users from directly accessing the `article.html` file through the `/article.html` route.


    2.Next we start to write the `article.html` file:
    ---
    
    ```html
    <!-- article.html -->
    {{$res := httpGetJson (print "http://localhost:8080/api/articles/" (.Request.GetParam "id")) }}
    {{if eq $res.StatusCode 200}}
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>{{$res.Data.Title}}-Article</title>
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
    Here we use the `httpGetJson` function to initiate an internal interface request (you need to prepare a back-end interface service to receive the request), the incoming parameter is a URL, and this URL is made by `print "http:// localhost:8080/api/articles/" (.Request.GetParam "id")` stitched together.

    `print` is a function for concatenating strings, the type is `JsonResponse func({any}...)`, any type and any number of parameters can be passed in.
    Because the second parameter of `print` is to get the second half of the current request route through `(.Request.GetParam "id")`, the value obtained is `"my-first-article"`, so the final splicing The result is: `"http://localhost:8080/api/articles/my-first-article"`

    ---

    After initiating the request, we define a `$res` variable to receive the response data of the `httpGetJson` function. (The `:=` symbol is used to declare a variable and assign a value). Next, we judge the response body `if eq $res.StatusCode 200` if the status code of the interface response is 200, indicating that the request for the article content data is successful, then the specific article content data will be displayed.

    The file content is obtained through `$res.Data`, which is a JSON object that contains the article title `$res.Data.Title` and the article content `$res.Data.Content`.

    If the response HTTP status code is not 200, the content after `else` is displayed

    ---