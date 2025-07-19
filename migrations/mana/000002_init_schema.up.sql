
-- Create indices for better query performance
CREATE INDEX IF NOT EXISTS idx_mana_transactions_from_wizard ON mana_transactions(from_wizard_id);
CREATE INDEX IF NOT EXISTS idx_mana_transactions_to_wizard ON mana_transactions(to_wizard_id);
CREATE INDEX IF NOT EXISTS idx_mana_transactions_created_at ON mana_transactions(created_at);

-- Create investment_types table
CREATE TABLE IF NOT EXISTS investment_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    min_amount BIGINT NOT NULL,
    max_amount BIGINT,
    duration_hours INTEGER NOT NULL,
    base_return_rate DECIMAL(5,2) NOT NULL,
    risk_level INTEGER NOT NULL CHECK (risk_level BETWEEN 1 AND 5),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create wizard_investments table
CREATE TABLE IF NOT EXISTS wizard_investments (
    id SERIAL PRIMARY KEY,
    wizard_id INTEGER NOT NULL,
    investment_type_id INTEGER NOT NULL REFERENCES investment_types(id),
    amount BIGINT NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'completed', 'failed')),
    actual_return_rate DECIMAL(5,2),
    returned_amount BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indices for investment tables
CREATE INDEX IF NOT EXISTS idx_investment_types_risk_level ON investment_types(risk_level);
CREATE INDEX IF NOT EXISTS idx_wizard_investments_wizard_id ON wizard_investments(wizard_id);
CREATE INDEX IF NOT EXISTS idx_wizard_investments_status ON wizard_investments(status);
CREATE INDEX IF NOT EXISTS idx_wizard_investments_end_time ON wizard_investments(end_time);

-- Create update trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
CREATE TRIGGER update_investment_types_updated_at
    BEFORE UPDATE ON investment_types
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_wizard_investments_updated_at
    BEFORE UPDATE ON wizard_investments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();