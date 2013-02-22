{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>User list board for admin - Hypermind</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="The homepage of Hypermind.">
    <meta name="author" content="hyper-carrot">

    {{template "header-import"}}
    {{template "js-import"}}

<script type="text/javascript">
$(document).ready(function() {
    url = "/user_list"
    $.ajax({ url: url, success: function(data) {
        var userTbody = $("#user_list_table tbody")
        if(!data){
            return 
        }
        $.each(data, function(i, user){
            var addRowTemplete = "<tr>" + 
                "<td id='" + (user.LoginName + "_login_name") + "'>" + user.LoginName + "</td>" +
                "<td id='" + (user.LoginName + "_email") + "'>" + user.Email + "</td>" + 
                "<td id='" + (user.LoginName + "_mobile_phone") + "'>" + user.MobilePhone + "</td>" + 
                "<td id='" + (user.LoginName + "_group") + "'>" + user.Group + "</td>" +
                "<td id='" + (user.LoginName + "_remark") + "'>" + user.Remark + "</td>" + 
                "</tr>";
            userTbody.append(addRowTemplete);
        });
        var innerBoard = $("#inner_board")
        for (i = 0; i < data.length; i++) {
            innerBoard.append("<br>")
        }  
    }, dataType: "json", timeout: 2000 });
 });
</script>

</head>
<body>

{{template "top-navbar" .}}

<div class="container-fluid">
    <div class="row-fluid">
        <div class="span2">
            {{template "admin-navbar" .}}
        </div>
        <div class="span10">
            <div id="inner_board" class="hero-unit">
                <table id="user_list_table" class="table table-striped table-bordered table-condensed span10">
                    <thead>
                      <tr>
                        <th class="span2">Login Name</th>
                        <th class="span3">Email</th>
                        <th class="span3">Mobile Phone</th>
                        <th class="span2">Group</th>
                        <th class="span3">Remark</th>
                      </tr>
                    </thead>
                    <tbody>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</div>

</body>
</html>
{{end}}
