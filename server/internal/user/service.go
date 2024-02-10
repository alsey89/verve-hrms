package user

import (
	"fmt"
	"verve-hrms/internal/schema"
)

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

//! Auth ------------------------------------------------------------

func (us *UserService) CreateNewUser(newUser *schema.User) (*schema.User, error) {
	var createdUser *schema.User
	var err error

	createdUser, err = us.userRepository.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_account: %w", err)
	}

	return createdUser, nil
}

func (us *UserService) GetUserByEmail(email string) (*schema.User, error) {
	var existingUser *schema.User
	var err error

	existingUser, err = us.userRepository.ReadByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user.s.get_user_by_email: %w", err)
	}

	return existingUser, nil
}

//! User -----------------------------------------------------------

func (us *UserService) GetAllUsersAndExpand(companyID uint) ([]*schema.User, error) {

	existingUsers, err := us.userRepository.ReadAllAndExpand()
	if err != nil {
		return nil, fmt.Errorf("user.s.get_all_users: %w", err)
	}

	filteredUsers := make([]*schema.User, 0)
	for _, user := range existingUsers {
		// Check for nil pointers before accessing properties
		if user.AssignedJob != nil && user.AssignedJob.Job.CompanyID == companyID {
			filteredUsers = append(filteredUsers, user)
		}
	}

	return filteredUsers, nil // Return the filtered users
}

func (us *UserService) GetUserByIDAndExpand(ID uint) (*schema.User, error) {
	var existingUser *schema.User
	var err error

	existingUser, err = us.userRepository.ReadAndExpand(ID)
	if err != nil {
		return nil, fmt.Errorf("user.s.get_user_by_id: %w", err)
	}

	return existingUser, nil
}

func (us *UserService) CreateNewUserAndGetAllUsersAndExpand(companyID uint, newUser *schema.User) ([]*schema.User, error) {
	_, err := us.userRepository.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_user: %w", err)
	}

	expandedUserList, err := us.GetAllUsersAndExpand(companyID)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_user: %w", err)
	}

	return expandedUserList, nil
}

func (us *UserService) UpdateUserAndGetAllUsersAndExpand(companyID uint, userID uint, updateData schema.User) ([]*schema.User, error) {
	_, err := us.userRepository.Update(userID, updateData)
	if err != nil {
		return nil, fmt.Errorf("user.s.update_user: %w", err)
	}

	expandedUserList, err := us.GetAllUsersAndExpand(companyID)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_user: %w", err)
	}

	return expandedUserList, nil
}

func (us *UserService) DeleteUserAndGetAllUsersAndExpand(companyID uint, userID uint) ([]*schema.User, error) {
	err := us.userRepository.Delete(userID)
	if err != nil {
		return nil, fmt.Errorf("s.delete_user: %w", err)
	}

	expandedUserList, err := us.GetAllUsersAndExpand(companyID)
	if err != nil {
		return nil, fmt.Errorf("user.s.create_new_user: %w", err)
	}

	return expandedUserList, nil
}
