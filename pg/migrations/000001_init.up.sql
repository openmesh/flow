CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users
(
    id            UUID                 DEFAULT uuid_generate_v4()
        CONSTRAINT users_pkey
            PRIMARY KEY,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name          VARCHAR     NULL,
    email         VARCHAR     NULL,
    password_hash VARCHAR     NULL,
    api_key       VARCHAR     NOT NULL
);

CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();

CREATE TABLE auths
(
    id            UUID                 DEFAULT uuid_generate_v4()
        CONSTRAINT auths_pkey
            PRIMARY KEY,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id       UUID
        CONSTRAINT auths_users_user
            REFERENCES users
            ON DELETE CASCADE,
    source        VARCHAR     NOT NULL,
    source_id     VARCHAR     NULL,
    access_token  VARCHAR     NULL,
    refresh_token VARCHAR     NULL,
    expires_at    TIMESTAMPTZ NULL
);

CREATE TRIGGER auths_set_updated_at
    BEFORE UPDATE
    ON auths
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();

CREATE TABLE workflows
(
    id          UUID                 DEFAULT uuid_generate_v4()
        CONSTRAINT workflows_pkey
            PRIMARY KEY,
    user_id     UUID
        CONSTRAINT workflows_users_user
            REFERENCES users
            ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name        VARCHAR     NULL,
    description VARCHAR     NULL
);

CREATE TRIGGER workflows_set_updated_at
    BEFORE UPDATE
    ON workflows
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();

CREATE TABLE nodes
(
    id          UUID                 DEFAULT uuid_generate_v4()
        CONSTRAINT nodes_pkey
            PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    workflow_id UUID
        CONSTRAINT nodes_workflows_workflow
            REFERENCES workflows
            ON DELETE CASCADE,
    integration VARCHAR     NOT NULL,
    action      VARCHAR     NOT NULL
);

CREATE TRIGGER nodes_set_updated_at
    BEFORE UPDATE
    ON nodes
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();

CREATE TABLE edges
(
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    head_id    UUID                 DEFAULT uuid_generate_v4()
        CONSTRAINT edges_nodes_head_id
            REFERENCES nodes
            ON DELETE CASCADE,
    tail_id    UUID                 DEFAULT uuid_generate_v4()
        CONSTRAINT edges_nodes_tail_id
            REFERENCES nodes
            ON DELETE CASCADE,
    CONSTRAINT edges_pkey
        PRIMARY KEY (head_id, tail_id)
);

CREATE TABLE params
(
    id         UUID                 DEFAULT uuid_generate_v4()
        CONSTRAINT params_pkey
            PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    key        VARCHAR     NOT NULL,
    value      VARCHAR     NOT NULL,
    type       VARCHAR     NOT NULL
);

CREATE TRIGGER params_set_up_updated_at
    BEFORE UPDATE
    ON params
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();