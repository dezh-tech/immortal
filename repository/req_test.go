package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	grpcclient "github.com/dezh-tech/immortal/infrastructure/grpc_client/gen"
	infra "github.com/dezh-tech/immortal/infrastructure/meilisearch"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"

	meilisearchGo "github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/meilisearch"
)

type MockGRPC struct {
	mock.Mock
}

func (m *MockGRPC) UpdateParameters(ctx context.Context, newParams *grpcclient.GetParametersResponse) error {
	args := m.Called(ctx, newParams)

	return args.Error(0)
}

func (m *MockGRPC) RegisterService(ctx context.Context, port, region string) (*grpcclient.RegisterServiceResponse, error) {
	args := m.Called(ctx, port, region)

	return args.Get(0).(*grpcclient.RegisterServiceResponse), args.Error(1)
}

func (m *MockGRPC) GetParameters(ctx context.Context) (*grpcclient.GetParametersResponse, error) {
	args := m.Called(ctx)

	return args.Get(0).(*grpcclient.GetParametersResponse), args.Error(1)
}

func (m *MockGRPC) AddLog(ctx context.Context, msg, stack string) (*grpcclient.AddLogResponse, error) {
	args := m.Called(ctx, msg, stack)

	return args.Get(0).(*grpcclient.AddLogResponse), args.Error(1)
}

func (m *MockGRPC) SetID(id string) {
	m.Called(id)
}

var (
	meiliAPIKey = os.Getenv("MEILI_API_KEY")
	currentTime = time.Now().Unix()
	testEvents  = []event.Event{
		{
			ID:        "event1",
			PublicKey: "author1",
			CreatedAt: currentTime - 500,
			Kind:      types.KindShortTextNote,
			Tags:      types.Tags{{"category", "sports"}, {"mood", "happy"}},
			Content:   "This is a test event about sports",
			Signature: "sig1",
		},
		{
			ID:        "event2",
			PublicKey: "author2",
			CreatedAt: currentTime - 1000,
			Kind:      types.KindDirectMessage,
			Tags:      types.Tags{{"category", "news"}},
			Content:   "Breaking news: Go is awesome!",
			Signature: "sig2",
		},
		{
			ID:        "event3",
			PublicKey: "author3",
			CreatedAt: currentTime - 2000,
			Kind:      types.KindShortTextNote,
			Tags:      types.Tags{{"category", "sports"}, {"importance", "high"}},
			Content:   "sport finals are coming up soon!",
			Signature: "sig3",
		},
		{
			ID:        "event4",
			PublicKey: "author4",
			CreatedAt: currentTime - 300,
			Kind:      types.KindDirectMessage,
			Tags:      types.Tags{{"mood", "excited"}},
			Content:   "Can't wait for the upcoming concert!",
			Signature: "sig4",
		},
		{
			ID:        "event5",
			PublicKey: "author5",
			CreatedAt: currentTime - 400,
			Kind:      types.KindShortTextNote,
			Tags:      types.Tags{{"topic", "AI"}, {"tech", "future"}},
			Content:   "The future of AI looks promising!",
		},
		{
			ID:        "event6",
			PublicKey: "author6",
			CreatedAt: currentTime - 800,
			Kind:      types.KindShortTextNote,
			Tags:      types.Tags{{"category", "testing"}},
			Content:   "Sorting test event A",
			Signature: "sig6",
		},
		{
			ID:        "event7",
			PublicKey: "author7",
			CreatedAt: currentTime - 800,
			Kind:      types.KindShortTextNote,
			Tags:      types.Tags{{"category", "testing"}},
			Content:   "Sorting test event B",
			Signature: "sig7",
		},
		{
			ID:        "event8",
			PublicKey: "author8",
			CreatedAt: currentTime - 900,
			Kind:      types.KindGiftWrap,
			Tags:      types.Tags{{"p", "pubkey"}},
			Content:   "GiftWrap event that should be included",
			Signature: "sig8",
		},
		{
			ID:        "event9",
			PublicKey: "author9",
			CreatedAt: currentTime - 950,
			Kind:      types.KindGiftWrap,
			Tags:      types.Tags{{"p", "another_pubkey"}},
			Content:   "GiftWrap event that should be excluded",
			Signature: "sig9",
		},
	}
)

type requestHandlerTest struct {
	name     string
	filter   filter.Filter
	expected []string
}

var testCases = []requestHandlerTest{
	{
		name: "Multiple Tags Filter",
		filter: filter.Filter{
			Tags: map[string][]string{
				"category": {"sports"},
				"mood":     {"happy"},
			},
			Limit: 5,
		},
		expected: []string{"event1"},
	},
	{
		name: "Since & Until Filter",
		filter: filter.Filter{
			Since: currentTime - 950,
			Until: currentTime - 350,
			Limit: 5,
		},
		expected: []string{"event5", "event1", "event6", "event7", "event8"},
	},
	{
		name: "Multiple Kinds Filter",
		filter: filter.Filter{
			Kinds: []types.Kind{types.KindShortTextNote, types.KindDirectMessage},
			Limit: 7,
		},
		expected: []string{"event4", "event5", "event1", "event6", "event7", "event2", "event3"},
	},
	{
		name: "Multiple Authors Filter",
		filter: filter.Filter{
			Authors: []string{"author1", "author3"},
			Limit:   5,
		},
		expected: []string{"event1", "event3"},
	},
	{
		name: "Search in Content - Sports",
		filter: filter.Filter{
			Search: "sports",
			Limit:  5,
		},
		expected: []string{"event1", "event3"},
	},
	{
		name: "Search in Content - Concert",
		filter: filter.Filter{
			Search: "concert",
			Limit:  5,
		},
		expected: []string{"event4"},
	},
	{
		name: "Search in Content - AI",
		filter: filter.Filter{
			Search: "AI",
			Limit:  5,
		},
		expected: []string{"event5"},
	},
	{
		name: "Search in Content - Breaking News",
		filter: filter.Filter{
			Search: "Breaking news",
			Limit:  5,
		},
		expected: []string{"event2"},
	},
	{
		name: "Sorting by id:asc if created_at is the same",
		filter: filter.Filter{
			Tags: map[string][]string{
				"category": {"testing"},
			},
			Limit: 5,
		},
		expected: []string{"event6", "event7"},
	},
	{
		name: "Filter GiftWrap events by pubkey in tag",
		filter: filter.Filter{
			Kinds: []types.Kind{types.KindGiftWrap},
			Limit: 5,
		},
		expected: []string{"event8"},
	},
}

func TestHandleReq(t *testing.T) {
	// Arrange
	ctx := context.Background()
	meiliContainer, meiliAddr := setupMeiliContainer(ctx, t)
	defer terminateMeiliContainer(t, meiliContainer)
	host, port := parseMeiliAddress(t, meiliAddr)
	indexName := "events"
	meili := setupMeiliClient(host, port, indexName)
	activateExperimentalFeatures(t, meili.Client)
	manager := setupMeiliIndex(t, meili.Client, indexName)
	mockGRPC := setupMockGRPC()
	configureIndexAttributes(t, manager)
	taskID := addTestDocuments(t, manager, testEvents)
	require.NoError(t, waitForMeiliIndexing(manager, taskID, 20, 500*time.Millisecond),
		"timeout")

	handler := Handler{
		db:     nil,
		meili:  meili,
		grpc:   mockGRPC,
		config: &Config{},
	}
	handler.config.SetDefaultQueryLimit(0)
	handler.config.SetMaxQueryLimit(10)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := handler.HandleReq(&tc.filter, "pubkey")
			require.NoError(t, err, "HandleReq failed for test: %s", tc.name)

			// Assert
			var resultIDs []string
			for _, e := range result {
				resultIDs = append(resultIDs, e.ID)
			}
			require.Equal(t, tc.expected, resultIDs, "Unexpected search results")
		})
	}
}

func setupMeiliContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	t.Helper()
	meiliContainer, err := meilisearch.Run(
		ctx,
		"getmeili/meilisearch",
		meilisearch.WithMasterKey(meiliAPIKey),
	)
	require.NoError(t, err, "failed to start Meilisearch container")

	meiliAddr, err := meiliContainer.Address(ctx)
	require.NoError(t, err, "failed to get Meilisearch address")

	return meiliContainer, meiliAddr
}

func terminateMeiliContainer(t *testing.T, meiliContainer testcontainers.Container) {
	t.Helper()

	if err := testcontainers.TerminateContainer(meiliContainer); err != nil {
		t.Logf("failed to terminate container: %s", err)
	}
}

func parseMeiliAddress(t *testing.T, meiliAddr string) (string, uint16) {
	t.Helper()

	parts := strings.Split(meiliAddr, ":")
	portStr := parts[len(parts)-1]
	num, err := strconv.ParseUint(portStr, 10, 16)
	require.NoError(t, err, "failed to parse port")
	port := uint16(num)
	host := meiliAddr[0 : len(meiliAddr)-len(portStr)-1]

	return host, port
}

func setupMeiliClient(host string, port uint16, indexName string) *infra.Meili {
	meili := infra.New(infra.Config{
		Host:              host,
		Port:              port,
		Timeout:           5000,
		DefaultCollection: indexName,
		APIKey:            meiliAPIKey,
	})

	return meili
}

func activateExperimentalFeatures(t *testing.T, client meilisearchGo.ServiceManager) {
	t.Helper()

	_, err := client.ExperimentalFeatures().SetContainsFilter(true).Update()
	require.NoError(t, err, "Failed to activate experimental features")
}

func setupMeiliIndex(t *testing.T, client meilisearchGo.ServiceManager, indexName string) meilisearchGo.IndexManager {
	t.Helper()

	manager := client.Index(indexName)
	require.NoError(t, nil, "failed to create index in Meilisearch")

	return manager
}

func setupMockGRPC() *MockGRPC {
	mockGRPC := new(MockGRPC)
	mockGRPC.On("RegisterService", mock.Anything, mock.Anything, mock.Anything).Return(&grpcclient.RegisterServiceResponse{}, nil)
	mockGRPC.On("GetParameters", mock.Anything).Return(&grpcclient.GetParametersResponse{}, nil)
	mockGRPC.On("AddLog", mock.Anything, mock.Anything, mock.Anything).Return(&grpcclient.AddLogResponse{}, nil)

	return mockGRPC
}

func configureIndexAttributes(t *testing.T, manager meilisearchGo.IndexManager) {
	t.Helper()

	_, err := manager.UpdateSortableAttributes(&[]string{"created_at", "id"})
	require.NoError(t, err, "failed to specify sortable attributes")

	_, err = manager.UpdateSearchableAttributes(&[]string{"content"})
	require.NoError(t, err, "failed to specify searchable attributes")

	_, err = manager.UpdateFilterableAttributes(&[]string{"id", "pubkey", "created_at", "kind", "tags"})
	require.NoError(t, err, "failed to specify filterable attributes")
}

func addTestDocuments(t *testing.T, manager meilisearchGo.IndexManager, events []event.Event) int64 {
	t.Helper()

	addDocsResponse, err := manager.AddDocuments(events, "id")
	require.NoError(t, err, "failed to insert events")

	return addDocsResponse.TaskUID
}

func waitForMeiliIndexing(manager meilisearchGo.IndexManager, taskID int64, maxRetries int, interval time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		task, err := manager.GetTask(taskID)
		if err != nil {
			return err
		}
		if task.Status == "succeeded" {
			return nil
		}
		time.Sleep(interval)
	}

	return fmt.Errorf("adding documents to meilisearch timed out after %v", time.Duration(maxRetries)*interval)
}
