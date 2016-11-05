package common

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// JobLogger is used to log messages from within the context of a scheduler job.
type JobLogger interface {
	// Info logs the args at the info log level.
	Info(args ...interface{})
	// Infof logs the template at the info log level.
	Infof(template string, args ...interface{})
	// Error logs the args at the error log level.
	Error(args ...interface{})
	// Errorf logs the template at the error log level.
	Errorf(template string, args ...interface{})

	// Start logs that the job started.
	Start()
	// Finish logs that the job finished.
	Finish()
}

// NewJobLogger external proxy for newJobLoggerImpl that reads parameters from
// an HTTP request.
func NewJobLogger(jobName string, jobParams JobParams) JobLogger {
	return newJobLoggerImpl(jobParams.ID, jobName, jobParams.StartTime)
}

// jobLoggerImpl is the internal implementation of the job logger.
type jobLoggerImpl struct {
	logPrefix string
	startTime time.Time
}

// newJobLoggerImpl creates a new jobLoggerImpl.
func newJobLoggerImpl(
	jobID string,
	jobName string,
	jobStartTime time.Time,
) *jobLoggerImpl {
	if strings.IndexByte(jobID, '%') != -1 {
		panic("jobID cannot have '%'")
	} else if strings.IndexByte(jobName, '%') != -1 {
		panic("jobName cannot have '%'")
	}

	return &jobLoggerImpl{
		logPrefix: fmt.Sprintf(
			`[scheduler:job:%s (id: %s, started: %s)] `,
			jobName,
			jobID,
			jobStartTime.Format(time.RFC3339)),
	}
}

// Info logs the args at the info log level.
func (j *jobLoggerImpl) Info(args ...interface{}) {
	if len(args) > 0 {
		args[0] = j.logPrefix + "INFO " + fmt.Sprint(args[0])
	} else {
		args = append(args, j.logPrefix+"INFO ")
	}

	log.Println(args...)
}

// Infof logs the template at the info log level.
func (j *jobLoggerImpl) Infof(template string, args ...interface{}) {
	template = j.logPrefix + "INFO " + template
	log.Printf(template, args...)
}

// Error logs the args at the error log level.
func (j *jobLoggerImpl) Error(args ...interface{}) {
	if len(args) > 0 {
		args[0] = j.logPrefix + "ERR  " + fmt.Sprint(args[0])
	} else {
		args = append(args, j.logPrefix+"ERR  ")
	}

	log.Println(args...)
}

// Errorf logs the template at the error log level.
func (j *jobLoggerImpl) Errorf(template string, args ...interface{}) {
	template = j.logPrefix + "ERR  " + template
	log.Printf(template, args...)
}

// Start logs that the job started.
func (j *jobLoggerImpl) Start() {
	j.startTime = time.Now()
	j.Info("Job execution started.")
}

// Start logs that the job started.
func (j *jobLoggerImpl) Finish() {
	j.Infof("Job execution finished in %s.\n", time.Since(j.startTime).String())
}
