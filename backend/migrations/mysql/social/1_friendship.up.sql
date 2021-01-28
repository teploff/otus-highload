CREATE TABLE IF NOT EXISTS friendship (
    id             VARCHAR(100) NOT NULL DEFAULT (uuid()),
    master_user_id VARCHAR(100) NOT NULL,
    slave_user_id  VARCHAR(100) NOT NULL,
    status         VARCHAR(100) NOT NULL CONSTRAINT friendship_status_field CHECK (status in ('expected', 'accepted')),
    create_time    TIMESTAMP NOT NULL,
    update_time    TIMESTAMP,
    CONSTRAINT     uc_master_slave_ids UNIQUE(master_user_id, slave_user_id),
) ENGINE=InnoDB;
