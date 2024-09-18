BEGIN;

ALTER TABLE public.reactions
ADD COLUMN r_tags TEXT[], 
ADD COLUMN content VARCHAR(8);

COMMIT;
