<!DOCTYPE html>
<html>

<head>
    <title>Upload</title>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css" integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO"
        crossorigin="anonymous">
    <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo"
        crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js" integrity="sha384-ZMP7rVo3mIykV+2+9J3UJ46jBk0WLaUAdn689aCwoqbBJiSnjAK/l8WvCWPIPm49"
        crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js" integrity="sha384-ChfqqxuZUCnJSK3+MXmPNIyE6ZbWh2IMqE241rYiqJxyMiZ6OW/JmZQ5stwEULTy"
        crossorigin="anonymous"></script>
    <link rel="stylesheet" href="./css/style.css">
</head>

<body>
    <h1>Upload Page</h1>
    <main>
        <h3>YourName: {{.}}</h3>

        <hr>
        <form action="/upload" enctype="multipart/form-data" method="post" class="input-group mb-3">
            <div class="custom-file">
                <input type="file" name="upload" id="upload" multiple="multiple" class="custom-file-input" id="inputGroupFile02" required>
                <label class="custom-file-label" for="inputGroupFile02" aria-describedby="inputGroupFileAddon02"></label>
            </div>
            <div class="input-group-append">
                <input type="submit" class="input-group-text" id="inputGroupFileAddon02" value="Upload">
            </div>
        </form>
    </main>
</body>

</html>