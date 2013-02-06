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
                    <li><a href="/?page={{.aboutWebsitePage}}">About Website</a></li>
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
                    <button id="cv_submit" type="submit" class="btn">Submit</button>
                </p>
            </div>
        </div>
    </div>
</div>

{{template "footer-import"}}

<script>
$(function () {
  $('#cv_submit').bind('click', function() {
    var auth_code = $('#cv_auth_code').val();
    //alert("Input: " + auth_code);
    $.post("/get-cv", "auth_code=" + auth_code,
        function (data, textStatus){
            if ((data.indexOf("FAIL:") == 0) || (data.indexOf("ERROR:") == 0)) {
                alert(data);  
            } else {
                openWindow = window.open("", "", "height=600, width=800,top=50,left=50,toolbar=no,menubar=no,scrollbars=auto,resizeable=no,location=no,status=no");  
                openWindow.document.write(data)
                openWindow.document.close(); 
            }
        }, "text");
    });
})
</script>

</body>
</html>
{{end}}
