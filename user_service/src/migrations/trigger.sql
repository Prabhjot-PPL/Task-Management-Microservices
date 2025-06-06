-- function that will be called by the trigger
CREATE OR REPLACE FUNCTION set_updated_at_timestamp()
RETURNS TRIGGER AS $$ 
BEGIN 
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- drop the trigger
DROP TRIGGER IF EXISTS update_user_updated_at ON user_details;

-- creating trigger to update user_details "updated_at" field
CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON user_details
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_timestamp();