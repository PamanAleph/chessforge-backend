CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS game_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bot_level INTEGER NOT NULL,
    result TEXT CHECK (result IN ('white', 'black', 'draw', 'ongoing')) DEFAULT 'ongoing',
    started_at TIMESTAMP DEFAULT now(),
    ended_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS moves (
    id SERIAL PRIMARY KEY,
    game_id UUID REFERENCES game_sessions(id) ON DELETE CASCADE,
    move_number INTEGER,
    color TEXT CHECK (color IN ('white', 'black')),
    from_square TEXT,
    to_square TEXT,
    san TEXT,
    fen TEXT,
    created_at TIMESTAMP DEFAULT now()
);
