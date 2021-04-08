-- Create a new UTF-8 `snippetbox` database.
-- CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Switch to using the `snippetbox` database.
-- USE snippetbox;

-- CREATE TABLE snippets (
-- 	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
-- 	title VARCHAR(100) NOT NULL,
--     content TEXT NOT NULL,
--     created DATETIME NOT NULL,
--     expires DATETIME NOT NULL
-- );

-- Add an index on the created column
-- CREATE INDEX idx_snippets_created ON snippets(created);

-- Add some dummy records
-- INSERT INTO snippets (title, content, created, expires) VALUES (
-- 	'An old silent pond',
--     'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
--     UTC_TIMESTAMP(), -- Get current UTC timestamp
--     DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY) -- Expires in one year
-- );

-- INSERT INTO snippets (title, content, created, expires) VALUES (
-- 	'Over the wintry forest',
--     'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
--     UTC_TIMESTAMP(),
--     DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
-- );

-- INSERT INTO snippets (title, content, created, expires) VALUES (
-- 	'First autumn morning',
--     'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
--     UTC_TIMESTAMP(),
--     DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
-- );

-- Create a new dummy user with restricted permissions (SELECT, INSERT, UPDATE)
-- CREATE USER 'web'@'localhost';
-- GRANT SELECT, INSERT, UPDATE ON snippetbox.* TO 'web'@'localhost';

-- Set a dummy password for user web
-- ALTER USER 'web'@'localhost' IDENTIFIED BY 'web1234';

-- More on connect user to database: https://dev.mysql.com/doc/mysql-shell/8.0/en/mysql-shell-connection-using-parameters.html