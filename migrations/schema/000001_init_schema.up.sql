---SITES---
CREATE TABLE public.sites
(
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "description" TEXT,
  "address" VARCHAR NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

---ROLES---
CREATE TABLE public.roles
(
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "code" VARCHAR NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

---USERS---
CREATE TABLE public.users
(
  "id" BIGSERIAL PRIMARY KEY,
  "first_name" VARCHAR NOT NULL,
  "last_name" VARCHAR,
  "full_name" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL,
  "phone_number" VARCHAR,
  "hashed_password" VARCHAR NOT NULL,
  "role_id" BIGINT NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

---USER_ACTION_LOGS---
CREATE TABLE public.user_action_logs
(
  "id" BIGSERIAL PRIMARY KEY,
  "target_id" BIGINT NOT NULL,
  "user_id" BIGINT NOT NULL,
  "action" VARCHAR NOT NULL,
  "target_type" VARCHAR NOT NULL,
  "old_data" JSONB,
  "new_data" JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

---INSPECTIONS---
CREATE TABLE public.inspections
(
  "id" BIGSERIAL PRIMARY KEY,
  "site_id" BIGINT NOT NULL,
  "user_id" BIGINT NOT NULL,
  "date" TIMESTAMP WITH TIME ZONE NOT NULL,
  "status" VARCHAR NOT NULL,
  "remark" VARCHAR,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

---ISSUES---
CREATE TABLE public.issues
(
  "id" BIGSERIAL PRIMARY KEY,
  "site_id" BIGINT NOT NULL,
  "user_id" BIGINT NOT NULL,
  "inspection_id" BIGINT NOT NULL,
  "status" VARCHAR NOT NULL,
  "description" TEXT,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

---SPARE_PARTS---
CREATE TABLE public.spare_parts
(
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "code" VARCHAR NOT NULL,
  "in_stock" BIGINT DEFAULT 0,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

---ISSUES_SPARE_PARTS---
CREATE TABLE public.issue_spare_parts
(
  "id" BIGSERIAL PRIMARY KEY,
  "issue_id" BIGINT NOT NULL,
  "spare_part_id" BIGINT NOT NULL,
  "quantity" BIGINT DEFAULT 1,
  "status" VARCHAR NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

---PURCHASE_REQUESTS---
CREATE TABLE public.purchase_requests
(
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "description" TEXT,
  "order_date" TIMESTAMP WITH TIME ZONE NOT NULL,
  "status" VARCHAR NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

---PURCHASE_SPARE_PARTS---
CREATE TABLE public.purchase_spare_parts
(
  "id" BIGSERIAL PRIMARY KEY,
  "purchase_request_id" BIGINT NOT NULL,
  "spare_part_id" BIGINT NOT NULL,
  "quantity" BIGINT DEFAULT 1,
  "spare_part_name" VARCHAR,
  "spare_part_code" VARCHAR,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "user_action_logs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "inspections" ADD FOREIGN KEY ("site_id") REFERENCES "sites" ("id");

ALTER TABLE "inspections" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "issues" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "issues" ADD FOREIGN KEY ("site_id") REFERENCES "sites" ("id");

ALTER TABLE "issues" ADD FOREIGN KEY ("inspection_id") REFERENCES "inspections" ("id");

ALTER TABLE "issue_spare_parts" ADD FOREIGN KEY ("issue_id") REFERENCES "issues" ("id");

ALTER TABLE "issue_spare_parts" ADD FOREIGN KEY ("spare_part_id") REFERENCES "spare_parts" ("id");

ALTER TABLE "purchase_requests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "purchase_spare_parts" ADD FOREIGN KEY ("purchase_request_id") REFERENCES "purchase_requests" ("id");

ALTER TABLE "purchase_spare_parts" ADD FOREIGN KEY ("spare_part_id") REFERENCES "spare_parts" ("id");