CREATE TABLE IF NOT EXISTS spells (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    mana_cost BIGINT NOT NULL,
    realm VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wizard_spells (
    wizard_id INTEGER NOT NULL,
    spell_id INTEGER NOT NULL,
    learned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (wizard_id, spell_id)
);

CREATE INDEX idx_spells_realm ON spells(realm);
CREATE INDEX idx_wizard_spells_wizard_id ON wizard_spells(wizard_id);
CREATE INDEX idx_wizard_spells_spell_id ON wizard_spells(spell_id);