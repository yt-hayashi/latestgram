<!DOCTYPE html>
<html>
<head>
    <title>Upload</title>
    <meta charset="utf-8">
</head>
<body>
<h1>Upload Page</h1>
<h3>YourName: {{.}}</h3>

<hr>
    <form action="/upload" enctype="multipart/form-data" method="post">
        <input type="file" name="upload" id="upload" multiple="multiple">
        <input type="submit" value="Upload image" />
    </form>
</body>
</html>