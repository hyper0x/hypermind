{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>CV for admin - Hypermind</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="The homepage of Hypermind.">
    <meta name="author" content="hyper-carrot">

    {{template "header-import"}}
    {{template "js-import"}}

<script type="text/javascript">
$(document).ready(function() {
    url = "/auth_code"
    count = 0
    function poll_auth_code() {
        $.ajax({ url: url, success: function(data) {
            if (count == 0) {
                $("#initial").text(data);
                url += "?type=lp"
            }
            $("#current").text(data);
            $("#count").text(count);
            count++
        }, dataType: "text", complete: poll_auth_code, timeout: (1000 * 60 * 10) });
    }
    poll_auth_code()
 });
</script>

</head>
<body>

{{template "top-navbar" .}}

<div class="container-fluid">
    <div class="row-fluid">
        <div class="span2">
            <div class="well sidebar-nav">
                <ul class="nav nav-list">
                    <li class="nav-header">About</li>
                {{if allTrue .admin_auth_code}}
                    <li {{if equal .currentPage .adminAuthCodePage}}class="active"{{end}}><a href="/?page={{.adminAuthCodePage}}">Auth Code</a></li>
                {{end}}
                </ul>
            </div>
        </div>
        <div class="span10">
            <div class="hero-unit">
                <p>
                    <table class="table table-striped table-bordered table-condensed span5">
                        <thead>
                          <tr>
                            <th class="span3">Initial</th>
                            <th class="span3">Current</th>
                            <th class="span2">Count</th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr>
                            <td><span id="initial">......</span></td>
                            <td><span id="current">......</span></td>
                            <td><span id="count">......</span></td>
                          </tr>
                        </tbody>
                    </table>
                </p>
                <br>
                <br>
                <br>
                <h4>Title Description:</h4>
                <dl class="dl-horizontal">
                <dt>Initial:</dt>
                <dd>The value of auth code when this page show.</dd>
                <dt>Current:</dt>
                <dd>The current value of auth code.</dd>
                <dt>Count:</dt>
                <dd>The change count of auth code since this page show.</dd>
                </ul>
                </dl>
            </div>
        </div>
    </div>
</div>}

</body>
</html>
{{end}}
