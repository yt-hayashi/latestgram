<!DOCTYPE html>
<html>
<head>
<title>SignUp</title>
<meta charset="utf-8">
</head>
<body>
<h1>SignUp Page</h1>
<hr>
<h3>{{.}}</h3>
<form action="/signup"　mthod="post">
    <p>Username</p>
    <p class="username">
     <input type="text" name="username" maxlength="32" autocomplete="OFF" />
    </p>
    <p>Password</p>
    <p class="password">
        <input type="password" name="password" maxlength="32" autocomplete="OFF" />
    </p>
    <p class="submit">
        <input type="submit" value="SignUp" />
    </p>
</form>

</body>
</html>