CREATE SEQUENCE IF NOT EXISTS tasks_id_seq
  AS INTEGER
  MAXVALUE 2147483647;

CREATE TABLE IF NOT EXISTS tasks
(
  id          INTEGER DEFAULT nextval('tasks_id_seq'::regclass) NOT NULL CONSTRAINT tasks_pkey PRIMARY KEY ,

  name        VARCHAR(50),
  details     TEXT,
  resolved_at TIMESTAMP WITH TIME ZONE,

  created_at  TIMESTAMP WITH TIME ZONE,
  updated_at  TIMESTAMP WITH TIME ZONE
);

-- create index idx_tasks_relation_id on tasks (relation_id);