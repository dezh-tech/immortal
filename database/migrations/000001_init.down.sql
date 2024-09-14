-- Drop foreign key constraints first
ALTER TABLE public.reactions DROP CONSTRAINT IF EXISTS FKreactions_users_metadata;
ALTER TABLE public.reactions DROP CONSTRAINT IF EXISTS FKreactions_text_notes;
ALTER TABLE public.follow_list DROP CONSTRAINT IF EXISTS FKfollow_list_following;
ALTER TABLE public.follow_list DROP CONSTRAINT IF EXISTS FKfollow_list_follower;
ALTER TABLE public.text_notes DROP CONSTRAINT IF EXISTS FKtext_notes_users_metadata;

-- Drop indexes
DROP INDEX IF EXISTS public.reactions_text_notesid;

-- Drop tables
DROP TABLE IF EXISTS public.follow_list CASCADE;
DROP TABLE IF EXISTS public.reactions CASCADE;
DROP TABLE IF EXISTS public.text_notes CASCADE;
DROP TABLE IF EXISTS public.users_metadata CASCADE;
