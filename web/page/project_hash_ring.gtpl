{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Hash Ring - Hypermind</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="The homepage of Hypermind.">
    <meta name="author" content="hyper-carrot">

    {{template "header-import"}}
    {{template "js-import"}}

</head>
<body>

{{template "top-navbar" .}}

<div class="container-fluid">
    <div class="row-fluid">
        <div class="span2">
            {{template "projects-navbar" .}}
        </div>
        <div class="span10">
            <div class="hero-unit">
                <p>
                    Hash ring is a kind of consistency hash realization. I had implemented hash ring model written by Go, Python and Java.
                    And, I compared their performance. (see picture below)
                </p>
                <p>
                    <img src="/img/chash_benchmark2.png">
                </p>
                <p>
                    The source code of hash ring implemented by me is here: <br>
                    Go Edition: <a href="https://github.com/hyper-carrot/chash4go">chash4go</a><br>
                    Python Edition: <a href="https://github.com/hyper-carrot/chash4py">chash4py</a><br>
                    Java Edition: <a href="https://github.com/hyper-carrot/chash4j">chash4j</a><br>
                </p>
            </div>
        </div>
    </div>
</div>

</body>
</html>
{{end}}

