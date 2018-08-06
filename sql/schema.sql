CREATE TABLE user_table
(
    usr_name CHAR(15) PRIMARY KEY,
    usr_password CHAR(20)  NOT NULL
);

CREATE TABLE post_table
(
    post_id CHAR(10) PRIMARY KEY,
    usr_name CHAR(15),
    img_name CHAR(100)
);

CREATE TABLE comment_table(
    comment_id CHAR(10) PRIMARY KEY,
    post_id CHAR(10),
    usr_name CHAR(15),
    comment_value CHAR(255)
);