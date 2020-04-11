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
  deleted boolean NOT NULL,
  originator_cache_guid varchar(255),
  originator_client_item_id varchar(255),
  specifics bytea,
  data_type_id int,
  folder boolean NOT NULL,
  client_defined_unique_tag varchar(255),
  unique_position bytea,
  client_id varchar(255) not null references clients(id),
  unique (client_id, server_defined_unique_tag),
  unique (client_id, client_defined_unique_tag)
);
