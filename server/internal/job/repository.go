package job

import (
	"fmt"
	"log"
	"verve-hrms/internal/schema"

	"gorm.io/gorm"
)

type Repository interface {
}

type JobRepository struct {
	client *gorm.DB
}

func NewJobRepository(client *gorm.DB) *JobRepository {
	return &JobRepository{client: client}
}

//! Job     ------------------------------------------------------

func (jr JobRepository) JobCreate(newJob *schema.Job) (*schema.Job, error) {
	log.Printf("job.r.job_create: newJob: %v", newJob)

	result := jr.client.Create(newJob)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_create: %w", result.Error)
	}

	return newJob, nil
}

func (jr JobRepository) JobRead(jobID uint) (*schema.Job, error) {
	var job schema.Job

	result := jr.client.First(&job, jobID)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read: %w", result.Error)
	}

	return &job, nil
}

func (jr JobRepository) JobReadAndExpand(jobID uint) (*schema.Job, error) {
	var job schema.Job

	result := jr.client.Preload("Subordinates").Preload("AssignedJobs").First(&job, jobID)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_and_expand: %w", result.Error)
	}

	return &job, nil
}

func (jr JobRepository) JobReadByCompany(companyID uint) ([]*schema.Job, error) {
	var jobs []*schema.Job

	result := jr.client.Where("company_id = ?", companyID).Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_by_company: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_by_company: %w", gorm.ErrRecordNotFound)
	}

	return jobs, nil
}

func (jr JobRepository) JobReadByCompanyAndExpand(companyID uint) ([]*schema.Job, error) {
	var jobs []*schema.Job

	result := jr.client.Preload("Subordinates").Preload("AssignedJobs").Where("company_id = ?", companyID).Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_by_company_and_expand: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_by_company_and_expand: %w", ErrNoRowsFound)
	}

	return jobs, nil
}

func (jr JobRepository) JobReadAll() ([]*schema.Job, error) {
	var jobs []*schema.Job

	result := jr.client.Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_all: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_all: %w", ErrEmptyTable)
	}

	return jobs, nil
}

func (jr JobRepository) JobReadAndExpandAll() ([]*schema.Job, error) {
	var jobs []*schema.Job

	result := jr.client.Preload("Subordinates").Preload("AssignedJobs").Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_and_expand_all: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_and_expand_all: %w", ErrEmptyTable)
	}

	return jobs, nil
}

func (jr JobRepository) JobUpdate(jobID uint, updateData *schema.Job) (*schema.Job, error) {
	var job schema.Job

	updateDataMap := map[string]interface{}{
		"Title":          updateData.Title,
		"Description":    updateData.Description,
		"Duties":         updateData.Duties,
		"Qualifications": updateData.Qualifications,
		"Experience":     updateData.Experience,
		"MinSalary":      updateData.MinSalary,
		"MaxSalary":      updateData.MaxSalary,

		// "CompanyID": updateData.CompanyID, //* not part of updateData, will always initialize to 0 since it's not a pointer
		"LocationID":   updateData.LocationID,
		"DepartmentID": updateData.DepartmentID,
		"ManagerID":    updateData.ManagerID,
	}

	//*Updates only updates non-zero values, need to select nil values explicitly
	result := jr.client.Model(&job).Where("ID = ?", jobID).Updates(updateDataMap)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_update: %w", result.Error)
	}

	// Check if any row was affected, if not, the job does not exist.
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("job.r.job_update: %w", gorm.ErrRecordNotFound)
	}

	return &job, nil
}

func (jr JobRepository) JobDelete(jobID uint) error {
	var job *schema.Job

	result := jr.client.Unscoped().Delete(&job, jobID)
	if result.Error != nil {
		return fmt.Errorf("job.r.job_delete: %w", result.Error)
	}

	return nil
}

//! AssignedJob ------------------------------------------------------

func (jr JobRepository) AssignedJobCreate(newAssignedJob *schema.AssignedJob) (*schema.AssignedJob, error) {
	result := jr.client.Create(newAssignedJob)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_create: %w", result.Error)
	}

	return newAssignedJob, nil
}

func (jr JobRepository) AssignedJobRead(assignedJobID uint) (*schema.AssignedJob, error) {
	var assignedJob schema.AssignedJob

	result := jr.client.First(&assignedJob, assignedJobID)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_read: %w", result.Error)
	}

	return &assignedJob, nil
}

func (jr JobRepository) AssignedJobReadAll() ([]*schema.AssignedJob, error) {
	var assignedJobs []*schema.AssignedJob

	result := jr.client.Find(&assignedJobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_read_all: %w", result.Error)
	}
	if len(assignedJobs) == 0 {
		return nil, fmt.Errorf("job.r.assigned_job_read_all: %w", ErrEmptyTable)
	}

	return assignedJobs, nil
}

func (jr JobRepository) AssignedJobUpdate(assignedJobID uint, updateData *schema.AssignedJob) error {
	var assignedJob *schema.AssignedJob

	result := jr.client.Model(&assignedJob).Where("ID = ?", assignedJobID).Updates(updateData)
	if result.Error != nil {
		return fmt.Errorf("job.r.assigned_job_update: %w", result.Error)
	}

	// Check if any row was affected, if not, the assignedJob does not exist.
	if result.RowsAffected == 0 {
		return fmt.Errorf("job.r.assigned_job_update: %w", gorm.ErrRecordNotFound)
	}

	return nil
}

func (jr JobRepository) AssignedJobDelete(AssignedJobID uint) error {
	var assignedJob *schema.AssignedJob

	result := jr.client.Unscoped().Delete(&assignedJob, AssignedJobID)
	if result.Error != nil {
		return fmt.Errorf("job.r.assigned_job_delete: %w", result.Error)
	}

	return nil
}
