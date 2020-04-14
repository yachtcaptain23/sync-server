CREATE TABLE clients (
  id varchar(255) PRIMARY KEY,
  token varchar(255) NOT NULL,
  expire_at bigint NOT NULL
);

CREATE TABLE sync_entities (
  id varchar(255) PRIMARY KEY,
  parent_id varchar(255),
  old_parent_id varchar(255),
  version bigint NOT NULL,
  mtime bigint NOT NULL,
  ctime bigint NOT NULL,
  name VARCHAR(255),
  non_unique_name varchar(255),
  server_defined_unique_tag varchar(255),
  deleted_at bigint DEFAULT NULL,
  originator_cache_guid varchar(255),
  originator_client_item_id varchar(255),
  specifics bytea NOT NULL,
  data_type_id int NOT NULL,
  folder boolean NOT NULL,
  client_defined_unique_tag varchar(255),
  unique_position bytea NOT NULL,
  client_id varchar(255) NOT NULL references clients(id)
);

CREATE UNIQUE INDEX unique_server_tag ON sync_entities(client_id, server_defined_unique_tag) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX unique_client_tag ON sync_entities(client_id, client_defined_unique_tag) WHERE deleted_at IS NULL;
