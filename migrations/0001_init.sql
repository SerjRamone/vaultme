-- +goose Up
BEGIN;
    
-- user -----------------------
CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    login VARCHAR(155) UNIQUE NOT NULL,
    password VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE "user" IS 'Users';
  
COMMENT ON COLUMN "user".id IS 'Unique user ID';
COMMENT ON COLUMN "user".login IS 'User login';
COMMENT ON COLUMN "user".password IS 'User password';
COMMENT ON COLUMN "user".created_at IS 'Row created date';

-- item ----------------------
DROP TYPE IF EXISTS item_type;
CREATE TYPE item_type as enum ('CREDENTIAL', 'TEXT', 'RAW', 'CARD');

CREATE TABLE IF NOT EXISTS "item" (
    id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    user_id UUID NOT NULL REFERENCES "user" (id),
    name VARCHAR(255) NOT NULL,
    type item_type NOT NULL,
    version INT NOT NULL DEFAULT 1, 
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE "item" IS 'Store user\s secret data items';

COMMENT ON COLUMN "item".id IS 'Unique item ID';
COMMENT ON COLUMN "item".user_id IS 'User ID';
COMMENT ON COLUMN "item".name IS 'Item name';
COMMENT ON COLUMN "item".type IS 'Item type';
COMMENT ON COLUMN "item".version IS 'Item version';
COMMENT ON COLUMN "item".created_at IS 'Row created date';
COMMENT ON COLUMN "item".updated_at IS 'Row updated date';

-- meta ----------------------
CREATE TABLE IF NOT EXISTS "meta" (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    item_id UUID NOT NULL REFERENCES "item" (id),
    tag varchar(255) NOT NULL,
    text varchar(255) NOT NULL
);

COMMENT ON TABLE "meta" IS 'Store meta data';

COMMENT ON COLUMN "meta".id IS 'Unique meta ID';
COMMENT ON COLUMN "meta".item_id IS 'Item ID';
COMMENT ON COLUMN "meta".tag IS 'Meta tag';
COMMENT ON COLUMN "meta".text IS 'Meta text';

COMMIT;

-- item_data ----------------
CREATE TABLE IF NOT EXISTS "item_data" (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    item_id UUID NOT NULL REFERENCES "item" (id),
    data BYTEA NOT NULL
);

COMMENT ON TABLE "items_data" IS 'Store items data';

COMMENT ON COLUMN "items_data".id IS 'Unique item data ID';
COMMENT ON COLUMN "items_data".item_id IS 'Item ID';
COMMENT ON COLUMN "items_data".data IS 'Item binary data';

-- +goose Down
BEGIN;

-- item_data ----------------
DROP TABLE IF EXISTS "item_data" CASCADE;

-- meta ----------------------
DROP TABLE IF EXISTS "meta" CASCADE;

-- item ----------------------
DROP TABLE IF EXISTS "item" CASCADE;
DROP TYPE IF EXISTS item_type;

-- users ----------------------
DROP TABLE IF EXISTS "user" CASCADE;

COMMIT;

