CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE IF NOT EXISTS members_permissions (
    member_id bigint NOT NULL REFERENCES members ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (member_id, permission_id)
);


-- Add the two permissions to the table.
INSERT INTO permissions (code)
VALUES
    ('posts:read'),
    ('posts:write');