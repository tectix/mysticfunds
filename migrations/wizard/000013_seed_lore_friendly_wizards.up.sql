-- Seed lore-friendly wizards for exploration
-- These wizards represent the diverse elemental masters from across the realms
-- Using a special system user_id (999) for exploration wizards

-- Insert lore-friendly wizards with system user_id
INSERT INTO wizards (user_id, name, realm, element, mana_balance, experience_points, level, created_at) VALUES
    
-- Fire Element - Pyrrhian Flame
(999, 'Ignis Pyroclast', 'Pyrrhian Flame', 'Fire', 8250, 12400, 15, NOW() - INTERVAL '45 days'),
(999, 'Salamandrine Emberheart', 'Pyrrhian Flame', 'Fire', 6750, 8900, 12, NOW() - INTERVAL '32 days'),
(999, 'Vulcan Ashforged', 'Pyrrhian Flame', 'Fire', 11200, 18600, 18, NOW() - INTERVAL '67 days'),

-- Water Element - Thalorion Depths  
(999, 'Luna Tidewhisper', 'Thalorion Depths', 'Water', 6800, 8900, 12, NOW() - INTERVAL '38 days'),
(999, 'Coral Deepcurrent', 'Thalorion Depths', 'Water', 9400, 14200, 16, NOW() - INTERVAL '52 days'),
(999, 'Nereid Moonbound', 'Thalorion Depths', 'Water', 12800, 21300, 19, NOW() - INTERVAL '78 days'),

-- Shadow Element - Umbros
(999, 'Vex Shadowbane', 'Umbros', 'Shadow', 15500, 24600, 18, NOW() - INTERVAL '89 days'),
(999, 'Nyx Voidwhisper', 'Umbros', 'Shadow', 13200, 19800, 17, NOW() - INTERVAL '71 days'),
(999, 'Shade Umbralmancer', 'Umbros', 'Shadow', 8900, 11400, 14, NOW() - INTERVAL '41 days'),

-- Air Element - Zepharion Heights
(999, 'Zephyr Stormcaller', 'Zepharion Heights', 'Air', 9200, 15300, 14, NOW() - INTERVAL '48 days'),
(999, 'Gale Cycloneborn', 'Zepharion Heights', 'Air', 10700, 17600, 16, NOW() - INTERVAL '59 days'),
(999, 'Tempest Skyweaver', 'Zepharion Heights', 'Air', 7800, 10200, 13, NOW() - INTERVAL '35 days'),

-- Earth Element - Terravine Hollow
(999, 'Terra Stoneforge', 'Terravine Hollow', 'Earth', 7400, 18700, 16, NOW() - INTERVAL '63 days'),
(999, 'Root Worldtree', 'Terravine Hollow', 'Earth', 14600, 23400, 19, NOW() - INTERVAL '85 days'),
(999, 'Moss Petrified', 'Terravine Hollow', 'Earth', 6200, 7800, 11, NOW() - INTERVAL '28 days'),

-- Light Element - Virelya
(999, 'Solaris Dawnbringer', 'Virelya', 'Light', 11600, 10800, 13, NOW() - INTERVAL '39 days'),
(999, 'Lumina Radiantweaver', 'Virelya', 'Light', 13900, 20100, 17, NOW() - INTERVAL '68 days'),
(999, 'Prism Truthseeker', 'Virelya', 'Light', 8500, 12700, 14, NOW() - INTERVAL '44 days'),

-- Time Element - Chronarxis
(999, 'Chronos Timekeeper', 'Chronarxis', 'Time', 18900, 32100, 20, NOW() - INTERVAL '112 days'),
(999, 'Aeon Spiralwalker', 'Chronarxis', 'Time', 16200, 27800, 19, NOW() - INTERVAL '95 days'),
(999, 'Temporal Fractured', 'Chronarxis', 'Time', 12300, 19500, 17, NOW() - INTERVAL '72 days'),

-- Void Element - Nyxthar
(999, 'Null the Voidwalker', 'Nyxthar', 'Void', 13200, 21800, 17, NOW() - INTERVAL '76 days'),
(999, 'Entropy Silencepriest', 'Nyxthar', 'Void', 15800, 26400, 18, NOW() - INTERVAL '91 days'),
(999, 'Nihil Forgottenone', 'Nyxthar', 'Void', 9700, 14900, 15, NOW() - INTERVAL '51 days'),

-- Spirit Element - Aetherion  
(999, 'Ethereal Soulweaver', 'Aetherion', 'Spirit', 9800, 7200, 11, NOW() - INTERVAL '26 days'),
(999, 'Phantom Dreamdancer', 'Aetherion', 'Spirit', 11400, 16800, 15, NOW() - INTERVAL '58 days'),
(999, 'Astral Thoughtbinder', 'Aetherion', 'Spirit', 14100, 22600, 18, NOW() - INTERVAL '81 days'),

-- Metal Element - Technarok
(999, 'Ferros Gearwright', 'Technarok', 'Metal', 16700, 28400, 19, NOW() - INTERVAL '97 days'),
(999, 'Steel Logicengine', 'Technarok', 'Metal', 10900, 16100, 15, NOW() - INTERVAL '55 days'),
(999, 'Brass Recursive', 'Technarok', 'Metal', 8600, 13200, 14, NOW() - INTERVAL '46 days');

-- Ensure guilds table has description column (defensive programming)
ALTER TABLE guilds ADD COLUMN IF NOT EXISTS description TEXT;

-- Create some sample guilds and assign wizards
INSERT INTO guilds (name, description, created_at) VALUES 
('Salamandrine Lords', 'Elite fire mages of Pyrrhian Flame', NOW()),
('Void Touched', 'Shadow masters who dance with entropy', NOW()),
('Temporal Seekers', 'Time manipulators seeking ultimate understanding', NOW()),
('Depth Wardens', 'Guardians of the deep ocean realms', NOW()),
('Sky Dancers', 'Masters of wind and storm from floating isles', NOW()),
('Stone Wardens', 'Ancient protectors of the buried realm', NOW()),
('Light Bearers', 'Pure souls who channel radiant truth', NOW()),
('Soul Weavers', 'Spirit guides between realms', NOW()),
('Logic Engineers', 'Techno-mages of mechanical precision', NOW())
ON CONFLICT (name) DO NOTHING;

-- Assign some wizards to guilds using a nested query
UPDATE wizards SET guild_id = (SELECT id FROM guilds WHERE name = 'Salamandrine Lords')
WHERE name IN ('Ignis Pyroclast', 'Vulcan Ashforged') AND user_id = 999;

UPDATE wizards SET guild_id = (SELECT id FROM guilds WHERE name = 'Void Touched')
WHERE name IN ('Vex Shadowbane', 'Nyx Voidwhisper') AND user_id = 999;

UPDATE wizards SET guild_id = (SELECT id FROM guilds WHERE name = 'Temporal Seekers')
WHERE name IN ('Chronos Timekeeper', 'Aeon Spiralwalker') AND user_id = 999;

UPDATE wizards SET guild_id = (SELECT id FROM guilds WHERE name = 'Depth Wardens')
WHERE name IN ('Luna Tidewhisper', 'Nereid Moonbound') AND user_id = 999;

UPDATE wizards SET guild_id = (SELECT id FROM guilds WHERE name = 'Sky Dancers')
WHERE name IN ('Zephyr Stormcaller', 'Gale Cycloneborn') AND user_id = 999;

UPDATE wizards SET guild_id = (SELECT id FROM guilds WHERE name = 'Stone Wardens')
WHERE name IN ('Terra Stoneforge', 'Root Worldtree') AND user_id = 999;

UPDATE wizards SET guild_id = (SELECT id FROM guilds WHERE name = 'Light Bearers')
WHERE name IN ('Solaris Dawnbringer', 'Lumina Radiantweaver') AND user_id = 999;

UPDATE wizards SET guild_id = (SELECT id FROM guilds WHERE name = 'Soul Weavers')
WHERE name IN ('Ethereal Soulweaver', 'Phantom Dreamdancer') AND user_id = 999;

UPDATE wizards SET guild_id = (SELECT id FROM guilds WHERE name = 'Logic Engineers')
WHERE name IN ('Ferros Gearwright', 'Steel Logicengine') AND user_id = 999;