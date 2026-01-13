DELETE FROM external_user_identities;
DELETE FROM wish_reservations;
DELETE FROM wishlist_permissions;
DELETE FROM wishes;
DELETE FROM wishlists;
DELETE FROM users;

INSERT INTO users (id, name)
VALUES 
('alice', 'Alice'),
('bob', 'Bob');

INSERT INTO wishlists (id, user_id, name)
VALUES
('alice_wishlist_1', 'alice', 'Alice Wishlist 1'),
('alice_wishlist_2', 'alice', 'Alice Wishlist 2'),
('bob_wishlist_1', 'bob', 'Bob Wishlist 1'),
('bob_wishlist_2', 'bob', 'Bob Wishlist 2');

INSERT INTO wishlist_permissions(user_id, wishlist_id, permission_level)
VALUES
('alice', 'alice_wishlist_1', 'owner'),
('alice', 'alice_wishlist_2', 'owner'),
('bob', 'bob_wishlist_1', 'owner'),
('bob', 'bob_wishlist_2', 'owner');

INSERT INTO wishes (id, wishlist_id, name)
VALUES
('alice_wish_1', 'alice_wishlist_1', 'Alice Wish 1'),
('alice_wish_2', 'alice_wishlist_1', 'Alice Wish 2'),
('bob_wish_1', 'bob_wishlist_1', 'Bob Wish 1');


