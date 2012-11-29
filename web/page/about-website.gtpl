{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Mainpage - Hypermind</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="A page of Hypermind.">
    <meta name="author" content="hyper-carrot">

    {{template "header-import"}}

</head>
<body>

{{template "top-navbar" .}}

<div class="container-fluid">
    <div class="row-fluid">
        <div class="span2">
            <div class="well sidebar-nav">
                <ul class="nav nav-list">
                    <li class="nav-header">About</li>
                    <li><a href="/?page={{.aboutMePage}}">About Me</a></li>
                    <li class="active"><a href="/?{{.aboutWebsitePage}}">About Website</a></li>
                </ul>
            </div>
        </div>
        <div class="span10">
            <div class="hero-unit">
                <p>
                    This is a experimental website. It constructed by
                    <a class="btn btn-small" href="http://golang.org">Golang</a>
                    &
                    <a class="btn btn-small" href="http://twitter.github.com/bootstrap/">Bootstrap</a>
                    &
                    <a class="btn btn-small" href="http://redis.io/">Redis</a>
                    .
                </p>
                <p>
                    The Web site source code in here:
                    <a class="btn btn-small" href="https://github.com/hyper-carrot/go-web-demo">go-web-demo</a>
                    .
                </p>
            </div>
        </div>
    </div>
</div>

{{template "footer-import"}}

</body>
</html>
{{end}}
