-- +goose Up
CREATE TABLE channels (
    id INTEGER PRIMARY KEY,
    type TEXT CHECK( type IN ('DM','GROUP','SCHAN') ),
    name TEXT,
    created_at TIMESTAMP,
    image BLOB
);

CREATE TABLE channelaccess (
    userID INTEGER NOT NULL,
    channelID INTEGER NOT NULL,
    CONSTRAINT FK_userID FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT FK_channelID FOREIGN KEY (channelID) REFERENCES channels(id) ON DELETE CASCADE,
    PRIMARY KEY (userID, channelID)
);

CREATE TABLE messages (
    id INTEGER PRIMARY KEY,
    senderID INTEGER NOT NULL,
    channelID INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT FK_senderID FOREIGN KEY (senderID) REFERENCES users(id),
    CONSTRAINT FK_channelID FOREIGN KEY (channelID) REFERENCES channels(id)
);

-- +goose Down
DROP TABLE channels;
DROP TABLE channelaccess;
DROP TABLE messages;