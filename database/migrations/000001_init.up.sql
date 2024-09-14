BEGIN;

CREATE TABLE public.follow_list (
    follower CHAR(32),  -- Optional, so no NOT NULL
    following CHAR(32),  -- Optional, so no NOT NULL
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- You may need a trigger for automatic updates
    deleted_at TIMESTAMP,
    PRIMARY KEY (follower, following),
    CONSTRAINT follower_following UNIQUE (follower, following)
);

CREATE TABLE public.reactions (
    id UUID NOT NULL,
    text_notesid UUID,  -- Optional, so no NOT NULL
    users_metadatapub_key CHAR(32),  -- Optional, so no NOT NULL
    e TEXT[],  -- Assuming e is an array of text
    p TEXT[],  -- Assuming p is an array of text
    a TEXT[],  -- Assuming a is an array of text
    event JSONB NOT NULL,
    k TEXT[],  -- Assuming k is an array of text
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE public.text_notes (
    id UUID NOT NULL,
    e TEXT[],  -- Assuming e is an array of text
    p TEXT[],  -- Assuming p is an array of text
    content VARCHAR(65535),
    event JSONB NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    users_metadatapub_key CHAR(32),  -- Optional, so no NOT NULL
    PRIMARY KEY (id)
);

CREATE TABLE public.users_metadata (
    pub_key CHAR(32) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    content VARCHAR(255),
    follow_list_event JSONB,
    PRIMARY KEY (pub_key)
);

-- Index for better query performance
CREATE INDEX reactions_text_notesid ON public.reactions (text_notesid);

-- Foreign key constraints with optional references
ALTER TABLE public.reactions
    ADD CONSTRAINT FKreactions_users_metadata
    FOREIGN KEY (users_metadatapub_key) REFERENCES public.users_metadata (pub_key);

ALTER TABLE public.follow_list
    ADD CONSTRAINT FKfollow_list_following
    FOREIGN KEY (following) REFERENCES public.users_metadata (pub_key),
    ADD CONSTRAINT FKfollow_list_follower
    FOREIGN KEY (follower) REFERENCES public.users_metadata (pub_key);

ALTER TABLE public.reactions
    ADD CONSTRAINT FKreactions_text_notes
    FOREIGN KEY (text_notesid) REFERENCES public.text_notes (id);

ALTER TABLE public.text_notes
    ADD CONSTRAINT FKtext_notes_users_metadata
    FOREIGN KEY (users_metadatapub_key) REFERENCES public.users_metadata (pub_key);

COMMIT;