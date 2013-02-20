{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Meeting Kanban - Hypermind</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="The homepage of Hypermind.">
    <meta name="author" content="hyper-carrot">

    {{template "header-import"}}

</head>
<body>

{{template "top-navbar" .}}

<div class="container">
    <div class="hero-unit">
        <h2>Come soon ...</h2>
        <p>
        <br>
        <h4>The scheduled functions:</h4>
        <ul>
            <li>Publish meeting</li>
            <li>Manage own meeting</li>
            <li>Meeting list show</li>
            <li>Meeting detail show</li>
            <li>Register for meeting</li>
            <li>* Registration confirmation & Reminding</li>
            <li>* Meeting tag & aggregation & recommendation</li>
            <li>* Meeting static/vertical analysis</li>
        </ul>
        <b>* Vision and long-term goals</b>
        </p>
    </div>
</div>
</div>

{{template "footer-import"}}

</body>
</html>
{{end}}
