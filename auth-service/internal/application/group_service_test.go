package application

import (
	"testing"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"gorm.io/gorm"
)

type mockGroupRepo struct {
	groups map[string]*entity.Group
}

func newMockGroupRepo() *mockGroupRepo {
	return &mockGroupRepo{groups: make(map[string]*entity.Group)}
}

func (m *mockGroupRepo) Create(group *entity.Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepo) GetByID(id string) (*entity.Group, error) {
	g, ok := m.groups[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return g, nil
}

func (m *mockGroupRepo) Update(group *entity.Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepo) Delete(id string) error {
	delete(m.groups, id)
	return nil
}

func (m *mockGroupRepo) List(page, pageSize int) ([]*entity.Group, int, error) {
	var all []*entity.Group
	for _, g := range m.groups {
		all = append(all, g)
	}
	total := len(all)
	offset := (page - 1) * pageSize
	if offset >= total {
		return nil, total, nil
	}
	end := offset + pageSize
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func TestGroupService_CreateGroup(t *testing.T) {
	repo := newMockGroupRepo()
	svc := NewGroupService(repo)

	group, err := svc.CreateGroup("developers", "Dev team", "", nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateGroup() error = %v", err)
	}
	if group.Name != "developers" {
		t.Errorf("Expected name 'developers', got %s", group.Name)
	}
	if group.ID == "" {
		t.Error("Expected non-empty ID")
	}
}

func TestGroupService_CreateGroupWithConfig(t *testing.T) {
	repo := newMockGroupRepo()
	svc := NewGroupService(repo)

	group, err := svc.CreateGroup("power-users", "", "",
		[]string{"gpt-4", "claude-*"},
		&entity.TokenLimit{PromptTokens: 100000, CompletionTokens: 100000, Period: "daily"},
		&entity.RateLimit{RequestsPerMinute: 60, RequestsPerDay: 10000},
	)
	if err != nil {
		t.Fatalf("CreateGroup() error = %v", err)
	}
	if len(group.ModelPatterns) != 2 {
		t.Errorf("Expected 2 model patterns, got %d", len(group.ModelPatterns))
	}
	if group.TokenLimit.PromptTokens != 100000 {
		t.Errorf("Expected prompt_tokens 100000, got %d", group.TokenLimit.PromptTokens)
	}
}

func TestGroupService_UpdateGroup(t *testing.T) {
	repo := newMockGroupRepo()
	svc := NewGroupService(repo)

	group, _ := svc.CreateGroup("developers", "", "", nil, nil, nil)

	updated, err := svc.UpdateGroup(group.ID, "senior-devs", "")
	if err != nil {
		t.Fatalf("UpdateGroup() error = %v", err)
	}
	if updated.Name != "senior-devs" {
		t.Errorf("Expected name 'senior-devs', got %s", updated.Name)
	}
}

func TestGroupService_DeleteGroup(t *testing.T) {
	repo := newMockGroupRepo()
	svc := NewGroupService(repo)

	group, _ := svc.CreateGroup("to-delete", "", "", nil, nil, nil)

	err := svc.DeleteGroup(group.ID)
	if err != nil {
		t.Fatalf("DeleteGroup() error = %v", err)
	}

	_, err = svc.groupRepo.GetByID(group.ID)
	if err == nil {
		t.Error("Expected error after deletion")
	}
}

func TestGroupService_ListGroups(t *testing.T) {
	repo := newMockGroupRepo()
	svc := NewGroupService(repo)

	svc.CreateGroup("group-a", "", "", nil, nil, nil)
	svc.CreateGroup("group-b", "", "", nil, nil, nil)

	groups, total, err := svc.ListGroups(1, 10)
	if err != nil {
		t.Fatalf("ListGroups() error = %v", err)
	}
	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
}

// Integration test lives in repository package to access GORM implementations
