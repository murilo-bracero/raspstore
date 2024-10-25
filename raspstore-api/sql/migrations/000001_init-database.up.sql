CREATE TABLE IF NOT EXISTS files (
    file_id text primary key,
    file_name text not null,
    size int not null,
    is_secret boolean not null default FALSE,
    owner_id text not null,
    created_at int not null,
    updated_at int,
    created_by text not null,
    updated_by text
);

CREATE UNIQUE INDEX IF NOT EXISTS files_file_id_idx ON files (file_id);

CREATE TABLE IF NOT EXISTS files_permissions (
    permission_id text primary key,
    file_id text not null,
    permission text check(permission in ('EDITOR', 'VIEWER')) not null,
    user_id text not null,
    FOREIGN KEY(file_id) REFERENCES files(file_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS files_permissions_permission_id_idx ON files_permissions (permission_id);