{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>About Website - Hypermind</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="A page of Hypermind.">
    <meta name="author" content="hyper-carrot">

    {{template "header-import"}}
    {{template "js-import"}}

</head>
<body>

{{template "top-navbar" .}}

<div class="container-fluid">
    <div class="row-fluid">
        <div class="span2">
            {{template "about-navbar" .}}
        </div>
        <div class="span10">
            <div class="hero-unit">
                <p>
                    This website constructed by
                    <a class="btn btn-small" href="http://golang.org" target="_blank">Golang</a>
                    &
                    <a class="btn btn-small" href="http://twitter.github.com/bootstrap/" target="_blank">Bootstrap</a>
                    &
                    <a class="btn btn-small" href="http://redis.io/" target="_blank">Redis</a>
                    .
                </p>
                <p>
                    <a href="https://github.com/hyper-carrot/hypermind" target="_blank" class="btn btn-primary">Source code</a>
                </p>
            </div>
        </div>
    </div>
</div>

</body>
</html>
{{end}}
