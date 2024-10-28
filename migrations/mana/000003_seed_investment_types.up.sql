-- Seed initial investment types
INSERT INTO investment_types 
(name, description, min_amount, max_amount, duration_hours, base_return_rate, risk_level)
VALUES 
    ('Novice Spell Bond', 
     'Low-risk, short-term investment suitable for beginning wizards. Guaranteed small returns with minimal risk.',
     100, 1000, 24, 2.5, 1),
    
    ('Mystic Market Fund', 
     'Balanced investment with moderate returns. Good for wizards looking for steady growth.',
     500, 5000, 72, 5.0, 2),
    
    ('Elemental Ventures', 
     'Higher risk investment with potentially greater returns. Recommended for experienced wizards.',
     1000, 10000, 168, 8.5, 3),
    
    ('Dragon''s Hoard', 
     'High-risk, high-reward long-term investment. For wizards with strong risk tolerance.',
     5000, 50000, 336, 12.0, 4),
    
    ('Phoenix Rising', 
     'Extremely volatile investment with maximum potential returns. Only for the bravest wizards.',
     10000, NULL, 720, 20.0, 5);