<!DOCTYPE html>
<html>
<head>
    <title>Upload</title>
    <meta charset="utf-8">
</head>
<body>
<h1>Upload Page</h1>
<<<<<<< HEAD
<h3>YourName{{.}}</h3>
=======
<h3>YourName: {{.}}</h3>
>>>>>>> 085fbfc4a9553589741e8cd81a6963079dabd166

<hr>
    <form action="/upload" enctype="multipart/form-data" method="post">
        <input type="file" name="upload" id="upload" multiple="multiple">
        <input type="submit" value="Upload image" />
    </form>
</body>
</html>