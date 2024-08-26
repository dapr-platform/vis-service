-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE o_topology(
                           id VARCHAR(32) NOT NULL,
                           created_by VARCHAR(32),
                           created_time TIMESTAMP,
                           updated_by VARCHAR(32),
                           updated_time TIMESTAMP,
                           group_id VARCHAR(255),
                           file_data TEXT,
                           key_name VARCHAR(255),
                           name VARCHAR(255),
                           parent_id VARCHAR(255),
                           remark VARCHAR(255),
                           type VARCHAR(255),
                           PRIMARY KEY (id)
);

COMMENT ON TABLE o_topology IS '组态';
COMMENT ON COLUMN o_topology.id IS '唯一标识';
COMMENT ON COLUMN o_topology.created_by IS '创建人';
COMMENT ON COLUMN o_topology.created_time IS '创建时间';
COMMENT ON COLUMN o_topology.updated_by IS '更新人';
COMMENT ON COLUMN o_topology.updated_time IS '更新时间';
COMMENT ON COLUMN o_topology.group_id IS '分组id';
COMMENT ON COLUMN o_topology.file_data IS '文件内容';
COMMENT ON COLUMN o_topology.key_name IS '标识';
COMMENT ON COLUMN o_topology.name IS '名称';
COMMENT ON COLUMN o_topology.parent_id IS '父id';
COMMENT ON COLUMN o_topology.remark IS '备注';
COMMENT ON COLUMN o_topology.type IS '类型';

create table o_dashboard(
    id varchar(255) not null,
    name varchar(255) not null,
    app varchar(255) not null,
    layout text not null,
    created_by VARCHAR(32),
    created_time TIMESTAMP,
    updated_by VARCHAR(32),
    updated_time TIMESTAMP,
    primary key (id)
);
COMMENT ON COLUMN o_dashboard.id IS 'id,(app+"_"+name)';
COMMENT ON COLUMN o_dashboard.name IS '名称';
COMMENT ON COLUMN o_dashboard.app IS '应用';
COMMENT ON COLUMN o_dashboard.layout IS '布局';
-- +goose Down
-- +goose StatementBegin
    drop table if exists o_topology;
    drop table if exists o_dashboard;
-- +goose StatementEnd
