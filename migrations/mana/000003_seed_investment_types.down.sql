-- Remove seeded investment types
DELETE FROM investment_types 
WHERE name IN (
    'Novice Spell Bond',
    'Mystic Market Fund',
    'Elemental Ventures',
    'Dragon''s Hoard',
    'Phoenix Rising'
);