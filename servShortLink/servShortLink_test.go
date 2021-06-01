package servShortLink

import
(
	"net/http"
	"net/http/httptest"
	"testing"
    "io"
)

// проверяет перенаправление с короткого урла на длинный google.com
func redirectGoogleCom( t *testing.T ) {
	req, err := http.NewRequest("GET", "/get_full_url/rGu2", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirect)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("Возвращен неверный статус: получил %v ожидал %v",
			status, http.StatusOK)
	}
}

// проверяет получение короткой ссылки для google.com
func TestGoogleCom(t *testing.T) {
	req, err := http.NewRequest("GET", "/reg_new_link/https://www.google.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerNewLink)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Возвращен неверный статус: получил %v ожидал %v",
			status, http.StatusOK)

	} else {
		body, err := io.ReadAll(rr.Body)

		if err != nil {
			t.Fatal(err)
		}

		expected:= "Короткая ссылка для https://www.google.com - rGu2 \n"
		got := string(body)

		if expected != got {
			t.Errorf("ожидал: %v получил: %v",
				expected, got)
		}

		redirectGoogleCom(t)
	}
}

// проверяет возвращаемую ошибку, есть дать пустой урл
func TestEmptyUrl(t *testing.T) {
	req, err := http.NewRequest("GET", "/reg_new_link/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerNewLink)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Возвращен неверный статус: получил %v ожидал %v",
			status, http.StatusOK)

	} else {
		body, err := io.ReadAll(rr.Body)

		if err != nil {
			t.Fatal(err)
		}

		expected:= "Введена пустая ссылка! Попробуйте еще раз!\n"
		got := string(body)

		if expected != got {
			t.Errorf("ожидал: %v получил: %v",
				expected, got)
		}
	}
}
