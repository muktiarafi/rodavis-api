package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/muktiarafi/rodavis-api/internal/entity"
	"github.com/muktiarafi/rodavis-api/internal/model"
)

func TestReportHandlerNewReport(t *testing.T) {
	t.Run("create new report normally", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bambankkk",
			PhoneNumber: "+6217344678910",
			Email:       "bambankk@gmail.com",
			Password:    "12345678",
		}
		userDTO, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		createReport := map[string]string{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "mataram",
		}

		for k, v := range createReport {
			writer.WriteField(k, v)
		}

		fileName := "jalan.jpg"
		file, err := os.Open(filepath.Join(imagePath, fileName))
		if err != nil {
			t.Error(err)
		}
		part, err := writer.CreateFormFile("image", fileName)
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Error(err)
		}

		writer.Close()
		req := httptest.NewRequest(http.MethodPost, "/api/reports", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+userDTO.Token)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusCreated, res.Code)

		resBody, _ := ioutil.ReadAll(res.Body)
		apiResponse := struct {
			Data *entity.Report `json:"data"`
		}{}
		json.Unmarshal(resBody, &apiResponse)

		if apiResponse.Data.ReporterName != createUserDTO.Name {
			t.Errorf("Expecting reporter name to be %q, but got %q instead", apiResponse.Data.ReporterName, createUserDTO.Name)
		}
		if len(apiResponse.Data.Classes) == 0 {
			t.Error("Expecting classes to be not empty")
		}

		if !strings.Contains(apiResponse.Data.ImageURL, "https://storage.googleapis.com") {
			t.Error("Expecting image url to contain https://storage.googleapis.com")
		}
	})

	t.Run("create new report without lat lng", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bambankkk",
			PhoneNumber: "+62173445678910",
			Email:       "ewrrew@gmail.com",
			Password:    "12345678",
		}
		userDTO, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		createReport := map[string]string{
			"note":    "",
			"address": "mataram",
		}

		for k, v := range createReport {
			writer.WriteField(k, v)
		}

		fileName := "jalan.jpg"
		file, err := os.Open(filepath.Join(imagePath, fileName))
		if err != nil {
			t.Error(err)
		}
		part, err := writer.CreateFormFile("image", fileName)
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Error(err)
		}

		writer.Close()
		req := httptest.NewRequest(http.MethodPost, "/api/reports", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+userDTO.Token)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("create new report without image", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bambankkk",
			PhoneNumber: "+6217768678910",
			Email:       "ytytryrtttt@gmail.com",
			Password:    "12345678",
		}
		userDTO, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		createReport := map[string]string{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "mataram",
		}

		for k, v := range createReport {
			writer.WriteField(k, v)
		}

		writer.Close()
		req := httptest.NewRequest(http.MethodPost, "/api/reports", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+userDTO.Token)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("create new report without body", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bambankkk",
			PhoneNumber: "+6217768998910",
			Email:       "lkfgdkret@gmail.com",
			Password:    "12345678",
		}
		userDTO, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		writer.Close()
		req := httptest.NewRequest(http.MethodPost, "/api/reports", nil)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+userDTO.Token)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})
}

func TestReportHandlerGetReports(t *testing.T) {
	createUserDTO := &model.CreateUserDTO{
		Name:        "bambankkk",
		PhoneNumber: "+6217344644410",
		Email:       "rewkkkjfer@gmail.com",
		Password:    "12345678",
	}
	userDTO, res := register(createUserDTO)

	assertResponseCode(t, http.StatusCreated, res.Code)

	reportFields := []map[string]string{
		{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "mataram",
		},
		{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "lawang sewu",
		},
		{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "area 51",
		},
		{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "bikini bottom",
		},
		{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "konoha",
		},
	}

	for _, v := range reportFields {
		res := sendReport(t, userDTO.Token, v)
		assertResponseCode(t, http.StatusCreated, res.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/reports", nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assertResponseCode(t, http.StatusOK, res.Code)

	resBody, _ := ioutil.ReadAll(res.Body)
	apiResponse := struct {
		Data []*entity.Report `json:"data"`
	}{}
	json.Unmarshal(resBody, &apiResponse)

	if len(apiResponse.Data) == 0 {
		t.Error("Expecting reports to not empty")
	}

	t.Run("get all reports without pagination", func(t *testing.T) {
		current := apiResponse.Data[0].ID
		for _, v := range apiResponse.Data {
			if current < v.ID {
				t.Error("Expecting reports to be ordered by ID descending")
			}
			current = v.ID
		}
	})

	t.Run("get all reports with limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/reports?limit=3", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		resBody, _ := ioutil.ReadAll(res.Body)
		apiResponse := struct {
			Data []*entity.Report `json:"data"`
		}{}
		json.Unmarshal(resBody, &apiResponse)

		if len(apiResponse.Data) != 3 {
			t.Errorf("Expecting the lenght of reports to be 3 but got %d instead", len(apiResponse.Data))
		}
	})

	t.Run("get all reports with limit and lastseenid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/reports?limit=3&lastseenid=4", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		resBody, _ := ioutil.ReadAll(res.Body)
		apiResponse := struct {
			Data []*entity.Report `json:"data"`
		}{}
		json.Unmarshal(resBody, &apiResponse)

		if len(apiResponse.Data) != 3 {
			t.Errorf("Expecting the lenght of reports to be 3 but got %d instead", len(apiResponse.Data))
		}

		lastSeenID := apiResponse.Data[len(apiResponse.Data)-1].ID
		if lastSeenID != 1 {
			t.Errorf("Expecting the id of reports to be 1 but got %d instead", lastSeenID)
		}
	})
}

func TestReportHandlerGetAllUserReport(t *testing.T) {
	createUserDTO := &model.CreateUserDTO{
		Name:        "paijower",
		PhoneNumber: "+6217347744410",
		Email:       "klkeruer@gmail.com",
		Password:    "12345678",
	}
	userDTO, res := register(createUserDTO)

	assertResponseCode(t, http.StatusCreated, res.Code)

	reportFields := []map[string]string{
		{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "mataram",
		},
		{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "lawang sewu",
		},
		{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "area 51",
		},
	}

	for _, v := range reportFields {
		res := sendReport(t, userDTO.Token, v)
		assertResponseCode(t, http.StatusCreated, res.Code)

	}

	req := httptest.NewRequest(http.MethodGet, "/api/reports/history", nil)
	req.Header.Set("Authorization", "Bearer "+userDTO.Token)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	resBody, _ := ioutil.ReadAll(res.Body)
	apiResponse := struct {
		Data []*entity.Report `json:"data"`
	}{}
	json.Unmarshal(resBody, &apiResponse)

	assertResponseCode(t, http.StatusOK, res.Code)

	if len(apiResponse.Data) != 3 {
		t.Errorf("Expecting reports length to be 3 but got %d instead", len(apiResponse.Data))
	}

	for _, v := range apiResponse.Data {
		if v.ReporterName != createUserDTO.Name {
			t.Errorf("Expecting reporter name to be %q but got %q instead", createUserDTO.Name, v.ReporterName)
		}
	}

	t.Run("get all reports without pagination", func(t *testing.T) {
		current := apiResponse.Data[0].ID
		for _, v := range apiResponse.Data {
			if current < v.ID {
				t.Error("Expecting reports to be ordered by ID descending")
			}
			current = v.ID
		}
	})

	t.Run("get all reports with limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/reports/history?limit=2", nil)
		req.Header.Set("Authorization", "Bearer "+userDTO.Token)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusOK, res.Code)

		resBody, _ := ioutil.ReadAll(res.Body)
		apiResponse := struct {
			Data []*entity.Report `json:"data"`
		}{}
		json.Unmarshal(resBody, &apiResponse)

		if len(apiResponse.Data) != 2 {
			t.Errorf("Expecting the lenght of reports to be 2 but got %d instead", len(apiResponse.Data))
		}
	})

	t.Run("get all reports with limit and lastseenid", func(t *testing.T) {
		lastseenID := 9
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/reports/history?limit=2&lastseenid=%d", lastseenID), nil)
		req.Header.Set("Authorization", "Bearer "+userDTO.Token)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusOK, res.Code)

		resBody, _ := ioutil.ReadAll(res.Body)
		apiResponse := struct {
			Data []*entity.Report `json:"data"`
		}{}
		json.Unmarshal(resBody, &apiResponse)

		if len(apiResponse.Data) != 2 {
			t.Errorf("Expecting the length of reports to be 2 but got %d instead", len(apiResponse.Data))
		}

		for _, v := range apiResponse.Data {
			if v.ID > lastseenID {
				t.Errorf("Expecting report id to be less than %d but got %d instead", lastseenID, v.ID)
			}
		}
	})
}

func sendReport(t *testing.T, token string, reportMap map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for k, v := range reportMap {
		writer.WriteField(k, v)
	}

	fileName := "jalan.jpg"
	file, err := os.Open(filepath.Join(imagePath, fileName))
	if err != nil {
		t.Error(err)
	}
	part, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		t.Error(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Error(err)
	}

	writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/reports", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	return res
}

func TestReportHandlerUpdateReport(t *testing.T) {
	b, _ := json.Marshal(admin)
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	assertResponseCode(t, http.StatusOK, res.Code)

	resBody, _ := ioutil.ReadAll(res.Body)
	adminDTO := struct {
		Data *model.UserDTO `json:"data"`
	}{}
	json.Unmarshal(resBody, &adminDTO)

	if adminDTO.Data.User.Email != admin.Email {
		t.Errorf("Expecting email to be %q but got %q instead", adminDTO.Data.User.Email, admin.Email)
	}

	t.Run("update report normally", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bambank",
			PhoneNumber: "+6213246467891",
			Email:       "ppouertmlser@gmail.com",
			Password:    "12345678",
		}
		userDTO, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		createReport := map[string]string{
			"lat":     "-7.666369905243495",
			"lng":     "110.66331442645793",
			"note":    "",
			"address": "mataram",
		}

		for k, v := range createReport {
			writer.WriteField(k, v)
		}

		fileName := "jalan.jpg"
		file, err := os.Open(filepath.Join(imagePath, fileName))
		if err != nil {
			t.Error(err)
		}
		part, err := writer.CreateFormFile("image", fileName)
		if err != nil {
			t.Error(err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Error(err)
		}

		writer.Close()
		req := httptest.NewRequest(http.MethodPost, "/api/reports", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+userDTO.Token)
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusCreated, res.Code)

		resBody, _ := ioutil.ReadAll(res.Body)
		createReportResponse := struct {
			Data *entity.Report `json:"data"`
		}{}
		json.Unmarshal(resBody, &createReportResponse)

		updateReportDTO := &model.UpdateReportDTO{
			Status: "Under Repair",
		}
		b, _ := json.Marshal(updateReportDTO)

		req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/reports/%d", createReportResponse.Data.ID), bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+adminDTO.Data.Token)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusOK, res.Code)
	})

	t.Run("update without admin role", func(t *testing.T) {
		createUserDTO := &model.CreateUserDTO{
			Name:        "bambank",
			PhoneNumber: "+62132465331891",
			Email:       "iytyetnnher@gmail.com",
			Password:    "12345678",
		}
		userDTO, res := register(createUserDTO)

		assertResponseCode(t, http.StatusCreated, res.Code)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/reports/%d", 1), bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+userDTO.Token)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("update with invalid payload", func(t *testing.T) {
		updateReportDTO := &model.UpdateReportDTO{
			Status: "anjay mabar",
		}
		b, _ := json.Marshal(updateReportDTO)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/reports/%d", 1), bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+adminDTO.Data.Token)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})

	t.Run("update nonexistent report", func(t *testing.T) {
		updateReportDTO := &model.UpdateReportDTO{
			Status: "Completed",
		}
		b, _ := json.Marshal(updateReportDTO)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/reports/%d", 168), bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+adminDTO.Data.Token)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusNotFound, res.Code)
	})

	t.Run("update with invalid param", func(t *testing.T) {
		updateReportDTO := &model.UpdateReportDTO{
			Status: "Completed",
		}
		b, _ := json.Marshal(updateReportDTO)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/reports/%s", "yoloo"), bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+adminDTO.Data.Token)
		req.Header.Set("Content-Type", "application/json")
		res = httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assertResponseCode(t, http.StatusBadRequest, res.Code)
	})
}
