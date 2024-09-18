BEGIN;

ALTER TABLE public.reactions
DROP COLUMN r_tag,
DROP COLUMN content;

COMMIT;
