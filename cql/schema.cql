CREATE KEYSPACE IF NOT EXISTS reader WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1}; 
CREATE TABLE IF NOT EXISTS reader.contacts (mobile text PRIMARY KEY, name text, email text); 
CREATE TABLE IF NOT EXISTS reader.sources (source text PRIMARY KEY); 
