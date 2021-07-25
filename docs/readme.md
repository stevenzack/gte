# GTE - Golang Template Engine [中文](zh/readme.md)

![logo](https://repository-images.githubusercontent.com/383689103/64c8877a-8516-4f53-8851-abe89cc2a7be)

Write native HTML/JS/CSS websites in a JSP-like way

We support the following features of native HTML websites:

- Grammar
    - [x][Introduction to Template Grammar](grammar.md)
    - [x][Import html file](import.md)
    - [x][Request context](context.md)
    - [x][Internal Interface Request](api.md)
    - [-]Access the database
- Configuration
    - [x][Custom Route](route.md)
    - [x][Built-in language internationalization](globalization.md)
- Pack
    - [x][Support gzip compression support](gzip.md)
    - [x][Support image webp automatic compression](webp.md)

# Get to know GTE quickly

In order to build a componentized website, we treat the `.html` file as a component. Components and components can refer to each other.

For example, we create a `header.html` file and a `footer.html` file in the project directory:

```html
<!-- header.html -->
<header>Header</header>
```

```html
<!-- footer.html -->
<footer>Footer</footer>
```

Then you can reference the above two components in the `index.html` file:
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

Then use the GTE command to run the project:
```shell
gte serve
```

You will see the running result in the browser:
```
Header
Content
Footer
```

# Next

- grammar
    - [Introduction to Template Grammar](grammar.md)
    - [Import html file](import.md)
    - [Request context](context.md)