/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	full_name VARCHAR (60) NOT NULL,
	phone_number VARCHAR (20) UNIQUE NOT NULL,
	login_count INT DEFAULT 0,
	password TEXT NOT NULL,
	is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NULL
);

-- function for update updated_at
CREATE FUNCTION update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

-- create triggers
CREATE TRIGGER users_updated_at BEFORE
UPDATE
    ON
    users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


INSERT INTO users (full_name, phone_number, password) VALUES ('john doe', '+62812345677', crypt('test1', gen_salt('bf')));
INSERT INTO users (full_name, phone_number, password) VALUES ('mamang doe', '+62812345678', crypt('test2', gen_salt('bf')));
