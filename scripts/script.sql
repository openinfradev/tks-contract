CREATE DATABASE tks;
\c tks;
CREATE TABLE contracts
(
    id character varying(10) COLLATE pg_catalog."default" primary key,
    contractor_name character varying(50) COLLATE pg_catalog."default",
    available_services character varying(50)[] COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    created_at timestamp with time zone
);
CREATE UNIQUE INDEX idx_contractor_name ON contracts(contractor_name);
ALTER TABLE contracts CLUSTER ON idx_contractor_name;
INSERT INTO contracts(
	contractor_name, id, available_services, updated_at, created_at)
	VALUES ('tester', 'Pedcaa975', ARRAY['lma'], '2021-05-01'::timestamp, '2021-05-01'::timestamp);

CREATE TABLE resource_quota
(
    id uuid primary key,
    cpu bigint,
    memory bigint,
    block bigint,
    block_ssd bigint,
    fs bigint,
    fs_ssd bigint,
    contract_id character varying(10) COLLATE pg_catalog."default",
    updated_at timestamp with time zone,
    created_at timestamp with time zone
);
