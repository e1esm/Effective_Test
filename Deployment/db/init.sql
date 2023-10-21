CREATE USER el_esm WITH ENCRYPTED PASSWORD 'test_password';
GRANT insert on ALL TABLES IN SCHEMA public to el_esm;
GRANT update on ALL TABLES IN SCHEMA public to el_esm;
GRANT select on ALL TABLES IN SCHEMA public to el_esm;
GRANT CONNECT ON DATABASE people TO el_esm;
GRANT CREATE ON SCHEMA public to el_esm;
