package main

import (
	"api-gin/controllers"
	"api-gin/database"
	"api-gin/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupRotasDeTeste() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	rotas := gin.Default()
	return rotas
}

var ID int
var alunoControle models.Aluno

func CriaAlunoMock() {
	aluno := models.Aluno{Nome: "Aluno Teste", CPF: "99999999999", RG: "9999"}
	database.DB.Create(&aluno)
	ID = int(aluno.ID)
	alunoControle = aluno
}

func DeletaAlunoMock() {
	var aluno models.Aluno
	database.DB.Delete(&aluno, ID)
}

func TestVerificaStatusCodeSaudacao(t *testing.T) {
	r := SetupRotasDeTeste()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/gui", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	if resposta.Code != http.StatusOK {
		t.Fatalf("Status error: valor recebido foi %d e o esperado era %d", resposta.Code, http.StatusOK)
	}
	mockDaRasposta := `{"API diz":"Olá gui"}`
	respostaBody, _ := io.ReadAll(resposta.Body)
	assert.Equal(t, mockDaRasposta, string(respostaBody), "unexpected body value")
}

func TestListandoAlunosHandler(t *testing.T) {
	database.ConectaDB()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRotasDeTeste()
	r.GET("/alunos", controllers.RetornaAlunos)
	req, _ := http.NewRequest("GET", "/alunos", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code)
}

func TestBuscaPorCPF(t *testing.T) {
	database.ConectaDB()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRotasDeTeste()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	req, _ := http.NewRequest("GET", "/alunos/cpf/99999999999", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code)
}

func TestBuscaPorId(t *testing.T) {
	database.ConectaDB()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupRotasDeTeste()
	r.GET("/alunos/:id", controllers.BuscaAlunoPorId)
	path := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", path, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	var aluno models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &aluno)
	assert.Equal(t, alunoControle.Nome, aluno.Nome)
	assert.Equal(t, alunoControle.RG, aluno.RG)
	assert.Equal(t, alunoControle.CPF, aluno.CPF)
	assert.Equal(t, http.StatusOK, resposta.Code)
}

func TestDelete(t *testing.T) {
	database.ConectaDB()
	CriaAlunoMock()
	r := SetupRotasDeTeste()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	path := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", path, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code)
	DeletaAlunoMock()
}

func TestPatch(t *testing.T) {
	database.ConectaDB()
	CriaAlunoMock()
	r := SetupRotasDeTeste()
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	aluno := models.Aluno{Nome: "Aluno Teste2", CPF: "99999999990", RG: "99999"}
	alunoJson, _ := json.Marshal(aluno)
	path := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", path, bytes.NewBuffer(alunoJson))
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	var alunoComparacao models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &alunoComparacao)
	fmt.Println(resposta.Body)
	DeletaAlunoMock()
	// assert.Equal(t, aluno.Nome, alunoComparacao.Nome)
	// assert.Equal(t, aluno.CPF, alunoComparacao.CPF)
	// assert.Equal(t, aluno.RG, alunoComparacao.RG)
}
