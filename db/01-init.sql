-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS expenses_id_seq;

-- Table Definition
CREATE TABLE "expenses" (
    "id" int4 NOT NULL DEFAULT nextval('expenses_id_seq'::regclass),
    "title" text,
    "amount" float,
    "note" text,
    "tags" text[],
    PRIMARY KEY ("id")
);

INSERT INTO "expenses" ("title", "amount", "note", "tags") VALUES ('buy a new phone', 39000, 'buy a new phone', ARRAY ['gadget', 'shopping']);