-- +goose Up

DROP TRIGGER IF EXISTS update_user_updated_at ON userinfo;

DROP FUNCTION IF EXISTS set_updated_at_timestamp();

CREATE OR REPLACE FUNCTION set_updated_at_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON userinfo
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_timestamp();

-- +goose Down

DROP TRIGGER IF EXISTS update_user_updated_at ON userinfo;

DROP FUNCTION IF EXISTS set_updated_at_timestamp() CASCADE;
