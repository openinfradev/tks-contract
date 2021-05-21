CREATE DATABASE tks;
\c tks;
CREATE TABLE contract
(
    contractor_name character varying(50) COLLATE pg_catalog."default",
    contract_id uuid primary key,
    available_services character varying(10)[] COLLATE pg_catalog."default",
    csp_id uuid,
    last_updated_ts timestamp with time zone
);

INSERT INTO contract(
	contractor_name, contract_id, available_services, csp_id, last_updated_ts)
	VALUES ('tester', 'edcaa975-dde4-4c4d-94f7-36bc38fe7064', ARRAY['lma'], '3390f92b-0da8-4628-83e2-e266b1928e11', '2021-05-01'::timestamp);