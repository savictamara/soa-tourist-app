-- Blog seed (safe to run multiple times).

INSERT INTO blog_posts ("Title", "DescriptionMarkdown", "CreatedAtUtc", "AuthorUsername")
SELECT
    'Welcome to Tourist Blog',
    'Seed post created for KT docker verification.',
    NOW(),
    'jovana'
WHERE NOT EXISTS (
    SELECT 1 FROM blog_posts WHERE "Title" = 'Welcome to Tourist Blog'
);

INSERT INTO blog_posts ("Title", "DescriptionMarkdown", "CreatedAtUtc", "AuthorUsername")
SELECT
    'Belgrade Walking Tour',
    'Second seed post for comments and likes checks.',
    NOW(),
    'jovana'
WHERE NOT EXISTS (
    SELECT 1 FROM blog_posts WHERE "Title" = 'Belgrade Walking Tour'
);

INSERT INTO blog_post_images ("BlogPostId", "ImageUrl")
SELECT p."Id", 'https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1'
FROM blog_posts p
WHERE p."Title" = 'Welcome to Tourist Blog'
  AND NOT EXISTS (
      SELECT 1
      FROM blog_post_images i
      WHERE i."BlogPostId" = p."Id"
        AND i."ImageUrl" = 'https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1'
  );

INSERT INTO blog_comments ("BlogPostId", "AuthorUsername", "AuthorRole", "Text", "CreatedAtUtc", "LastModifiedAtUtc")
SELECT
    p."Id",
    'jovana',
    'Tourist',
    'Great post. Seed comment is visible.',
    NOW(),
    NOW()
FROM blog_posts p
WHERE p."Title" = 'Welcome to Tourist Blog'
  AND NOT EXISTS (
      SELECT 1
      FROM blog_comments c
      WHERE c."BlogPostId" = p."Id"
        AND c."AuthorUsername" = 'jovana'
        AND c."Text" = 'Great post. Seed comment is visible.'
  );

INSERT INTO blog_likes ("BlogPostId", "Username", "CreatedAtUtc")
SELECT p."Id", 'admin', NOW()
FROM blog_posts p
WHERE p."Title" = 'Belgrade Walking Tour'
  AND NOT EXISTS (
      SELECT 1
      FROM blog_likes l
      WHERE l."BlogPostId" = p."Id"
        AND l."Username" = 'admin'
  );

-- Ensure admin has no seeded posts or comments from previous runs.
UPDATE blog_posts
SET "AuthorUsername" = 'jovana'
WHERE "AuthorUsername" = 'admin';

DELETE FROM blog_comments
WHERE "AuthorUsername" = 'admin';
