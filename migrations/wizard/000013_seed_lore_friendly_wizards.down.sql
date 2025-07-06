-- Remove seeded lore-friendly wizards and related data

-- Remove wizards created by system user_id 999
DELETE FROM wizards WHERE user_id = 999;

-- Remove sample guilds (be careful not to remove user-created guilds)
DELETE FROM guilds WHERE name IN (
    'Salamandrine Lords',
    'Void Touched', 
    'Temporal Seekers',
    'Depth Wardens',
    'Sky Dancers',
    'Stone Wardens',
    'Light Bearers',
    'Soul Weavers',
    'Logic Engineers'
);