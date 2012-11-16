<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Mainpage - Hypermind</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="The homepage of Hypermind.">
    <meta name="author" content="hyper-carrot">

    <link href="../css/bootstrap.css" rel="stylesheet">
    <style>
        body {
        padding-top: 60px;
        }
    </style>
    <link href="../css/bootstrap-responsive.css" rel="stylesheet">

    <!-- HTML5 shim, for IE6-8 support of HTML5 elements -->
    <!--[if lt IE 9]>
    <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
    <![endif]-->

    <link rel="shortcut icon" href="../img/favicon.ico">
    <link rel="apple-touch-icon-precomposed" sizes="144x144" href="../img/apple-touch-icon-144-precomposed.png">
    <link rel="apple-touch-icon-precomposed" sizes="114x114" href="../img/apple-touch-icon-114-precomposed.png">
    <link rel="apple-touch-icon-precomposed" sizes="72x72" href="../img/apple-touch-icon-72-precomposed.png">
    <link rel="apple-touch-icon-precomposed" href="../img/apple-touch-icon-57-precomposed.png">
</head>
<body>
<div class="navbar navbar-fixed-top">
    <div class="navbar-inner">
        <div class="container-fluid">
            <a class="btn btn-navbar" data-toggle="collapse" data-target=".nav-collapse">
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </a>
            <a class="brand" href="#">Hypermind</a>
            <div class="nav-collapse collapse">
                <ul class="nav">
                    <li class="active"><a href="/?page={{.homePage}}">Home</a></li>
                    <li><a href="/?page={{.meetingKanbanPage}}">Meeting Kanban</a></li>
                    <li><a href="/?page={{.aboutMePage}}">About</a></li>
                </ul>
            </div>
            <ul class="nav pull-right">
                {{if .validLogin}}
                <li><a href="#">Hi, {{.loginName}}</a></li>
                <li class="divider-vertical"></li>
                <a class="btn navbar-form pull-right" href="http://{{.serverAddr}}:{{.serverPort}}/logout">Sign Out</a></p>
                {{else}}
                <li class="divider-vertical"></li>
                <a class="btn navbar-form pull-right" href="http://{{.serverAddr}}:{{.serverPort}}/login">I'm Admin</a></p>
                {{end}}
            </ul>

        </div>
    </div>
</div>

<div class="container">
    <div class="hero-unit">
      <h2>Welcome to Hypermind!</h2>
      <p>
        The web site is mainly in order to provide more convenience to myself (or people who have same needs with me).
        The function including meeting kanban etc.
      </p>
      <p>
        <a class="btn btn-primary btn-large" href="/?page={{.aboutWebsitePage}}">It's how to be constructed »</a>
      </p>
    </div>
  </div>
</div>

<script src="../js/jquery.js"></script>
<script src="../js/bootstrap.js"></script>

</body>
</html>

