package types

type (
	Kind  uint16
	Range uint8
)

const (
	// Ranges.
	Regular Range = iota
	Replaceable
	Ephemeral
	Addressable

	// Kinds.
	KindUserMetadata                 Kind = 0
	KindShortTextNote                Kind = 1
	KindRecommendRelay               Kind = 2
	KindFollows                      Kind = 3
	KindEncryptedDirectMessages      Kind = 4
	KindEventDeletionRequest         Kind = 5
	KindRepost                       Kind = 6
	KindReaction                     Kind = 7
	KindBadgeAward                   Kind = 8
	KindGroupChatMessage             Kind = 9
	KindGroupChatThreadedReply       Kind = 10
	KindGroupThread                  Kind = 11
	KindGroupThreadReply             Kind = 12
	KindSeal                         Kind = 13
	KindDirectMessage                Kind = 14
	KindGenericRepost                Kind = 16
	KindReactionToWebsite            Kind = 17
	KindChannelCreation              Kind = 40
	KindChannelMetadata              Kind = 41
	KindChannelMessage               Kind = 42
	KindChannelHideMessage           Kind = 43
	KindChannelMuteUser              Kind = 44
	KindChessPGN                     Kind = 64
	KindMergeRequests                Kind = 818
	KindBid                          Kind = 1021
	KindBidConfirmation              Kind = 1022
	KindOpenTimestamps               Kind = 1040
	KindGiftWrap                     Kind = 1059
	KindFileMetadata                 Kind = 1063
	KindLiveChatMessage              Kind = 1311
	KindPatches                      Kind = 1617
	KindIssues                       Kind = 1621
	KindReplies                      Kind = 1622
	KindStatus                       Kind = 1630 // to 1633. support: todo.
	KindProblemTracker               Kind = 1971
	KindReporting                    Kind = 1984
	KindLabel                        Kind = 1985
	KindRelayReviews                 Kind = 1986
	KindAIEmbeddingsVectorLists      Kind = 1987
	KindTorrent                      Kind = 2003
	KindTorrentComment               Kind = 2004
	KindCoinJoinPool                 Kind = 2022
	KindCommunityPostApproval        Kind = 4550
	KindJobRequest                   Kind = 5000 // to 5999. support: todo.
	KindJobResult                    Kind = 6000 // to 6999. support: todo.
	KindJobFeedback                  Kind = 7000
	KindGroupControlEvents           Kind = 9000 // to 9030. support: todo.
	KindZapGoal                      Kind = 9041
	KindTidalLogin                   Kind = 9467
	KindZapRequest                   Kind = 9734
	KindZap                          Kind = 9735
	KindHighlights                   Kind = 9802
	KindMuteList                     Kind = 10000
	KindPinList                      Kind = 10001
	KindRelayListMetadata            Kind = 10002
	KindBookmarkList                 Kind = 10003
	KindCommunitiesList              Kind = 10004
	KindPublicChatsList              Kind = 10005
	KindBlockedRelaysList            Kind = 10006
	KindSearchRelaysList             Kind = 10007
	KindUserGroups                   Kind = 10009
	KindInterestsList                Kind = 10015
	KindUserEmojiList                Kind = 10030
	KindRelayListToReceiveDMs        Kind = 10050
	KindUserServerList               Kind = 10063
	KindFileStorageServerList        Kind = 10096
	KindWalletInfo                   Kind = 13194
	KindLightningPubRPC              Kind = 21000
	KindClientAuthentication         Kind = 22242
	KindWalletRequest                Kind = 23194
	KindWalletResponse               Kind = 23195
	KindNostrConnect                 Kind = 24133
	KindBlobsStoredOnMediaServers    Kind = 24242
	KindHTTPAuth                     Kind = 27235
	KindFollowSets                   Kind = 30000
	KindGenericLists                 Kind = 30001
	KindRelaySets                    Kind = 30002
	KindBookmarkSets                 Kind = 30003
	KindCurationSets                 Kind = 30004
	KindVideoSets                    Kind = 30005
	KindKindMuteSets                 Kind = 30007
	KindProfileBadges                Kind = 30008
	KindBadgeDefinition              Kind = 30009
	KindInterestSets                 Kind = 30015
	KindCreateOrUpdateAStall         Kind = 30017
	KindCreateOrUpdateAProduct       Kind = 30018
	KindMarketplaceUIUx              Kind = 30019
	KindProductSoldAsAnAuction       Kind = 30020
	KindLongFormContent              Kind = 30023
	KindDraftLongFormContent         Kind = 30024
	KindEmojiSets                    Kind = 30030
	KindModularArticleHeader         Kind = 30040
	KindModularArticleContent        Kind = 30041
	KindReleaseArtifactSets          Kind = 30063
	KindApplicationSpecificData      Kind = 30078
	KindLiveEvent                    Kind = 30311
	KindUserStatuses                 Kind = 30315
	KindClassifiedListing            Kind = 30402
	KindDraftClassifiedListing       Kind = 30403
	KindRepositoryAnnouncements      Kind = 30617
	KindRepositoryStateAnnouncements Kind = 30618
	KindWikiArticle                  Kind = 30818
	KindRedirects                    Kind = 30819
	KindFeed                         Kind = 31890
	KindDateBasedCalendarEvent       Kind = 31922
	KindTimeBasedCalendarEvent       Kind = 31923
	KindCalendar                     Kind = 31924
	KindCalendarEventRSVP            Kind = 31925
	KindHandlerRecommendation        Kind = 31989
	KindHandlerInformation           Kind = 31990
	KindVideoEvent                   Kind = 34235
	KindShortFormPortraitVideoEvent  Kind = 34236
	KindVideoViewEvent               Kind = 34237
	KindCommunityDefinition          Kind = 34550
	KindGroupMetadataEvents          Kind = 39000 // to 39009. support: todo.
)

// IsRegular checks if the given kind is in Regular range.
func (k Kind) IsRegular() bool {
	return k < 10000 && k != 0 && k != 3
}

// IsReplaceable checks if the given kind is in Replaceable range.
func (k Kind) IsReplaceable() bool {
	return k == 0 || k == 3 ||
		(10000 <= k && k < 20000)
}

// IsEphemeral checks if the given kind is in Ephemeral range.
func (k Kind) IsEphemeral() bool {
	return 20000 <= k && k < 30000
}

// IsAddressable checks if the given kind is in Addressable range.
func (k Kind) IsAddressable() bool {
	return 30000 <= k && k < 40000
}

// Range returns the kind range based on NIP-01.
func (k Kind) Range() Range {
	if k.IsRegular() {
		return Regular
	} else if k.IsReplaceable() {
		return Replaceable
	} else if k.IsAddressable() {
		return Addressable
	}

	return Ephemeral
}
