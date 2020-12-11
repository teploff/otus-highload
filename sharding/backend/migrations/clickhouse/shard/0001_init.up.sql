CREATE TABLE message (
    date Date,
    time DateTime,
    event String,
    client String,
    value UInt32,
    nshard UInt8
) ENGINE = MergeTree(date, (event, client), 8192)