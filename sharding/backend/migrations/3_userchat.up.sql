CREATE TABLE IF NOT EXISTS user_chat (
     id      VARCHAR(100) NOT NULL DEFAULT (uuid()) PRIMARY KEY,
     user_id VARCHAR(100) NOT NULL,
     chat_id VARCHAR(100) NOT NULL,
     FOREIGN KEY(user_id) REFERENCES user(id),
     FOREIGN KEY(chat_id) REFERENCES chat(id)
) ENGINE=InnoDB;
