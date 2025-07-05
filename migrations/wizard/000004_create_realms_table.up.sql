-- Create realms table with lore and artifacts
CREATE TABLE realms (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    element VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    lore TEXT NOT NULL,
    artifact_name VARCHAR(100) NOT NULL,
    artifact_description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert the mystical realms
INSERT INTO realms (name, element, description, lore, artifact_name, artifact_description) VALUES 
(
    'Pyrrhian Flame',
    'Fire / Heat / Chaos',
    'Realm of eternal fire and volcanic titans',
    'Home to the Salamandrine Lords and the Eternal Forge. Pyrrhian Flame births volcanic titans and flame-bonded warriors. Time moves faster here, aging all who enter.',
    'Heart of Cinder',
    'A molten gem that grants destructive power and burns away lies.'
),
(
    'Zepharion Heights', 
    'Wind / Sky / Sound',
    'Floating islands around an eternal cyclone',
    'Floating islands encircle a permanent cyclone known as The Whisper. Skyborn sages ride wind-serpents and wield songs that bend reality.',
    'Aeon Harp',
    'Plays melodies that control storms and memories.'
),
(
    'Terravine Hollow',
    'Stone / Growth / Gravity', 
    'Ancient buried realm of stone titans',
    'An ancient, buried realm where roots grow like veins and sentient stone titans slumber. Once a great civilization, now petrified into time.',
    'Verdant Core',
    'Grants dominion over life, soil, and rebirth.'
),
(
    'Thalorion Depths',
    'Water / Ice / Depth',
    'Submerged empire of the Moonbound Court',
    'A submerged empire ruled by the Moonbound Court. Time slows here, and the ocean whispers ancient truths. Home to leviathans and drowned prophets.',
    'Tideglass Mirror',
    'Sees through illusions and to possible futures.'
),
(
    'Virelya',
    'Light / Purity / Illumination',
    'Blinding paradise of pure truth',
    'A blinding paradise where truth manifests as form. Ruled by beings known as Radiants. Mortals must wear veilshades to even look upon it.',
    'Lumen Shard',
    'Reveals the true name of anything it touches.'
),
(
    'Umbros',
    'Shadow / Secrets / Corruption',
    'Void-split realm where light cannot reach',
    'Light cannot reach this void-split realm. Every whisper is a thought stolen, every step a forgotten path. Shadowmages barter in memories.',
    'Eclipse Fang',
    'Severs light, binding a soul to darkness.'
),
(
    'Nyxthar',
    'Null / Anti-Matter / Entropy',
    'Realm where reality collapses inward',
    'A realm where reality collapses inward. Voidwalkers and Silence Priests seek ultimate release from being. To enter is to forget existence.',
    'Hollow Crown',
    'Nullifies all magic and erases history.'
),
(
    'Aetherion',
    'Spirit / Soul / Dream',
    'Realm between realms of dreaming dead',
    'The realm between realms, where the dreaming dead speak. Time is nonlinear, and the laws of logic bend to desire. Spirits travel as thought.',
    'Soulforge Locket',
    'Binds spirits to bodies or frees them eternally.'
),
(
    'Chronarxis',
    'Time / Fate / Chronomancy',
    'Spiral palace of fractured timelines',
    'A spiral palace where timelines fracture and reform. Timekeepers judge anomalies and anomalies fight back. Accessed only through ancient rituals.',
    'Clockheart Mechanism',
    'Rewinds one moment once, but at a cost.'
),
(
    'Technarok',
    'Metal / Machines / Order',
    'Fusion of steel gods and nano-intelligences',
    'A fusion of ancient steel gods and nano-intelligences. Run on logic and decay. Home to sentient forges and recursive codebeasts.',
    'Iron Synapse',
    'Merges user with machine intelligence.'
);