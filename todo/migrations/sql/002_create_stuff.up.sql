INSERT INTO boards
(id, user_id, title)
VALUES
('bbbbbbbb-0000-aaaa-0000-dddddddddddd', '00000000-0000-eeee-0000-000000000000', 'Initial Board');

INSERT INTO columns
(id, board_id, user_id, title, position)
VALUES
('cccccccc-0000-0000-0000-000000000000', 'bbbbbbbb-0000-aaaa-0000-dddddddddddd', '00000000-0000-eeee-0000-000000000000', 'Initial Column', 0);

INSERT INTO cards
(id, column_id, user_id, title, description, position)
VALUES
('cccccccc-aaaa-0000-dddd-dddddddddddd', 'cccccccc-0000-0000-0000-000000000000', '00000000-0000-eeee-0000-000000000000', 'Initial Card', 'Initial Description', 0);
