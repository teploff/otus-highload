ALTER TABLE message
    ADD COLUMN shard_key_id INT NOT NULL AFTER chat_id;