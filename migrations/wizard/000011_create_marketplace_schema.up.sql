-- Marketplace System: Artifacts, Scrolls, and Spells
-- This migration creates the foundation for the magical marketplace

-- Artifacts table: Legendary items tied to realms
CREATE TABLE IF NOT EXISTS artifacts (
    id SERIAL PRIMARY KEY,
    realm_id INTEGER NOT NULL REFERENCES realms(id),
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    lore TEXT NOT NULL,
    power_level INTEGER NOT NULL CHECK (power_level BETWEEN 1 AND 10),
    rarity VARCHAR(20) NOT NULL CHECK (rarity IN ('Common', 'Uncommon', 'Rare', 'Epic', 'Legendary', 'Mythical', 'Forbidden')),
    mana_cost BIGINT NOT NULL,
    artifact_type VARCHAR(50) NOT NULL, -- 'Weapon', 'Armor', 'Accessory', 'Tome', 'Relic'
    special_abilities TEXT[], -- Array of special abilities
    requirements TEXT, -- Level or other requirements
    image_url VARCHAR(255),
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Scrolls table: Knowledge that can be purchased with mana
CREATE TABLE IF NOT EXISTS scrolls (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    skill_type VARCHAR(50) NOT NULL, -- 'Combat', 'Magic', 'Crafting', 'Alchemy', 'Enchanting'
    skill_level INTEGER NOT NULL CHECK (skill_level BETWEEN 1 AND 5),
    mana_cost BIGINT NOT NULL,
    prerequisites TEXT[], -- Array of required skills/scrolls
    benefits TEXT NOT NULL, -- What the scroll teaches
    rarity VARCHAR(20) NOT NULL CHECK (rarity IN ('Common', 'Uncommon', 'Rare', 'Epic', 'Legendary')),
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Spells table: Magic that must be learned from other wizards
CREATE TABLE IF NOT EXISTS spells (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    spell_school VARCHAR(50) NOT NULL, -- 'Elemental', 'Illusion', 'Divination', 'Transmutation', 'Conjuration'
    element VARCHAR(50), -- Optional: 'Fire', 'Water', 'Earth', 'Air', 'Light', 'Shadow'
    power_level INTEGER NOT NULL CHECK (power_level BETWEEN 1 AND 10),
    mana_cost_to_learn BIGINT NOT NULL,
    mana_cost_to_cast INTEGER NOT NULL,
    requirements TEXT, -- Level, realm, or other requirements
    effects TEXT NOT NULL, -- What the spell does
    rarity VARCHAR(20) NOT NULL CHECK (rarity IN ('Common', 'Uncommon', 'Rare', 'Epic', 'Legendary', 'Forbidden')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Wizard inventory: Artifacts owned by wizards
CREATE TABLE IF NOT EXISTS wizard_artifacts (
    id SERIAL PRIMARY KEY,
    wizard_id INTEGER NOT NULL REFERENCES wizards(id) ON DELETE CASCADE,
    artifact_id INTEGER NOT NULL REFERENCES artifacts(id),
    acquired_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_equipped BOOLEAN NOT NULL DEFAULT false,
    UNIQUE(wizard_id, artifact_id) -- Wizard can only own one of each artifact
);

-- Wizard scrolls: Scrolls learned by wizards
CREATE TABLE IF NOT EXISTS wizard_scrolls (
    id SERIAL PRIMARY KEY,
    wizard_id INTEGER NOT NULL REFERENCES wizards(id) ON DELETE CASCADE,
    scroll_id INTEGER NOT NULL REFERENCES scrolls(id),
    learned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    mastery_level INTEGER NOT NULL DEFAULT 1 CHECK (mastery_level BETWEEN 1 AND 5),
    UNIQUE(wizard_id, scroll_id) -- Wizard can only learn each scroll once
);

-- Wizard spells: Spells known by wizards
CREATE TABLE IF NOT EXISTS wizard_spells (
    id SERIAL PRIMARY KEY,
    wizard_id INTEGER NOT NULL REFERENCES wizards(id) ON DELETE CASCADE,
    spell_id INTEGER NOT NULL REFERENCES spells(id),
    learned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    learned_from_wizard_id INTEGER REFERENCES wizards(id), -- Who taught this spell
    mastery_level INTEGER NOT NULL DEFAULT 1 CHECK (mastery_level BETWEEN 1 AND 10),
    times_cast INTEGER NOT NULL DEFAULT 0,
    UNIQUE(wizard_id, spell_id) -- Wizard can only learn each spell once
);

-- Spell teaching: Track which wizards can teach which spells
CREATE TABLE IF NOT EXISTS wizard_spell_teaching (
    id SERIAL PRIMARY KEY,
    wizard_id INTEGER NOT NULL REFERENCES wizards(id) ON DELETE CASCADE,
    spell_id INTEGER NOT NULL REFERENCES spells(id),
    can_teach BOOLEAN NOT NULL DEFAULT true,
    teaching_price BIGINT NOT NULL, -- Mana cost for others to learn
    max_students INTEGER DEFAULT NULL, -- NULL means unlimited
    students_taught INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(wizard_id, spell_id)
);

-- Marketplace transactions: Track all purchases and spell learning
CREATE TABLE IF NOT EXISTS marketplace_transactions (
    id SERIAL PRIMARY KEY,
    buyer_wizard_id INTEGER NOT NULL REFERENCES wizards(id),
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('artifact', 'scroll', 'spell_learning')),
    item_id INTEGER NOT NULL, -- artifact_id, scroll_id, or spell_id
    mana_spent BIGINT NOT NULL,
    seller_wizard_id INTEGER REFERENCES wizards(id), -- NULL for system sales, wizard ID for spell learning
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    notes TEXT -- Additional transaction details
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_artifacts_realm_id ON artifacts(realm_id);
CREATE INDEX IF NOT EXISTS idx_artifacts_rarity ON artifacts(rarity);
CREATE INDEX IF NOT EXISTS idx_artifacts_available ON artifacts(is_available);
CREATE INDEX IF NOT EXISTS idx_scrolls_skill_type ON scrolls(skill_type);
CREATE INDEX IF NOT EXISTS idx_scrolls_available ON scrolls(is_available);
CREATE INDEX IF NOT EXISTS idx_spells_school ON spells(spell_school);
CREATE INDEX IF NOT EXISTS idx_spells_element ON spells(element);
CREATE INDEX IF NOT EXISTS idx_wizard_artifacts_wizard_id ON wizard_artifacts(wizard_id);
CREATE INDEX IF NOT EXISTS idx_wizard_scrolls_wizard_id ON wizard_scrolls(wizard_id);
CREATE INDEX IF NOT EXISTS idx_wizard_spells_wizard_id ON wizard_spells(wizard_id);
CREATE INDEX IF NOT EXISTS idx_wizard_spells_learned_from ON wizard_spells(learned_from_wizard_id);
CREATE INDEX IF NOT EXISTS idx_wizard_spell_teaching_wizard_id ON wizard_spell_teaching(wizard_id);
CREATE INDEX IF NOT EXISTS idx_wizard_spell_teaching_spell_id ON wizard_spell_teaching(spell_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_transactions_buyer ON marketplace_transactions(buyer_wizard_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_transactions_seller ON marketplace_transactions(seller_wizard_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_transactions_type ON marketplace_transactions(transaction_type);