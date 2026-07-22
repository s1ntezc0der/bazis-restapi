package usecase

import (
	"github.com/s1ntezc0der/bazis-restapi/internal/services/auth/repository"
	authRepo "github.com/s1ntezc0der/bazis-restapi/internal/services/auth/repository"
	"github.com/s1ntezc0der/bazis-restapi/internal/services/teams/entity"
	"github.com/s1ntezc0der/bazis-restapi/internal/services/teams/repository"
	"github.com/s1ntezc0der/bazis-restapi/pkg/errors"
)

type TeamService interface {
	CreateTeam(userID int64, req *entity.CreateTeamRequest) (*entity.Team, error)
	GetUserTeams(userID int64) ([]entity.Team, error)
	InviteUser(teamID, inviterID, userID int64) error
}

type teamService struct {
	teamRepo repository.TeamRepository
	authRepo authRepo.AuthRepository
}

func NewTeamService(teamRepo repository.TeamRepository, authRepo authRepo.AuthRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
		authRepo: authRepo,
	}
}

func (s *teamService) CreateTeam(userID int64, req *entity.CreateTeamRequest) (*entity.Team, error) {
	team := &entity.Team{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	}

	if err := s.teamRepo.Create(team); err != nil {
		return nil, err
	}

	member := &entity.TeamMember{
		UserID: userID,
		TeamID: team.ID,
		Role:   "owner",
	}
	if err := s.teamRepo.AddMember(member); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *teamService) GetUserTeams(userID int64) ([]entity.Team, error) {
	return s.teamRepo.GetByUserID(userID)
}

func (s *teamService) InviteUser(teamID, inviterID, userID int64) error {
	member, err := s.teamRepo.GetMember(teamID, inviterID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.ErrNotMember
	}
	if member.Role != "owner" && member.Role != "admin" {
		return errors.ErrNotAdmin
	}

	_, err = s.authRepo.GetUserByID(userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	existing, _ := s.teamRepo.GetMember(teamID, userID)
	if existing != nil {
		return errors.ErrConflict
	}

	newMember := &entity.TeamMember{
		UserID: userID,
		TeamID: teamID,
		Role:   "member",
	}

	// Мок email service с Circuit Breaker
    emailService := middleware.NewCircuitBreaker(3, 2, 5*time.Second)
    err := emailService.Call(func() error {
        // Имитация отправки email
        log.Printf("📧 Приглашение отправлено пользователю %d в команду %d", userID, teamID)
        return nil
    })
    if err != nil {
        return errors.ErrInternal
    }
    
    // return nil
	
	return s.teamRepo.AddMember(newMember)
}

