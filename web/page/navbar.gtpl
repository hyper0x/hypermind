{{define "top-navbar"}}
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
                    <li {{if equal .currentPage .homePage}}class="active"{{end}}><a href="/?page={{.homePage}}">Home</a></li>
                    <li {{if equal .currentPage .meetingKanbanPage}}class="active"{{end}}><a href="/?page={{.meetingKanbanPage}}">Meeting Kanban</a></li>
                    <li {{if equal .currentPage .hashRingPage}}class="active"{{end}}><a href="/?page={{.hashRingPage}}">Hash Ring</a></li>
                    <li {{if match .currentPage "^about-.*"}}class="active"{{end}}><a href="/?page={{.aboutMePage}}">About</a></li>
                </ul>
            </div>
            <ul class="nav pull-right">
                {{if .validLogin}}
                <li><a href="#">Hi, {{.loginName}}</a></li>
                <li class="divider-vertical"></li>
                <a class="btn navbar-form pull-right" href="http://{{.serverAddr}}:{{.serverPort}}/logout">Sign Out</a></p>
                {{else}}
                <li class="divider-vertical"></li>
                <a class="btn navbar-form pull-right" href="http://{{.serverAddr}}:{{.serverPort}}/login">Sign In</a></p>
                {{end}}
            </ul>

        </div>
    </div>
</div>
{{end}}
