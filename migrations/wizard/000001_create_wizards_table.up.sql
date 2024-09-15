CREATE TABLE IF NOT EXISTS wizards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    realm VARCHAR(50) NOT NULL,
    element VARCHAR(50),
    guild_id INTEGER,
    mana_balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wizards_user_id ON wizards(user_id);
CREATE INDEX idx_wizards_realm ON wizards(realm);
CREATE INDEX idx_wizards_element ON wizards(element);
CREATE INDEX idx_wizards_guild_id ON wizards(guild_id);