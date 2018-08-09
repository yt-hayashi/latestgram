<!DOCTYPE html>
<html>
<head>
<title>Login</title>
<meta charset="utf-8">
</head>
<body>
<h1>Login Page</h1>
<hr>
<h3>{{.}}</h3>
<form action="/login" method="post">
    <p>Username</p>
    <p class="username">
     <input type="text" name="username" maxlength="32" autocomplete="OFF" />
    </p>
    <p>Password</p>
    <p class="password">
        <input type="password" name="password" maxlength="32" autocomplete="OFF" />
    </p>
    <p class="submit">
        <input type="submit" value="Login" />
    </p>
</form>

</body>
</html>