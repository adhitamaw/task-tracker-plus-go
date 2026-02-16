package client

import (
	"a21hc3NpZ25tZW50/config"
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UserClient interface {
	Login(email, password string) (respCode int, err error)
	Register(fullname, email, password string) (respCode int, err error)
	GetUserByEmail(email string, token string) (model.User, error)
	GetUserTaskCategory(token string) (*[]model.UserTaskCategory, error)
}

type userClient struct {
}

func NewUserClient() *userClient {
	return &userClient{}
}

func (u *userClient) Login(email, password string) (respCode int, err error) {
	datajson := map[string]string{
		"email":    email,
		"password": password,
	}

	data, err := json.Marshal(datajson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/user/login"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	// Baca response body untuk mendapatkan pesan error
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// Jika status bukan 200 OK, kembalikan error dengan pesan dari response
	if resp.StatusCode != http.StatusOK {
		var errResp model.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
			return resp.StatusCode, errors.New(errResp.Error)
		}
		return resp.StatusCode, errors.New("Login failed with status: " + strconv.Itoa(resp.StatusCode))
	}

	return resp.StatusCode, nil
}

func (u *userClient) Register(fullname, email, password string) (respCode int, err error) {
	datajson := map[string]string{
		"fullname": fullname,
		"email":    email,
		"password": password,
	}

	data, err := json.Marshal(datajson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/user/register"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	// Baca response body untuk mendapatkan pesan error
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// Jika status bukan 201 Created, kembalikan error dengan pesan dari response
	if resp.StatusCode != http.StatusCreated {
		var errResp model.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
			return resp.StatusCode, errors.New(errResp.Error)
		}
		return resp.StatusCode, errors.New("Unknown server error")
	}

	return resp.StatusCode, nil
}

func (u *userClient) GetUserTaskCategory(token string) (*[]model.UserTaskCategory, error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/v1/user/tasks"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("status code not 200")
	}

	var userTasks []model.UserTaskCategory
	err = json.Unmarshal(b, &userTasks)
	if err != nil {
		return nil, err
	}

	return &userTasks, nil
}

func (u *userClient) GetUserByEmail(email string, token string) (model.User, error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return model.User{}, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/v1/user/profile/"+email), nil)
	if err != nil {
		return model.User{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return model.User{}, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.User{}, err
	}

	if resp.StatusCode != 200 {
		return model.User{}, errors.New("failed to get user info: " + strconv.Itoa(resp.StatusCode))
	}

	var user model.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
