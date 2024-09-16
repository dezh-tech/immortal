BEGIN;

-- Revert event field in reactions table from TEXT to JSONB
ALTER TABLE public.reactions
    ALTER COLUMN event TYPE JSONB USING event::JSONB;

-- Revert event field in text_notes table from TEXT to JSONB
ALTER TABLE public.text_notes
    ALTER COLUMN event TYPE JSONB USING event::JSONB;

-- Revert follow_list_event field in users_metadata table from TEXT to JSONB
ALTER TABLE public.users_metadata
    ALTER COLUMN follow_list_event TYPE JSONB USING follow_list_event::JSONB;

COMMIT;
