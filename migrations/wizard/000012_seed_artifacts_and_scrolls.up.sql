-- Seed data for artifacts, scrolls, and spells
-- This migration populates the marketplace with initial magical items

-- First, ensure the artifacts table constraint allows 'Forbidden' rarity
ALTER TABLE artifacts DROP CONSTRAINT IF EXISTS artifacts_rarity_check;
ALTER TABLE artifacts ADD CONSTRAINT artifacts_rarity_check 
CHECK (rarity IN ('Common', 'Uncommon', 'Rare', 'Epic', 'Legendary', 'Mythical', 'Forbidden'));

-- Insert Artifacts for each realm
-- Pyrrhian Flame (realm_id = 1)
INSERT INTO artifacts (realm_id, name, description, lore, power_level, rarity, mana_cost, artifact_type, special_abilities, requirements, image_url) VALUES
(1, 'Emberforge Gauntlets', 'Gauntlets that channel the pure essence of volcanic fire', 'Forged in the heart of Mount Pyrrhia during the Great Eruption. These gauntlets were worn by Ignis the Flamebearer, first ruler of the Pyrrhian Flame realm. Legend says they can melt any metal and grant immunity to fire.', 8, 'Legendary', 15000, 'Armor', '{"Fire Immunity", "Metal Melting", "+50% Fire Spell Power"}', 'Level 5+, Fire Affinity', '/images/artifacts/emberforge_gauntlets.png'),
(1, 'Volcanic Crown', 'A crown carved from obsidian and embedded with fire gems', 'The crown of the Pyrrhian rulers, said to contain the wisdom of ancient fire dragons. When worn, it grants the ability to see through flame and smoke.', 9, 'Mythical', 25000, 'Accessory', '{"True Sight through Fire", "Dragon Communication", "+100% Fire Resistance"}', 'Level 8+, Realm Leader', '/images/artifacts/volcanic_crown.png'),
(1, 'Salamander Blade', 'A sword that burns with eternal flame', 'This blade was tempered in salamander breath and has never been extinguished. It grows stronger with each enemy defeated in combat.', 6, 'Epic', 8000, 'Weapon', '{"Eternal Flame", "Growing Power", "Cauterizing Strikes"}', 'Level 3+, Combat Experience', '/images/artifacts/salamander_blade.png'),

-- Zepharion Heights (realm_id = 2)  
(2, 'Stormcaller Staff', 'A staff that commands the winds and lightning', 'Carved from a tree struck by lightning a thousand times, this staff was wielded by the Storm King of Zepharion. It can summon cyclones and control weather patterns.', 9, 'Mythical', 20000, 'Weapon', '{"Weather Control", "Lightning Mastery", "Flight Enhancement"}', 'Level 7+, Air Affinity', '/images/artifacts/stormcaller_staff.png'),
(2, 'Windwalker Boots', 'Boots that allow walking on air currents', 'These boots were crafted by the Sky Smiths using crystallized wind. They grant the wearer the ability to walk on air and move at incredible speeds.', 7, 'Legendary', 12000, 'Armor', '{"Air Walking", "Enhanced Speed", "Perfect Balance"}', 'Level 4+, Dexterity 15+', '/images/artifacts/windwalker_boots.png'),
(2, 'Cyclone Compass', 'A compass that always points to the center of storms', 'This mystical compass was created to navigate the ever-changing winds of Zepharion. It can predict weather and locate magical disturbances.', 5, 'Rare', 5000, 'Accessory', '{"Weather Prediction", "Magic Detection", "Storm Navigation"}', 'Level 2+', '/images/artifacts/cyclone_compass.png'),

-- Terravine Hollow (realm_id = 3)
(3, 'Earthheart Shield', 'A shield made from the core of a stone titan', 'This shield was carved from the heart of an ancient stone titan. It pulses with earth magic and can create protective barriers of stone and root.', 8, 'Legendary', 14000, 'Armor', '{"Stone Barrier", "Root Entanglement", "Earthquake Resistance"}', 'Level 5+, Earth Affinity', '/images/artifacts/earthheart_shield.png'),
(3, 'Vinewarden Circlet', 'A circlet that commands plant life', 'Grown rather than made, this living circlet grants dominion over all plant life. It was worn by the first Druids of Terravine Hollow.', 7, 'Epic', 10000, 'Accessory', '{"Plant Control", "Accelerated Growth", "Nature Communication"}', 'Level 4+, Nature Affinity', '/images/artifacts/vinewarden_circlet.png'),
(3, 'Granite Warhammer', 'A massive hammer that never dulls', 'Forged from a single piece of enchanted granite, this warhammer grows heavier with each successful strike, making it devastating in prolonged combat.', 6, 'Epic', 9000, 'Weapon', '{"Increasing Weight", "Stone Shattering", "Seismic Strike"}', 'Level 3+, Strength 18+', '/images/artifacts/granite_warhammer.png'),

-- Thalorion Depths (realm_id = 4)
(4, 'Tideglass Mirror', 'A mirror that shows possible futures', 'This ancient mirror was crafted by the Moonbound Court from crystallized moonlight and deep-sea glass. It reveals glimpses of potential futures and can see through illusions.', 10, 'Mythical', 30000, 'Relic', '{"Future Sight", "Illusion Piercing", "Moonlight Channeling"}', 'Level 10+, Divination Mastery', '/images/artifacts/tideglass_mirror.png'),
(4, 'Leviathan Scale Armor', 'Armor made from the scales of an ancient sea beast', 'This armor was crafted from the scales of the Great Leviathan that once ruled the deepest trenches. It provides unmatched protection underwater.', 8, 'Legendary', 16000, 'Armor', '{"Water Breathing", "Pressure Immunity", "Aquatic Speed"}', 'Level 6+, Water Affinity', '/images/artifacts/leviathan_scale_armor.png'),
(4, 'Coral Trident', 'A trident that commands the tides', 'Grown in the deepest coral gardens over centuries, this trident can control water currents and summon sea creatures to aid its wielder.', 7, 'Epic', 11000, 'Weapon', '{"Tide Control", "Sea Creature Summoning", "Water Spear"}', 'Level 4+, Swimming Proficiency', '/images/artifacts/coral_trident.png'),

-- Virelya (realm_id = 5)
(5, 'Lumen Shard', 'A crystal that reveals the true name of anything', 'This shard of pure light was formed at the moment of Virelya''s creation. It can pierce any deception and reveal the true nature of all things.', 9, 'Mythical', 22000, 'Relic', '{"True Name Revelation", "Deception Detection", "Pure Light"}', 'Level 8+, Pure Heart', '/images/artifacts/lumen_shard.png'),
(5, 'Radiant Robes', 'Robes woven from concentrated sunlight', 'These robes were woven by the Radiants themselves from threads of pure light. They protect against all forms of darkness and corruption.', 7, 'Legendary', 13000, 'Armor', '{"Darkness Immunity", "Corruption Resistance", "Blinding Flash"}', 'Level 5+, Light Affinity', '/images/artifacts/radiant_robes.png'),
(5, 'Prism Wand', 'A wand that splits light into rainbow magic', 'This crystalline wand refracts light into seven distinct magical frequencies, allowing the wielder to cast multiple types of spells simultaneously.', 6, 'Epic', 8500, 'Weapon', '{"Multi-cast", "Light Manipulation", "Rainbow Bridge"}', 'Level 3+, Spell Variety', '/images/artifacts/prism_wand.png'),

-- Umbros (realm_id = 6)
(6, 'Eclipse Fang', 'A dagger that severs light and binds souls to darkness', 'Forged in the void between light and shadow, this dagger can cut through any illumination and trap souls in eternal darkness.', 9, 'Forbidden', 35000, 'Weapon', '{"Light Severing", "Soul Binding", "Shadow Step"}', 'Level 9+, Dark Pact', '/images/artifacts/eclipse_fang.png'),
(6, 'Voidwalker Cloak', 'A cloak that grants passage through shadows', 'Woven from the essence of the void itself, this cloak allows its wearer to step through shadows and become one with darkness.', 8, 'Legendary', 15000, 'Armor', '{"Shadow Travel", "Invisibility", "Void Immunity"}', 'Level 6+, Shadow Affinity', '/images/artifacts/voidwalker_cloak.png'),
(6, 'Memory Thief Pendant', 'A pendant that steals and stores memories', 'This dark pendant can extract memories from victims and store them for later use. It was created by the Shadow Mages of Umbros.', 7, 'Epic', 12000, 'Accessory', '{"Memory Extraction", "Knowledge Storage", "Mind Reading"}', 'Level 4+, Mental Resistance', '/images/artifacts/memory_thief_pendant.png'),

-- Nyxthar (realm_id = 7)
(7, 'Hollow Crown', 'A crown that nullifies all magic and erases history', 'The ultimate artifact of negation, this crown can unmake reality itself. It was worn by the Null Kings before they erased themselves from existence.', 10, 'Forbidden', 50000, 'Relic', '{"Magic Nullification", "Reality Erasure", "History Deletion"}', 'Level 10+, Void Master', '/images/artifacts/hollow_crown.png'),
(7, 'Entropy Gauntlet', 'A gauntlet that accelerates decay and destruction', 'This gauntlet harnesses the power of entropy itself, causing whatever it touches to age and decay at an accelerated rate.', 8, 'Legendary', 18000, 'Armor', '{"Accelerated Decay", "Time Manipulation", "Destruction Enhancement"}', 'Level 7+, Entropy Affinity', '/images/artifacts/entropy_gauntlet.png'),
(7, 'Null Blade', 'A sword that exists and doesn''t exist simultaneously', 'This paradoxical weapon phases in and out of reality, making it impossible to defend against or predict its strikes.', 7, 'Epic', 14000, 'Weapon', '{"Phase Strikes", "Reality Shift", "Paradox Creation"}', 'Level 5+, Logic Resistance', '/images/artifacts/null_blade.png'),

-- Aetherion (realm_id = 8)
(8, 'Soulforge Locket', 'A locket that binds spirits to bodies or frees them eternally', 'This ancient locket controls the boundary between life and death, spirit and flesh. It was created by the first Spirit Walkers.', 9, 'Mythical', 24000, 'Relic', '{"Soul Binding", "Spirit Liberation", "Life/Death Balance"}', 'Level 8+, Spirit Affinity', '/images/artifacts/soulforge_locket.png'),
(8, 'Dreamweaver Robes', 'Robes that allow travel through dreams and thoughts', 'These ethereal robes shift between reality and dream, allowing the wearer to enter the dreams of others and walk through thoughts.', 7, 'Legendary', 13500, 'Armor', '{"Dream Walking", "Thought Travel", "Astral Projection"}', 'Level 5+, Psychic Resistance', '/images/artifacts/dreamweaver_robes.png'),
(8, 'Ghost Blade', 'A sword that can strike incorporeal beings', 'This spectral weapon exists partially in the spirit realm, allowing it to harm ghosts, spirits, and other incorporeal creatures.', 6, 'Epic', 9500, 'Weapon', '{"Spirit Strike", "Incorporeal Damage", "Ectoplasmic Edge"}', 'Level 3+, Spirit Sight', '/images/artifacts/ghost_blade.png'),

-- Chronarxis (realm_id = 9)
(9, 'Clockheart Mechanism', 'A device that rewinds one moment once, but at great cost', 'This intricate clockwork device can undo a single moment in time, but each use ages the wielder significantly.', 10, 'Forbidden', 40000, 'Relic', '{"Time Rewind", "Moment Manipulation", "Temporal Anchor"}', 'Level 10+, Time Master', '/images/artifacts/clockheart_mechanism.png'),
(9, 'Timekeeper Robes', 'Robes that slow time around the wearer', 'These robes are woven from crystallized time itself, allowing the wearer to move faster relative to everything else.', 8, 'Legendary', 17000, 'Armor', '{"Time Dilation", "Accelerated Reactions", "Temporal Shield"}', 'Level 6+, Time Affinity', '/images/artifacts/timekeeper_robes.png'),
(9, 'Paradox Staff', 'A staff that creates temporal anomalies', 'This staff can create small tears in time, causing temporal paradoxes and anomalies that confuse and disorient enemies.', 7, 'Epic', 12500, 'Weapon', '{"Temporal Anomalies", "Paradox Creation", "Time Loop"}', 'Level 4+, Temporal Stability', '/images/artifacts/paradox_staff.png'),

-- Technarok (realm_id = 10)
(10, 'Iron Synapse', 'A device that merges user with machine intelligence', 'This cybernetic implant connects the user''s mind to the great machine network of Technarok, granting access to vast computational power.', 9, 'Mythical', 26000, 'Relic', '{"Machine Interface", "Computational Enhancement", "Network Access"}', 'Level 8+, Tech Affinity', '/images/artifacts/iron_synapse.png'),
(10, 'Nanoweave Armor', 'Armor that adapts and repairs itself', 'This armor is composed of billions of nanomachines that can reshape themselves to counter any threat and repair any damage.', 8, 'Legendary', 16500, 'Armor', '{"Adaptive Defense", "Self-Repair", "Threat Analysis"}', 'Level 6+, Technology Integration', '/images/artifacts/nanoweave_armor.png'),
(10, 'Data Sword', 'A blade made of crystallized information', 'This weapon exists as pure information given physical form. It can cut through any defense by rewriting the data that defines it.', 7, 'Epic', 11500, 'Weapon', '{"Data Manipulation", "Reality Rewriting", "Information Strike"}', 'Level 4+, Code Mastery', '/images/artifacts/data_sword.png');

-- Insert Scrolls for various skills
INSERT INTO scrolls (name, description, skill_type, skill_level, mana_cost, prerequisites, benefits, rarity) VALUES
('Scroll of Basic Combat', 'Teaches fundamental combat techniques and weapon handling', 'Combat', 1, 500, '{}', 'Learn basic weapon proficiency and combat stances', 'Common'),
('Scroll of Advanced Swordplay', 'Advanced techniques for sword fighting and blade mastery', 'Combat', 3, 2000, '{"Basic Combat"}', 'Master sword techniques and unlock special attacks', 'Uncommon'),
('Scroll of Elemental Basics', 'Introduction to elemental magic and spell casting', 'Magic', 1, 750, '{}', 'Learn to cast basic elemental spells', 'Common'),
('Scroll of Fire Mastery', 'Advanced fire magic techniques and spells', 'Magic', 4, 4000, '{"Elemental Basics"}', 'Master fire magic and unlock powerful fire spells', 'Rare'),
('Scroll of Alchemy Fundamentals', 'Basic potion brewing and ingredient identification', 'Alchemy', 1, 600, '{}', 'Learn to brew basic potions and identify magical ingredients', 'Common'),
('Scroll of Master Alchemist', 'Advanced alchemy including transmutation and rare potions', 'Alchemy', 5, 8000, '{"Alchemy Fundamentals", "Elemental Basics"}', 'Create legendary potions and perform transmutation', 'Legendary'),
('Scroll of Enchanting Basics', 'Introduction to item enchantment and magical enhancement', 'Enchanting', 2, 1200, '{"Elemental Basics"}', 'Learn to enchant weapons and armor with basic effects', 'Uncommon'),
('Scroll of Arcane Engineering', 'Crafting magical constructs and enchanted items', 'Crafting', 4, 5000, '{"Enchanting Basics"}', 'Create magical constructs and advanced enchanted items', 'Epic'),
('Scroll of Meditation Mastery', 'Advanced meditation techniques for mana regeneration', 'Magic', 2, 1000, '{}', 'Increase mana regeneration and spell efficiency', 'Uncommon'),
('Scroll of Battle Tactics', 'Strategic combat and leadership skills', 'Combat', 3, 2500, '{"Basic Combat"}', 'Lead groups in combat and use advanced battle strategies', 'Rare'),
('Scroll of Nature''s Whisper', 'Communication with animals and plants', 'Magic', 2, 1500, '{}', 'Speak with animals and understand plant needs', 'Uncommon'),
('Scroll of Shadow Walking', 'Stealth and shadow manipulation techniques', 'Magic', 3, 3000, '{}', 'Move unseen and manipulate shadows', 'Rare'),
('Scroll of Time Perception', 'Understanding temporal magic and time flow', 'Magic', 5, 10000, '{"Elemental Basics", "Meditation Mastery"}', 'Perceive time anomalies and resist temporal effects', 'Legendary'),
('Scroll of Machine Speech', 'Communication with mechanical and artificial beings', 'Magic', 3, 2800, '{}', 'Interface with machines and understand artificial intelligence', 'Rare'),
('Scroll of Spirit Binding', 'Techniques for working with spirits and souls', 'Magic', 4, 6000, '{"Meditation Mastery"}', 'Communicate with spirits and perform soul magic', 'Epic');

-- Insert initial spells that can be learned
INSERT INTO spells (name, description, spell_school, element, power_level, mana_cost_to_learn, mana_cost_to_cast, requirements, effects, rarity) VALUES
('Fireball', 'A classic spell that hurls a ball of fire at enemies', 'Elemental', 'Fire', 2, 1000, 25, 'Level 1+', 'Deals fire damage to a single target', 'Common'),
('Lightning Bolt', 'Strikes a target with a bolt of lightning', 'Elemental', 'Air', 3, 1500, 35, 'Level 2+', 'Deals lightning damage with chance to stun', 'Uncommon'),
('Healing Light', 'Channels positive energy to heal wounds', 'Elemental', 'Light', 2, 1200, 30, 'Level 1+', 'Restores health to target', 'Common'),
('Shadowmeld', 'Allows caster to become one with shadows', 'Illusion', 'Shadow', 4, 3000, 45, 'Level 3+', 'Grants temporary invisibility', 'Rare'),
('Teleportation', 'Instantly transport to a nearby location', 'Transmutation', 'Void', 5, 5000, 60, 'Level 4+', 'Instantly move to visible location', 'Epic'),
('Mind Read', 'Peer into the thoughts of another being', 'Divination', 'Spirit', 3, 2500, 40, 'Level 2+', 'Read surface thoughts of target', 'Uncommon'),
('Summon Familiar', 'Calls a magical creature to serve as companion', 'Conjuration', 'Spirit', 3, 2000, 50, 'Level 2+', 'Summons a magical familiar for 1 hour', 'Uncommon'),
('Time Stop', 'Briefly stops time for everyone except the caster', 'Transmutation', 'Time', 9, 25000, 200, 'Level 8+, Time Affinity', 'Stops time for 6 seconds', 'Forbidden'),
('Meteor Strike', 'Calls down a meteor from the heavens', 'Elemental', 'Fire', 8, 15000, 120, 'Level 7+, Fire Mastery', 'Devastating area fire damage', 'Legendary'),
('Soul Drain', 'Drains life force from target to heal caster', 'Elemental', 'Shadow', 6, 8000, 80, 'Level 5+, Dark Affinity', 'Damages target and heals caster', 'Epic'),
('Dispel Magic', 'Removes magical effects from target', 'Transmutation', 'Void', 3, 1800, 35, 'Level 2+', 'Removes magical effects', 'Uncommon'),
('Ice Prison', 'Encases target in magical ice', 'Elemental', 'Water', 4, 3500, 50, 'Level 3+', 'Immobilizes target in ice', 'Rare'),
('Earthquake', 'Causes the ground to shake violently', 'Elemental', 'Earth', 6, 7000, 90, 'Level 5+, Earth Affinity', 'Area earth damage and knockdown', 'Epic'),
('Astral Projection', 'Separates spirit from body for exploration', 'Divination', 'Spirit', 7, 12000, 100, 'Level 6+, Spirit Affinity', 'Explore in spirit form', 'Legendary'),
('Reality Rift', 'Tears a hole in reality itself', 'Transmutation', 'Void', 10, 50000, 300, 'Level 10+, Void Mastery', 'Creates dangerous spatial anomaly', 'Forbidden');