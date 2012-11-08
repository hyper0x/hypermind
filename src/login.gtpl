<html>
<head>
<title>Login form demo By Go</title>
</head>
<body>
<form action="http://{{.serverAddr}}:{{.serverPort}}/login" method="post">
    Username: <input type="text" name="login_name">
    Passport: <input type="password" name="password">
    <input type="hidden" name="token" value="{{.token}}">
    <input type="submit" value="Login">
</form>
</body>
</html>

