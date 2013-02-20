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
        The function including meeting kanban etc.
      </p>
      <p>
        <a class="btn btn-primary btn-large" href="/?page={{.aboutWebsitePage}}">Detail »</a>
      </p>
    </div>
  </div>
</div>

{{template "footer-import"}}

</body>
</html>
{{end}}

