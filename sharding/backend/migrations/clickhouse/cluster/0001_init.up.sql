CREATE TABLE message (
    datetime    DateTime,
    id          UUID,
    text        String,
    status      Enum8('created' = 0, 'received' = 1, 'read' = 2, 'deleted' = 3),
    create_time Int64,
    user_id     UUID,
    chat_id     UUID,
    shard_id    UInt8
) ENGINE = Distributed(messages, default, message, shard_id);

CREATE TABLE chat (
    datetime    DateTime,
    id          UUID,
    create_time Int64
) ENGINE = MergeTree()
ORDER BY (datetime, id, create_time);