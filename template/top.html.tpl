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
<input type="button" onclick="location.href='/signup'"value="SignUp">
<input type="button" onclick="location.href='/login'"value="Login">
<input type="button" onclick="location.href='/upload'"value="Upload">

<hr>
<div>
    {{range .}}
    <div class="post">
        <img src={{.ImgPath}} alt="post_img">
        <h4>{{.NameText}}</h4>
    </div>
    {{end}}
</div>

</body>
</html>