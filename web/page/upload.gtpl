{{define "page"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Upload file</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post"> 
  <input type="file" name="file" /> 
  <input type="hidden" name="token" value="{{.}}">
  <input type="submit" value="upload" /> 
</form>
</body>
</html>
{{end}}
