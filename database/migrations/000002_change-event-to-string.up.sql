BEGIN;

-- Alter event field in reactions table from JSONB to TEXT
ALTER TABLE public.reactions
    ALTER COLUMN event TYPE TEXT USING event::TEXT;

-- Alter event field in text_notes table from JSONB to TEXT
ALTER TABLE public.text_notes
    ALTER COLUMN event TYPE TEXT USING event::TEXT;

-- Alter follow_list_event field in users_metadata table from JSONB to TEXT
ALTER TABLE public.users_metadata
    ALTER COLUMN follow_list_event TYPE TEXT USING follow_list_event::TEXT;

COMMIT;
