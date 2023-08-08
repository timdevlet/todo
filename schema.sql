-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."todos" (
    "uuid" uuid NOT NULL DEFAULT gen_random_uuid(),
    "owner_uuid" uuid NOT NULL,
    "title" varchar(200) NOT NULL,
    "description" varchar(2000),
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "done_at" timestamptz
);

```
```sql
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."logs" (
    "uuid" uuid NOT NULL DEFAULT gen_random_uuid(),
    "payload" jsonb NOT NULL DEFAULT '{}'::jsonb,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "address" varchar(15)
);
