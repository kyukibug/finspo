-- Users table
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255),
  email VARCHAR(255) UNIQUE,
  google_id VARCHAR(255) UNIQUE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Categories table
CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  user_id INT,
  name VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Clothing Items table
CREATE TABLE clothing_items (
  id SERIAL PRIMARY KEY,
  user_id INT,
  category_id INT,
  image_url VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Sandbox table
CREATE TABLE sandbox (
  id SERIAL PRIMARY KEY,
  user_id INT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Tags table
CREATE TABLE tags (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Sandbox Positions table
CREATE TABLE sandbox_positions (
  id SERIAL PRIMARY KEY,
  sandbox_id INT,
  clothing_item_id INT,
  position_x FLOAT,
  position_y FLOAT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (sandbox_id) REFERENCES sandbox(id),
  FOREIGN KEY (clothing_item_id) REFERENCES clothing_items(id)
);

-- Clothing Item Tags Junction Table
CREATE TABLE clothing_item_tags (
  clothing_item_id INT,
  tag_id INT,
  PRIMARY KEY (clothing_item_id, tag_id),
  FOREIGN KEY (clothing_item_id) REFERENCES clothing_items(id),
  FOREIGN KEY (tag_id) REFERENCES tags(id)
);

-- Insert users
INSERT INTO users (username, email, google_id, created_at, updated_at) VALUES
('user1', 'user1@example.com', 'google1', '2024-07-01 12:00:00', '2024-07-01 12:00:00'),
('user2', 'user2@example.com', 'google2', '2024-07-01 12:00:00', '2024-07-01 12:00:00');

-- Insert categories
INSERT INTO categories (user_id, name, created_at, updated_at) VALUES
(1, 'Tops', '2024-07-01 12:00:00', '2024-07-01 12:00:00'),
(1, 'Bottoms', '2024-07-01 12:00:00', '2024-07-01 12:00:00'),
(2, 'Accessories', '2024-07-01 12:00:00', '2024-07-01 12:00:00');

-- Insert clothing items
INSERT INTO clothing_items (user_id, category_id, image_url, created_at, updated_at) VALUES
(1, 1, 'http://example.com/clothing1.jpg', '2024-07-01 12:00:00', '2024-07-01 12:00:00'),
(1, 2, 'http://example.com/clothing2.jpg', '2024-07-01 12:00:00', '2024-07-01 12:00:00'),
(2, 3, 'http://example.com/clothing3.jpg', '2024-07-01 12:00:00', '2024-07-01 12:00:00');

-- Insert sandbox
INSERT INTO sandbox (user_id, created_at, updated_at) VALUES
(1, '2024-07-01 12:00:00', '2024-07-01 12:00:00'),
(2, '2024-07-01 12:00:00', '2024-07-01 12:00:00');

-- Insert tags
INSERT INTO tags (name, created_at, updated_at) VALUES
('Summer', '2024-07-01 12:00:00', '2024-07-01 12:00:00'),
('Winter', '2024-07-01 12:00:00', '2024-07-01 12:00:00');

-- Insert sandbox positions
INSERT INTO sandbox_positions (sandbox_id, clothing_item_id, position_x, position_y, created_at, updated_at) VALUES
(1, 1, 50.0, 50.0, '2024-07-01 12:00:00', '2024-07-01 12:00:00'),
(1, 2, 100.0, 150.0, '2024-07-01 12:00:00', '2024-07-01 12:00:00'),
(2, 3, 200.0, 250.0, '2024-07-01 12:00:00', '2024-07-01 12:00:00');

-- Insert clothing item tags
INSERT INTO clothing_item_tags (clothing_item_id, tag_id) VALUES
(1, 1),
(2, 2),
(3, 1);


-- Functions

CREATE OR REPLACE FUNCTION get_clothes_by_user(_user_id INT)
	returns TABLE (
    id INT, 
    user_id INT, 
    category_id INT, 
    image_url VARCHAR(255), 
    created_at TIMESTAMP, 
    updated_at TIMESTAMP, 
    tags text)
	language sql
	security definer
as $$
  SELECT c.id,
    c.user_id,
    c.category_id,
    c.image_url,
    c.created_at,
    c.updated_at,
     array_to_json(array_agg(tags)) as tags_list
  FROM clothing_items c
  LEFT JOIN clothing_item_tags cit
    ON cit.clothing_item_id = c.id
  LEFT JOIN tags
    ON tags.id = cit.tag_id
  WHERE c.user_id = _user_id
  GROUP BY c.id
$$;

CREATE OR REPLACE FUNCTION get_clothes_by_user_and_id(_user_id INT, _id INT)
	returns TABLE (
    id INT, 
    user_id INT, 
    category_id INT, 
    image_url VARCHAR(255), 
    created_at TIMESTAMP, 
    updated_at TIMESTAMP, 
    tags text)
	language sql
	security definer
as $$
  SELECT c.id,
    c.user_id,
    c.category_id,
    c.image_url,
    c.created_at,
    c.updated_at,
     array_to_json(array_agg(tags)) as tags_list
  FROM clothing_items c
  LEFT JOIN clothing_item_tags cit
    ON cit.clothing_item_id = c.id
  LEFT JOIN tags
    ON tags.id = cit.tag_id
  WHERE c.user_id = _user_id
    AND c.id = _id
  GROUP BY c.id
$$;

CREATE OR REPLACE FUNCTION get_tags_by_user(_user_id INT) 
  RETURNS TABLE (
    id INT,
    name VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
  )
	language sql
	security definer
  as $$
    SELECT t.id,
      t.name,
      t.created_at,
      t.updated_at
    FROM tags t
    LEFT JOIN clothing_item_tags cit
      ON cit.tag_id = t.id
    LEFT JOIN clothing_items c
      ON c.id = cit.clothing_item_id
    WHERE c.user_id = _user_id
    GROUP BY t.id
  $$;