<!DOCTYPE html>
<html>
<head>
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<title>Top</title>
<meta charset="utf-8">
</head>
<body>
<h1>Top Page</h1>
<h3>Post</h3>

<hr>
<div>
    {{range .}}
    <h4>{{.NameText}}</h4>
    <p>IMG={{.ImgPath}}</p>
    {{end}}
</div>

</body>
</html>