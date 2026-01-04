CREATE TABLE mountains (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  elevation_m INTEGER NOT NULL,
  province TEXT,
  latitude DOUBLE PRECISION,
  longitude DOUBLE PRECISION,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE routes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  mountain_id UUID NOT NULL REFERENCES mountains(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  start_point_name TEXT,
  difficulty TEXT,
  length_m INTEGER,
  elevation_gain_m INTEGER,
  is_official BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  UNIQUE (mountain_id, name)
);

CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  order_index INTEGER NOT NULL,
  altitude_m INTEGER,
  distance_from_start_m INTEGER,
  segment_to_next_m INTEGER,
  distance_to_summit_m INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  UNIQUE (route_id, order_index)
);

CREATE INDEX idx_routes_mountain_id ON routes(mountain_id);
CREATE INDEX idx_posts_route_id ON posts(route_id);
