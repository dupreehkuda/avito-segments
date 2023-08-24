CREATE TABLE IF NOT EXISTS users (
                                     id text PRIMARY KEY NOT NULL UNIQUE,
                                     created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS segments (
                                        tag text PRIMARY KEY NOT NULL UNIQUE,
                                        created_at timestamptz NOT NULL,
                                        deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS user_segments (
                                             tag text,
                                             user_id text,
                                             created_at timestamptz NOT NULL,
                                             expired_at timestamptz,
                                             deleted_at timestamptz,
                                             PRIMARY KEY(tag, user_id)
);

ALTER TABLE user_segments ADD FOREIGN KEY (tag) REFERENCES segments (tag);
ALTER TABLE user_segments ADD FOREIGN KEY (user_id) REFERENCES users (id);