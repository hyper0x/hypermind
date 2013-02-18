{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>About Me - Hypermind</title>
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
                    <li class="active"><a href="/?page={{.aboutMePage}}">About Me</a></li>
                {{if allTrue .aboutWebsitePage}}
                    <li><a href="/?page={{.aboutWebsitePage}}">About Website</a></li>
                {{end}}
                </ul>
            </div>
        </div>
        <div class="span10">
            <div class="hero-unit">
                <p>
                    My name is Harry Hao.
                    I live in Beijing.
                    I am in Sohu Inc (NSDQ: SOHU) as the position of Dev Leader.
                </p>
                <p>
                    I'm a broad interests software developer. I'm a open source fan, and pay attention to the agile methods and software process improvement.
                    I focus on Clojure and Go programming language, and contribute strength in order to the popularization of them in Chinese community.
                </p>
                <p>
                    My homepage in GitHub is
                    <a class="btn btn-small" href="https://github.com/hyper-carrot">hyper-carrot</a>
                    .
                </p>
                <p>
                    <label><h3>CV Authorization Code</h3></label>
                    <input id="cv_auth_code"type="text" class="input-medium search-query" placeholder="Please inquire with me.">
                    <button id="cv_auth_code_submit" type="submit" class="btn">Submit</button>
                </p>
            </div>
        </div>
    </div>
</div>

{{template "footer-import"}}

<script src="../js/hypermind.js"></script>

</body>
</html>
{{end}}
