CREATE TABLE IF NOT EXISTS users (
                                     id text PRIMARY KEY NOT NULL UNIQUE,
                                     created_at timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS segments (
                                        tag text PRIMARY KEY NOT NULL UNIQUE,
                                        description text,
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

CREATE OR REPLACE  FUNCTION insert_user_if_not_exists() RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM users WHERE id = NEW.user_id) THEN
        INSERT INTO users (id, created_at)
        VALUES (NEW.user_id, now());
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER insert_user
    BEFORE INSERT ON user_segments
    FOR EACH ROW
EXECUTE FUNCTION insert_user_if_not_exists();

