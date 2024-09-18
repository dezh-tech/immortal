BEGIN;

CREATE TABLE public.follow_list (
    follower CHAR(64),  -- Optional, so no NOT NULL
    following CHAR(64),  -- Optional, so no NOT NULL
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- You may need a trigger for automatic updates
    deleted_at TIMESTAMP,
    PRIMARY KEY (follower, following)
);

CREATE TABLE public.reactions (
    id CHAR(64) NOT NULL,
    text_notesid CHAR(64),  -- Optional, so no NOT NULL
    users_metadatapub_key CHAR(64),  -- Optional, so no NOT NULL
    e_tags TEXT[],  -- Assuming e is an array of text
    p_tags TEXT[],  -- Assuming p is an array of text
    a_tags TEXT[],  -- Assuming a is an array of text
    k_tags TEXT[],  -- Assuming k is an array of text
    event TEXT NOT NULL,
    event_created_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE public.text_notes (
    id CHAR(64) NOT NULL,
    e_tags TEXT[],  -- Assuming e is an array of text
    p_tags TEXT[],  -- Assuming p is an array of text
    content VARCHAR(65535),
    event TEXT NOT NULL UNIQUE,
    event_created_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    users_metadatapub_key CHAR(64),  -- Optional, so no NOT NULL
    PRIMARY KEY (id)
);

CREATE TABLE public.users_metadata (
    pub_key CHAR(64) NOT NULL,
    event_created_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    content VARCHAR(65535),
    follow_list_event TEXT,
    PRIMARY KEY (pub_key)
);

-- Index for better query performance
CREATE INDEX reactions_text_notesid ON public.reactions (text_notesid);

-- Index for deleted entities in follow_list
CREATE INDEX follow_list_deleted_at_idx ON public.follow_list (deleted_at);

-- Index for deleted entities in reactions
CREATE INDEX reactions_deleted_at_idx ON public.reactions (deleted_at);

-- Index for deleted entities in text_notes
CREATE INDEX text_notes_deleted_at_idx ON public.text_notes (deleted_at);

-- Index for deleted entities in users_metadata
CREATE INDEX users_metadata_deleted_at_idx ON public.users_metadata (deleted_at);

-- Composite index on follow_list by follower and following
CREATE INDEX idx_follow_list_follower_following ON public.follow_list (follower, following);

COMMIT;
