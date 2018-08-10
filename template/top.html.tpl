<!DOCTYPE html>
<html>

<head>
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Top</title>
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
    <h1>Latestgram</h1>
    <main>
        <h2>Top Page</h2>
        <h3>Post</h3>
        <input class="btn btn-primary btn-lg" type="button" onclick="location.href='/signup'" value="SignUp">
        <input class="btn btn-primary btn-lg" type="button" onclick="location.href='/login'" value="Login">
        <input class="btn btn-success btn-lg" type="button" onclick="location.href='/upload'" value="Upload">
        <input class="btn btn-danger btn-lg" type="button" onclick="location.href='/logout'" value="Logout">
        <hr>
        <div class="card-columns">
            {{range .}}
            <div class="card" style="width: 25rem;">
                <img class="card-img-top" src={{.ImgPath}} alt="post_img">
                <div class="card-body">
                    <h5 class="card-title">Post User: {{.NameText}}</h5>
                    <p class="card-text">Comments</p>
                    <ul class="list-group list-group-flush">
                        {{range $var := .Comments}}
                        <li class="list-group-item">{{$var}}</li>
                        {{end}}
                    </ul>

                    <form action="/comment?id={{.PostID}}" method="post" class="input-group mb-3">
                        <input type="text" name="comment_text" class="form-control" placeholder="Add comment..." aria-label="Add comment..." aria-describedby="button-addon2" autocomplete="OFF" required>
                        <div class="input-group-append">
                            <button class="btn btn-outline-secondary" type="submit" id="button-addon2">Post</button>
                        </div>

                    </form>

                </div>
            </div>
            {{end}}
    </main>Ω
</body>

</html>