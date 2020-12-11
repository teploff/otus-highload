ALTER TABLE message
    ADD COLUMN shard_id INT NOT NULL AFTER chat_id;