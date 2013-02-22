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
                    <li {{if match .currentPage "^project_.*"}}class="active"{{end}}><a href="/?page={{.projectHashRingPage}}">Projects</a></li>
                {{end}}
                {{if allTrue .about_me .about_website}}
                    <li {{if match .currentPage "^about_.*"}}class="active"{{end}}><a href="/?page={{.aboutMePage}}">About</a></li>
                {{end}}
                {{if allTrue .admin_auth_code .admin_user_list}}
                    <li {{if match .currentPage "^admin_.*"}}class="active"{{end}}><a href="/?page={{.adminAuthCodePage}}">Admin</a></li>
                {{end}}
                </ul>
            </div>
            <ul class="nav pull-right">
                {{if .loginName}}
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

{{define "projects-navbar"}}
<div class="well sidebar-nav">
    <ul class="nav nav-list">
        <li class="nav-header">Projects</li>
    {{if allTrue .project_hash_ring}}
        <li {{if equal .currentPage .projectHashRingPage}}class="active"{{end}}><a href="/?page={{.projectHashRingPage}}">Hash Ring</a></li>
    {{end}}
    </ul>
</div>
{{end}}

{{define "about-navbar"}}
<div class="well sidebar-nav">
    <ul class="nav nav-list">
        <li class="nav-header">About</li>
        <li {{if equal .currentPage .aboutMePage}}class="active"{{end}}><a href="/?page={{.aboutMePage}}">About Me</a></li>
    {{if allTrue .about_website}}
        <li {{if equal .currentPage .aboutWebsitePage}}class="active"{{end}}><a href="/?page={{.aboutWebsitePage}}">About Website</a></li>
    {{end}}
    </ul>
</div>
{{end}}

{{define "admin-navbar"}}
<div class="well sidebar-nav">
    <ul class="nav nav-list">
        <li class="nav-header">Admin Board</li>
    {{if allTrue .admin_auth_code}}
        <li {{if equal .currentPage .adminAuthCodePage}}class="active"{{end}}><a href="/?page={{.adminAuthCodePage}}">Auth Code</a></li>
    {{end}}
     {{if allTrue .admin_user_list}}
        <li {{if equal .currentPage .adminUserListPage}}class="active"{{end}}><a href="/?page={{.adminUserListPage}}">User List</a></li>
    {{end}}
    </ul>
</div>
{{end}}