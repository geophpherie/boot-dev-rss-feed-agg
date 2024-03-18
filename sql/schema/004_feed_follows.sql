-- +goose Up
CREATE TABLE FeedFollows(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL references Feeds(id) ON DELETE CASCADE,
    user_id UUID NOT NULL references Users(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE FeedFollows;