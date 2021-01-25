CREATE TABLE IF NOT EXISTS news (
    id          VARCHAR(100) NOT NULL DEFAULT (uuid()),
    owner_id    VARCHAR(100) NOT NULL,
    content     VARCHAR(500) NOT NULL,
    create_time TIMESTAMP NOT NULL,
    FOREIGN     KEY(owner_id) REFERENCES user(id)
) ENGINE=InnoDB;
