package mail

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	grpc "notification-service/infrastructure/gRPC"
	"notification-service/internal/model"
	"time"
)

type MailParams struct {
	Email    string
	Password string
	GRPC     *grpc.GrpcComm
}
type Mail struct {
	Email    string
	Password string
	grpc     *grpc.GrpcComm
}

func Init(params MailParams) (*Mail, error) {
	var mail *Mail = &Mail{
		Email:    params.Email,
		Password: params.Password,
		grpc:     params.GRPC,
	}
	return mail, nil
}

func (m *Mail) SendTaskMsgEmail(actionTask, body string) error {
	var response *model.Response = &model.Response{}

	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		return fmt.Errorf("failed convert message task from broker. err = %v", err)
	}

	var task *model.Task = &model.Task{}
	if response.Data != nil {
		dataBytes, err := json.Marshal(response.Data)
		if err != nil {
			return fmt.Errorf("failed convert message data of task to json. err = %v", err)
		}

		// Unmarshal the data into the Task struct
		err = json.Unmarshal(dataBytes, &task)
		if err != nil {
			return fmt.Errorf("failed convert message data of task to struct. err = %v", err)
		}
	}

	fmt.Println("task = ", task)
	return m.processSendTaskEmail(task, actionTask)
}

func (m *Mail) processSendTaskEmail(task *model.Task, actionTask string) error {
	responseGrpc, err := m.grpc.UserGrpcClient.RequestUserInfo(task.AssignedTo, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed request email to user service. err = %s", err.Error())
	}

	subject := m.createSubjectMailTask(actionTask, task.Title)
	bodyMail := m.createBodyMailTask(actionTask, task)

	if err = m.sendMessageGmail(responseGrpc.Email, subject, bodyMail); err != nil {
		return fmt.Errorf("failed send notification task to email. err = %v", err)
	}

	return nil
}

func (m *Mail) createSubjectMailTask(actionTask, taskTitle string) string {
	var subject string = ""
	switch actionTask {
	case model.ACTION_TASK_CREATE:
		subject = fmt.Sprintf("[task-service] Create Task '%s'", taskTitle)
	case model.ACTION_TASK_UPDATE:
		subject = fmt.Sprintf("[task-service] Update Task '%s'", taskTitle)
	case model.ACTION_TASK_READ:
		subject = fmt.Sprintf("[task-service] Read Task '%s'", taskTitle)
	case model.ACTION_TASK_DELETE:
		subject = fmt.Sprintf("[task-service] Delete Task '%s'", taskTitle)
	default:
		subject = "[task-service] Task Notification"
	}

	return subject
}

func (m *Mail) createBodyMailTask(actionTask string, task *model.Task) string {
	var body string = ""
	switch actionTask {
	case model.ACTION_TASK_CREATE:
		body = fmt.Sprintf(`You have assigned a task '%s'.
Information about task :
	id          = %d
	title       = %s
	description = %s
	status      = %s
	due date    = %s
Please dont forget update your task 🤞.`, task.Title, task.Id, task.Title, task.Description, task.Status, task.DueDate.Format("2006-01-02 15:04:05"))
	case model.ACTION_TASK_UPDATE:
		body = fmt.Sprintf(`You have updated a task task '%s'.
Information about task :
	id          = %d
	title       = %s
	description = %s
	status      = %s
	due date    = %s
Please dont forget update your task 🤞.`, task.Title, task.Id, task.Title, task.Description, task.Status, task.DueDate.Format("2006-01-02 15:04:05"))
	case model.ACTION_TASK_DELETE:
		body = fmt.Sprintf(`You have delete a task task '%s'. Please dont forget update your other tasks 🤞.`, task.Title)

	default:
		body = "[task-service] Task Notification"
	}

	return body
}

func (m *Mail) sendMessageGmail(emailPurpose, subject, body string) error {
	var outputEmail []string = []string{emailPurpose}
	// Set up authentication information
	auth := smtp.PlainAuth("", m.Email, m.Password, "smtp.gmail.com")
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	// Mengirim email via Gmail SMTP
	return smtp.SendMail("smtp.gmail.com:587", auth, m.Email, outputEmail, msg)
}
