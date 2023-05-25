package server

// fazer um teste pra quando o id existe e outro pra quando não existe

// func TestUpdateProblemReceivesProblemId(t *testing.T) {
// 	db := &mocks.Database{}
// 	db.On("UpdateProblem", mock.MatchedBy(func(problem models.Problem) bool {
// 		if problem.ID != "test" {
// 			t.Errorf("expected %v, got %v", "test", problem.ID)
// 		}
// 		return true
// 	})).Return(nil)

// 	router := buildAPIRouter(db)

// 	w := httptest.NewRecorder()
// 	body := `{"statement": "Pergunta?", "answer": "Resposta"}`
// 	req, _ := http.NewRequest("POST", "/api/problem/test", strings.NewReader(body))
// 	router.ServeHTTP(w, req)
// }
