CREATE TABLE IF NOT EXISTS users(
    id text NOT NULL,
    name text NOT NULL,

    PRIMARY KEY(id),

    UNIQUE(name),

    CHECK(char_length(name) < 32)
);

CREATE TABLE IF NOT EXISTS wishlists(
    id text NOT NULL,
    user_id text NOT NULL,
    name text NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(name),

    CHECK(char_length(name) BETWEEN 4 AND 32)
);

CREATE TABLE IF NOT EXISTS wishes(
    id text NOT NULL,
    wishlist_id text NOT NULL,
    name text NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY(wishlist_id) REFERENCES wishlists(id) ON DELETE CASCADE,
    UNIQUE(wishlist_id, name),

    CHECK(char_length(name) BETWEEN 4 AND 32)
);

-- ACL
CREATE TABLE IF NOT EXISTS wishlist_permissions(
    id bigserial NOT NULL,
    user_id text NOT NULL,
    wishlist_id text NOT NULL,
    permission_level text NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(wishlist_id) REFERENCES wishlists(id) ON DELETE CASCADE,
    UNIQUE (user_id, wishlist_id)
);

CREATE TABLE IF NOT EXISTS wish_reservations(
    id text NOT NULL,
    wish_id text NOT NULL,
    reserved_by_user_id text NOT NULL,
    reserved_at timestamptz NOT NULL DEFAULT now(),

    PRIMARY KEY(id),
    FOREIGN KEY(wish_id) REFERENCES wishes(id) ON DELETE CASCADE,
    FOREIGN KEY(reserved_by_user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(wish_id)
)