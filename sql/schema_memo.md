# latestgramのデータベーステーブルのメモです。

## user_table
|カラム名|型|内容|
|---|---|---|
|id | int | ID|
|name  | char(15) | ユーザー名(PRIMARY KEY)|
|password|varchar(20) | パスワード |
|created_at | dateatime | 更新日時|

## post_table
|カラム名|型|内容|
|---|---|---|
|id | int | ID|
|user_id | int| ユーザーID|
|img_name |varchar(100) |画像の名前|
|created_at | dateatime | 更新日時|

## comment_table
|カラム名|型|内容|
|---|---|---|
|id | int | ID|
|post_id | int|post_tableから|
|usr_id | int|ユーザーID|
|comment_text |varchar(255) |コメントの内容|
|created_at | dateatime | 更新日時|
