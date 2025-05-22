CREATE TABLE image_scans (
    id SERIAL PRIMARY KEY,
    image_name VARCHAR(255) NOT NULL,
    scan_start_time TIMESTAMP WITH TIME ZONE,
    scan_finish_time TIMESTAMP WITH TIME ZONE,
    scan_result TEXT,
    files_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);