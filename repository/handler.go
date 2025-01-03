package repository

import (
	"github.com/dezh-tech/immortal/infrastructure/database"
	"github.com/dezh-tech/immortal/types"
)

var KindToCollectionName = map[types.Kind]string{
	types.KindUserMetadata:                "user_metadatas",
	types.KindShortTextNote:               "short_text_notes",
	types.KindRecommendRelay:              "recommend_relays",
	types.KindFollows:                     "follows",
	types.KindEncryptedDirectMessages:     "encrypted_direct_messages",
	types.KindEventDeletionRequest:        "event_deletion_requests",
	types.KindRepost:                      "reposts",
	types.KindReaction:                    "reactions",
	types.KindBadgeAward:                  "badge_awards",
	types.KindGroupChatMessage:            "group_chat_messages",
	types.KindGroupChatThreadedReply:      "group_chat_threaded_replies",
	types.KindGroupThread:                 "group_threads",
	types.KindGroupThreadReply:            "group_thread_replies",
	types.KindSeal:                        "seals",
	types.KindDirectMessage:               "direct_messages",
	types.KindGenericRepost:               "generic_reposts",
	types.KindReactionToWebsite:           "reactions_to_websites",
	types.KindChannelCreation:             "channel_creations",
	types.KindChannelMetadata:             "channel_metadatas",
	types.KindChannelMessage:              "channel_messages",
	types.KindChannelHideMessage:          "channel_hide_messages",
	types.KindChannelMuteUser:             "channel_mute_users",
	types.KindChessPGN:                    "chess_pgn",
	types.KindMergeRequests:               "merge_requests",
	types.KindBid:                         "bids",
	types.KindBidConfirmation:             "bid_confirmations",
	types.KindOpenTimestamps:              "open_timestamps",
	types.KindGiftWrap:                    "gift_wraps",
	types.KindFileMetadata:                "file_metadatas",
	types.KindLiveChatMessage:             "live_chat_messages",
	types.KindPatches:                     "patches",
	types.KindIssues:                      "issues",
	types.KindReplies:                     "replies",
	types.KindStatus:                      "status",
	types.KindProblemTracker:              "problem_trackers",
	types.KindReporting:                   "reportings",
	types.KindLabel:                       "labels",
	types.KindRelayReviews:                "relay_reviews",
	types.KindAIEmbeddingsVectorLists:     "ai_embeddings_vector_lists",
	types.KindTorrent:                     "torrents",
	types.KindTorrentComment:              "torrent_comments",
	types.KindCoinJoinPool:                "coin_join_pools",
	types.KindCommunityPostApproval:       "community_post_approvals",
	types.KindJobRequest:                  "dvm",
	types.KindJobResult:                   "dvm",
	types.KindJobFeedback:                 "dvm",
	types.KindGroups:                      "groups",
	types.KindZapGoal:                     "zap_goals",
	types.KindTidalLogin:                  "tidal_logins",
	types.KindZapRequest:                  "zap_requests",
	types.KindZap:                         "zaps",
	types.KindHighlights:                  "highlights",
	types.KindMuteList:                    "mute_lists",
	types.KindPinList:                     "pin_lists",
	types.KindRelayListMetadata:           "relay_list_metadatas",
	types.KindBookmarkList:                "bookmark_lists",
	types.KindCommunitiesList:             "communities_lists",
	types.KindPublicChatsList:             "public_chats_lists",
	types.KindBlockedRelaysList:           "blocked_relays_lists",
	types.KindSearchRelaysList:            "search_relays_lists",
	types.KindUserGroups:                  "user_groups",
	types.KindInterestsList:               "interests_lists",
	types.KindUserEmojiList:               "user_emoji_lists",
	types.KindRelayListToReceiveDMs:       "relay_list_to_receive_dms",
	types.KindUserServerList:              "user_server_lists",
	types.KindFileStorageServerList:       "file_storage_server_lists",
	types.KindWalletInfo:                  "wallet_infos",
	types.KindLightningPubRPC:             "lightning_pub_rpcs",
	types.KindClientAuthentication:        "client_authentications",
	types.KindWalletRequest:               "wallet_requests",
	types.KindWalletResponse:              "wallet_responses",
	types.KindNostrConnect:                "nostr_connects",
	types.KindBlobsStoredOnMediaServers:   "blobs_stored_on_media_servers",
	types.KindHTTPAuth:                    "http_auths",
	types.KindFollowSets:                  "follow_sets",
	types.KindGenericLists:                "generic_lists",
	types.KindRelaySets:                   "relay_sets",
	types.KindBookmarkSets:                "bookmark_sets",
	types.KindCurationSets:                "curation_sets",
	types.KindVideoSets:                   "video_sets",
	types.KindKindMuteSets:                "kind_mute_sets",
	types.KindProfileBadges:               "profile_badges",
	types.KindBadgeDefinition:             "badge_definitions",
	types.KindLiveEvent:                   "live_events",
	types.KindShortFormPortraitVideoEvent: "short_form_portrait_video_events",
	types.KindVideoViewEvent:              "video_view_events",
	types.KindCommunityDefinition:         "community_definitions",
	types.KindGroupsMetadata:              "groups_metadata",
}

func getCollectionName(k types.Kind) string {
	collName, ok := KindToCollectionName[k]
	if ok {
		return collName
	}

	if k >= 9000 && k <= 9030 {
		return "groups"
	}

	if k >= 1630 && k <= 1633 {
		return "status"
	}

	if k >= 39000 && k <= 39009 {
		return "groups_metadata"
	}

	if k >= 5000 && k <= 5999 || k >= 6000 && k <= 6999 || k == 7000 {
		return "dvm"
	}

	return "unknown"
}

type Handler struct {
	db     *database.Database
	config Config
}

func New(db *database.Database, cfg Config) *Handler {
	return &Handler{
		db:     db,
		config: cfg,
	}
}
