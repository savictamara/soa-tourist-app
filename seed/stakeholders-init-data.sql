-- Stakeholders seed (safe to run multiple times).
-- Runs independently from application startup logic.

INSERT INTO users (
    "Username",
    "PasswordHash",
    "Email",
    "Role",
    "IsBlocked",
    "FirstName",
    "LastName",
    "Biography",
    "Motto",
    "ProfileImage"
)
VALUES (
    'admin',
    'AQAAAAIAAYagAAAAEJ8EZ+TKCKLcZwSMjdnXfmUkzwYa6lf4GmjRshTFQn5q+Zm+qDThqnUM6MvghwG32Q==',
    'admin@stakeholders.local',
    'Administrator',
    FALSE,
    'System',
    'Admin',
    'Administrator account seeded for KT docker testing.',
    'Keep services healthy',
    'https://images.unsplash.com/photo-1527980965255-d3b416303d12'
)
ON CONFLICT ("Username") DO NOTHING;

UPDATE users
SET
    "FirstName" = 'System',
    "LastName" = 'Admin',
    "Biography" = 'Administrator account seeded for KT docker testing.',
    "Motto" = 'Keep services healthy',
    "ProfileImage" = 'https://images.unsplash.com/photo-1527980965255-d3b416303d12'
WHERE "Username" = 'admin';

INSERT INTO users (
    "Username",
    "PasswordHash",
    "Email",
    "Role",
    "IsBlocked",
    "FirstName",
    "LastName",
    "Biography",
    "Motto",
    "ProfileImage"
)
VALUES (
    'jovana',
    'AQAAAAIAAYagAAAAEBYKqToG83/YvhvvUIgmSDIKUwjbu2/l7iAtqFci6A5yCe9sAPNBAJNVK5AL8v1Cdw==',
    'jovana@stakeholders.local',
    'Tourist',
    FALSE,
    'Jovana',
    'Ilic',
    'Sample tourist user added through SQL seed.',
    'Ready for KT testing',
    'https://images.unsplash.com/photo-1494790108377-be9c29b29330'
)
ON CONFLICT ("Username") DO NOTHING;

UPDATE users
SET "PasswordHash" = 'AQAAAAIAAYagAAAAEBYKqToG83/YvhvvUIgmSDIKUwjbu2/l7iAtqFci6A5yCe9sAPNBAJNVK5AL8v1Cdw=='
WHERE "Username" = 'jovana';

UPDATE users
SET
    "FirstName" = 'Jovana',
    "LastName" = 'Ilic'
WHERE "Username" = 'jovana';

INSERT INTO users (
    "Username",
    "PasswordHash",
    "Email",
    "Role",
    "IsBlocked",
    "FirstName",
    "LastName",
    "Biography",
    "Motto",
    "ProfileImage"
)
VALUES (
    'marko',
    'AQAAAAIAAYagAAAAEI1a1D0Mij0QVmEgc09+uhJsNhzJ/XIs8n1XYqXEQnecUfud2yQsTj3EEyF3yG4I6A==',
    'marko@stakeholders.local',
    'Guide',
    FALSE,
    'Marko',
    'Petrovic',
    'Seeded tour author account for creating and publishing tours.',
    'Build memorable routes',
    'https://images.unsplash.com/photo-1500648767791-00dcc994a43e'
)
ON CONFLICT ("Username") DO NOTHING;

UPDATE users
SET
    "PasswordHash" = 'AQAAAAIAAYagAAAAEI1a1D0Mij0QVmEgc09+uhJsNhzJ/XIs8n1XYqXEQnecUfud2yQsTj3EEyF3yG4I6A==',
    "Email" = 'marko@stakeholders.local',
    "Role" = 'Guide',
    "IsBlocked" = FALSE,
    "FirstName" = 'Marko',
    "LastName" = 'Petrovic',
    "Biography" = 'Seeded tour author account for creating and publishing tours.',
    "Motto" = 'Build memorable routes',
    "ProfileImage" = 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e'
WHERE "Username" = 'marko';

DELETE FROM users WHERE "Username" = 'hashseed_jovana123';
