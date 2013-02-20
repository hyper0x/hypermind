{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Mainpage - Hypermind</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="The homepage of Hypermind.">
    <meta name="author" content="hyper-carrot">

    {{template "header-import"}}

</head>
<body>

{{template "top-navbar" .}}

<div class="container">
    <div class="hero-unit">
      <h2>Welcome to Hypermind!</h2>
      <p>
        The web site is mainly in order to provide more convenience to myself (or people who have same needs with me).
        <br>
        You can find out some open source projects create by me in <a href="/?page={{.projectHashRingPage}}">here</a>.
        <br>
        If you would know the author of this website, please try to click <a href="/?page={{.aboutMePage}}">here</a>.
      </p>
      <p>
        
      </p>
    </div>
  </div>
</div>

{{template "footer-import"}}

</body>
</html>
{{end}}

