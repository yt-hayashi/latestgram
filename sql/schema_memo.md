# latestgramのデータベーステーブルのメモです。

## user_table
|レコード名|型|内容|
|---|---|---|
|usr_name  | char(15) | ユーザー名(PRIMARY KEY)|
|usr_password|char(20) | パスワード |

## post_table
|レコード名|型|内容|
|---|---|---|
|post_id  |char(10)|ランダムに割り当て(PRIMARY KEY)|
|usr_name | char(15)| ユーザー名|
|img_name |char(100) |画像の名前|

## comment_table
|レコード名|型|内容|
|---|---|---|
|comment_id  | char(10) | ランダムに割り当て(PRYMARY kEY) |
|post_id | CHAR(10)|post_tableから|
|usr_name |char(15) |ユーザー名|
|comment_value |char(255) |コメントの内容|
