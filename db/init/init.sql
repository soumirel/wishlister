CREATE TABLE IF NOT EXISTS users(
    id text NOT NULL,
    name text NOT NULL,

    PRIMARY KEY(id),

    UNIQUE(name),

    CHECK(char_length(name) BETWEEN 4 AND 32)
);

CREATE TABLE IF NOT EXISTS wishlists(
    id text NOT NULL,
    user_id text NOT NULL,
    name text NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,

    UNIQUE(user_id, name),

    CHECK(char_length(name) BETWEEN 4 AND 32)
);

CREATE TABLE IF NOT EXISTS wishes(
    id text NOT NULL,
    wishlist_id text NOT NULL,
    user_id text NOT NULL,
    name text NOT NULL,

    PRIMARY KEY(id),

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(wishlist_id) REFERENCES wishlists(id) ON DELETE CASCADE,

    UNIQUE(user_id, wishlist_id, name),

    CHECK(char_length(name) BETWEEN 4 AND 32)
);

-- CREATE TABLE IF NOT EXISTS reservations(
--     wish_id text NOT NULL,
--     user_id text NOT NULL,

--     PRIMARY KEY(wish_id, user_id)

--     FOREIGN KEY(wish_id) REFERENCES wishes(id)ON DELETE CASCADE,
--     FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
-- );

-- ACL
CREATE TABLE IF NOT EXISTS wishlists_permissions(
    id bigserial NOT NULL,
    user_id text NOT NULL,
    wishlist_id text NOT NULL,
    permission_level text NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(wishlist_id) REFERENCES wishlists(id) ON DELETE CASCADE,
    UNIQUE (user_id, wishlist_id)
);
