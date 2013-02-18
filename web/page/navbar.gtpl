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
                {{if allTrue .home}}
                    <li {{if equal .currentPage .homePage}}class="active"{{end}}><a href="/?page={{.homePage}}">Home</a></li>
                {{end}}
                {{if allTrue .meeting_kanban}}
                    <li {{if equal .currentPage .meetingKanbanPage}}class="active"{{end}}><a href="/?page={{.meetingKanbanPage}}">Meeting Kanban</a></li>
                {{end}}
                {{if allTrue .project_hash_ring}}
                    <li {{if match .currentPage "^project-.*"}}class="active"{{end}}><a href="/?page={{.projectHashRingPage}}">Projects</a></li>
                {{end}}
                {{if allTrue .about_me .about_website}}
                    <li {{if match .currentPage "^about-.*"}}class="active"{{end}}><a href="/?page={{.aboutMePage}}">About</a></li>
                {{end}} 
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
