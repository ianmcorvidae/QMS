BEGIN;

SET search_path = public, pg_catalog;

--
-- A table used to store resource usage limits for users and plans.
--
CREATE TABLE IF NOT EXISTS quotas (
    id uuid NOT NULL DEFAULT uuid_generate_v1(),
    quota numeric NOT NULL,
    resource_type_id uuid NOT NULL,
    user_plan_id uuid NOT NULL,
    created_by text NOT NULL,
    created_at timestamp with timezone NOT NULL,
    last_modified_by text NOT NULL,
    last_modified_at timestamp with time zone NOT NULL,
    FOREIGN KEY (resource_type_id) REFERENCES resource_types(id) ON DELETE CASCADE,
    FOREIGN KEY (user_plan_id) REFERENCES user_plans(id) ON DELETE CASCADE,
    PRIMARY KEY (id)
);

--
-- There can only be one quota for each user plan and resource type.
--
CREATE UNIQUE INDEX IF NOT EXISTS quotas_resource_type_user_plan_index
    ON quotas (resource_type_id, user_plan_id);

--
-- A trigger to set the created_by field when a new row is added to the quotas table.
--
DROP TRIGGER IF EXISTS quotas_created_by_trigger ON quotas CASCADE;
CREATE TRIGGER quotas_created_by_trigger
    BEFORE INSERT ON quotas
    FOR EACH ROW
    EXECUTE PROCEDURE insert_username(created_by);

--
-- A trigger to set the created_at field when a new row is added to the quotas table.
--
DROP TRIGGER IF EXISTS quotas_created_at_trigger ON quotas CASCADE;
CREATE TRIGGER quotas_created_at_trigger
    BEFORE INSERT ON quotas
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(created_at);

--
-- A trigger to set the last_modified_by field when a row is added to the quotas table.
--
DROP TRIGGER IF EXISTS quotas_last_modified_by_insert_trigger ON quotas CASCADE;
CREATE TRIGGER quotas_last_modified_by_insert_trigger
    BEFORE INSERT ON quotas
    FOR EACH ROW
    EXECUTE PROCEDURE insert_username(last_modified_by);

--
-- A trigger to set the last_modified_at field when a row is added to the quotas table.
--
DROP TRIGGER IF EXISTS quotas_last_modified_at_insert_trigger ON quotas CASCADE;
CREATE TRIGGER quotas_last_modified_at_insert_trigger
    BEFORE UPDATE ON quotas
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(last_modified_at);

--
-- A trigger to set the last_modified_by field when a row is modified in the quotas table.
--
DROP TRIGGER IF EXISTS quotas_last_modified_by_trigger ON quotas CASCADE;
CREATE TRIGGER quotas_last_modified_by_trigger
    BEFORE UPDATE ON quotas
    FOR EACH ROW
    EXECUTE PROCEDURE insert_username(last_modified_by);

--
-- A trigger to set the last_modified_at field when a row is modified in the quotas table.
--
DROP TRIGGER IF EXISTS quotas_last_modified_at_trigger ON quotas CASCADE;
CREATE TRIGGER quotas_last_modified_at_trigger
    BEFORE UPDATE ON quotas
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(last_modified_at);

--
-- A table to track a user's resource usage.
--
CREATE TABLE IF NOT EXISTS usages (
    id uuid NOT NULL DEFAULT uuid_generate_v1(),
    "usage" numeric NOT NULL
    resource_type_id uuid NOT NULL,
    user_plan_id uuid NOT NULL,
    created_by text NOT NULL,
    created_at timestamp with timezone NOT NULL,
    last_modified_by text NOT NULL,
    last_modified_at timestamp with time zone NOT NULL,
    FOREIGN KEY (resource_type_id) REFERENCES resource_types(id) ON DELETE CASCADE,
    FOREIGN KEY (user_plan_id) REFERENCES user_plans(id) ON DELETE CASCADE,
    PRIMARY KEY (id)
);

--
-- There can only be one usage for each user plan and resource type.
--
CREATE UNIQUE INDEX IF NOT EXISTS usages_resource_type_user_plan_index
    ON usages (resource_type_id, user_plan_id);

--
-- A trigger to set the created_by field when a new row is added to the usages table.
--
DROP TRIGGER IF EXISTS usages_created_by_trigger ON usages CASCADE;
CREATE TRIGGER usages_created_by_trigger
    BEFORE INSERT ON usages
    FOR EACH ROW
    EXECUTE PROCEDURE insert_username(created_by);

--
-- A trigger to set the created_at field when a new row is added to the usages table.
--
DROP TRIGGER IF EXISTS usages_created_at_trigger ON usages CASCADE;
CREATE TRIGGER usages_created_at_trigger
    BEFORE INSERT ON usages
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(created_at);

--
-- A trigger to set the last_modified_by field when a row is added to the usages table.
--
DROP TRIGGER IF EXISTS usages_last_modified_by_insert_trigger ON usages CASCADE;
CREATE TRIGGER usages_last_modified_by_insert_trigger
    BEFORE INSERT ON usages
    FOR EACH ROW
    EXECUTE PROCEDURE insert_username(last_modified_by);

--
-- A trigger to set the last_modified_at field when a row is added to the usages table.
--
DROP TRIGGER IF EXISTS usages_last_modified_at_insert_trigger ON usages CASCADE;
CREATE TRIGGER usages_last_modified_at_insert_trigger
    BEFORE UPDATE ON usages
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(last_modified_at);

--
-- A trigger to set the last_modified_by field when a row is modified in the usages table.
--
DROP TRIGGER IF EXISTS usages_last_modified_by_trigger ON usages CASCADE;
CREATE TRIGGER usages_last_modified_by_trigger
    BEFORE UPDATE ON usages
    FOR EACH ROW
    EXECUTE PROCEDURE insert_username(last_modified_by);

--
-- A trigger to set the last_modified_at field when a row is modified in the usages table.
--
DROP TRIGGER IF EXISTS usages_last_modified_at_trigger ON usages CASCADE;
CREATE TRIGGER usages_last_modified_at_trigger
    BEFORE UPDATE ON usages
    FOR EACH ROW
    EXECUTE PROCEDURE moddatetime(last_modified_at);

--
-- A table listing the types of updates that can be performed on a usage value.
--
CREATE TABLE IF NOT EXISTS update_operations (
    id uuid NOT NULL DEFAULT uuid_generate_v1(),
    "name" text NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

--
-- Tracked metrics is an enumeration indicating the types of values for which updates are tracked in the updates table.
--
