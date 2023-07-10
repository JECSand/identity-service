CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS user_groups CASCADE;
DROP TABLE IF EXISTS memberships CASCADE;
DROP TABLE IF EXISTS blacklists CASCADE;


CREATE TABLE users
(
    id              UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    username        VARCHAR(250)  NOT NULL CHECK ( username <> '' ),
    email           VARCHAR(250)  NOT NULL CHECK ( email <> '' ),
    password        VARCHAR(250) NOT NULL CHECK ( password <> '' ),
    root            BOOLEAN       NOT NULL,
    active          BOOLEAN       NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_groups
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_name  VARCHAR(250)  NOT NULL CHECK ( group_name <> '' ),
    description VARCHAR(250) NOT NULL CHECK ( description <> '' ),
    creator_id  UUID NOT NULL,
    active      BOOLEAN       NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

CREATE TABLE memberships
(
    id          UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL,
    group_id    UUID NOT NULL,
    status      INTEGER       NOT NULL,
    member_role INTEGER       NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES user_groups(id)
);

CREATE TABLE blacklists
(
    id                 UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    access_token       VARCHAR(2500)  NOT NULL CHECK ( access_token <> '' ),
    created_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);