CREATE TABLE IF NOT EXISTS friendship (
    id           VARCHAR(100) NOT NULL DEFAULT (uuid()),
    subj_user_id VARCHAR(100) NOT NULL,
    obj_user_id  VARCHAR(100) NOT NULL,
    status       INT NOT NULL CONSTRAINT friendship_status_field CHECK (status in (0, 1, 2)),
    create_time  TIMESTAMP NOT NULL,
    update_time  TIMESTAMP,
    CONSTRAINT   uc_id_status_ui_ci UNIQUE(subj_user_id, obj_user_id),
    FOREIGN      KEY(subj_user_id) REFERENCES user(id),
    FOREIGN      KEY(obj_user_id) REFERENCES chat(id)
) ENGINE=InnoDB;
