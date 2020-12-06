CREATE TABLE IF NOT EXISTS message (
    id          VARCHAR(100) NOT NULL DEFAULT (uuid()),
    text        TEXT NOT NULL,
    status      INT NOT NULL CONSTRAINT message_status_field CHECK (status in (0, 1, 2, 3, 4)),
    create_time TIMESTAMP NOT NULL,
    user_id     VARCHAR(100) NOT NULL,
    chat_id     VARCHAR(100) NOT NULL,
    CONSTRAINT  uc_id_status_ui_ci UNIQUE(id, status, user_id, chat_id),
    FOREIGN     KEY(user_id) REFERENCES user(id),
    FOREIGN     KEY(chat_id) REFERENCES chat(id)
) ENGINE=InnoDB;
