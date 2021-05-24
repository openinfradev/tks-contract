CREATE DATABASE tks;
\c tks;
CREATE TABLE contracts
(
    contractor_name character varying(50) COLLATE pg_catalog."default",
    id uuid primary key,
    available_services character varying(50)[] COLLATE pg_catalog."default",
    csp_id uuid,
    updated_at timestamp with time zone,
    created_at timestamp with time zone
);
CREATE UNIQUE INDEX idx_contractor_name ON contracts(contractor_name);
ALTER TABLE contracts CLUSTER ON idx_contractor_name;
INSERT INTO contracts(
	contractor_name, id, available_services, csp_id, updated_at, created_at)
	VALUES ('tester', 'edcaa975-dde4-4c4d-94f7-36bc38fe7064', ARRAY['lma'], '3390f92b-0da8-4628-83e2-e266b1928e11', '2021-05-01'::timestamp, '2021-05-01'::timestamp);