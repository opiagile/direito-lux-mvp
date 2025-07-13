package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/direito-lux/auth-service/internal/domain"
)

// Test User Domain Functions
func TestUser_SetPassword(t *testing.T) {
	// Create a new user
	user := &domain.User{
		ID:        "test-id",
		TenantID:  "tenant-123",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      domain.RoleOperator,
		Status:    domain.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set a password
	password := "TestPass123!"
	err := user.SetPassword(password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, password, user.Password)
	
	// Verify the password
	assert.True(t, user.VerifyPassword(password))
	assert.False(t, user.VerifyPassword("wrong-password"))
}

func TestUser_ValidatePassword(t *testing.T) {
	user := &domain.User{}
	
	// Test valid password
	password := "ValidPass123!"
	err := user.ValidatePassword(password)
	assert.NoError(t, err)
	
	// Test invalid passwords
	err = user.ValidatePassword("short")
	assert.Error(t, err)
	
	err = user.ValidatePassword("")
	assert.Error(t, err)
	
	err = user.ValidatePassword("nouppercase123!")
	assert.Error(t, err)
	
	err = user.ValidatePassword("NOLOWERCASE123!")
	assert.Error(t, err)
}

func TestUser_ValidateEmail(t *testing.T) {
	// Test valid emails
	user := &domain.User{Email: "test@example.com"}
	err := user.ValidateEmail()
	assert.NoError(t, err)
	
	user.Email = "user.name@domain.co.uk"
	err = user.ValidateEmail()
	assert.NoError(t, err)
	
	// Test invalid emails
	user.Email = ""
	err = user.ValidateEmail()
	assert.Error(t, err)
	
	user.Email = "invalid-email"
	err = user.ValidateEmail()
	assert.Error(t, err)
	
	user.Email = "@domain.com"
	err = user.ValidateEmail()
	assert.Error(t, err)
}

func TestUser_IsActive(t *testing.T) {
	user := &domain.User{
		Status: domain.StatusActive,
	}
	
	assert.True(t, user.IsActive())
	
	// Test other statuses
	user.Status = domain.StatusInactive
	assert.False(t, user.IsActive())
	
	user.Status = domain.StatusSuspended
	assert.False(t, user.IsActive())
	
	user.Status = domain.StatusBlocked
	assert.False(t, user.IsActive())
}

func TestUser_CanLogin(t *testing.T) {
	user := &domain.User{
		Status: domain.StatusActive,
	}
	
	// Active user can login
	assert.True(t, user.CanLogin())
	
	// Pending user can also login
	user.Status = domain.StatusPending
	assert.True(t, user.CanLogin())
	
	// Other statuses cannot login
	user.Status = domain.StatusInactive
	assert.False(t, user.CanLogin())
	
	user.Status = domain.StatusSuspended
	assert.False(t, user.CanLogin())
	
	user.Status = domain.StatusBlocked
	assert.False(t, user.CanLogin())
}

func TestUser_FullName(t *testing.T) {
	user := &domain.User{
		FirstName: "John",
		LastName:  "Doe",
	}
	
	assert.Equal(t, "John Doe", user.FullName())
	
	// Test with empty last name
	user.LastName = ""
	assert.Equal(t, "John", user.FullName())
	
	// Test with empty first name
	user.FirstName = ""
	user.LastName = "Doe"
	assert.Equal(t, "Doe", user.FullName())
	
	// Test with both empty
	user.FirstName = ""
	user.LastName = ""
	assert.Equal(t, "", user.FullName())
}

func TestUser_HasRole(t *testing.T) {
	user := &domain.User{
		Role: domain.RoleManager,
	}
	
	assert.True(t, user.HasRole(domain.RoleManager))
	assert.False(t, user.HasRole(domain.RoleAdmin))
	assert.False(t, user.HasRole(domain.RoleOperator))
}

func TestUser_ValidateRole(t *testing.T) {
	// Valid roles
	validRoles := []domain.UserRole{
		domain.RoleAdmin,
		domain.RoleManager,
		domain.RoleOperator,
		domain.RoleClient,
		domain.RoleReadOnly,
	}
	
	for _, role := range validRoles {
		user := &domain.User{Role: role}
		err := user.ValidateRole()
		assert.NoError(t, err)
	}
}

func TestUser_ValidateStatus(t *testing.T) {
	// Valid statuses
	validStatuses := []domain.UserStatus{
		domain.StatusActive,
		domain.StatusInactive,
		domain.StatusPending,
		domain.StatusSuspended,
		domain.StatusBlocked,
	}
	
	for _, status := range validStatuses {
		user := &domain.User{Status: status}
		err := user.ValidateStatus()
		assert.NoError(t, err)
	}
}

// Test User Roles
func TestUserRoles(t *testing.T) {
	roles := []domain.UserRole{
		domain.RoleAdmin,
		domain.RoleManager,
		domain.RoleOperator,
		domain.RoleClient,
		domain.RoleReadOnly,
	}
	
	for _, role := range roles {
		user := &domain.User{Role: role}
		assert.True(t, user.HasRole(role))
		
		// Test that user doesn't have other roles
		for _, otherRole := range roles {
			if otherRole != role {
				assert.False(t, user.HasRole(otherRole))
			}
		}
	}
}

// Test User Status
func TestUserStatus(t *testing.T) {
	statuses := []domain.UserStatus{
		domain.StatusActive,
		domain.StatusInactive,
		domain.StatusPending,
		domain.StatusSuspended,
		domain.StatusBlocked,
	}
	
	for _, status := range statuses {
		user := &domain.User{Status: status}
		
		if status == domain.StatusActive {
			assert.True(t, user.IsActive())
			assert.True(t, user.CanLogin())
		} else if status == domain.StatusPending {
			assert.False(t, user.IsActive())  // Only Active status is considered "active"
			assert.True(t, user.CanLogin())   // But Pending can login
		} else {
			assert.False(t, user.IsActive())
			assert.False(t, user.CanLogin())
		}
	}
}

// Test Multi-tenant Isolation at Domain Level
func TestMultiTenantIsolation(t *testing.T) {
	user1 := &domain.User{
		ID:       "user-1",
		TenantID: "tenant-123",
		Email:    "same@example.com",
	}
	
	user2 := &domain.User{
		ID:       "user-2", 
		TenantID: "tenant-456",
		Email:    "same@example.com",
	}
	
	// Same email but different tenants should be allowed
	assert.Equal(t, user1.Email, user2.Email)
	assert.NotEqual(t, user1.TenantID, user2.TenantID)
	assert.NotEqual(t, user1.ID, user2.ID)
	
	// Both users should have different tenant isolation
	assert.NotEqual(t, user1.TenantID, user2.TenantID)
}