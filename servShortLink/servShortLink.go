package main

import
(
    "crypto/sha256"
    "fmt"
    "net/http"
    "encoding/base64"
	"net/url"
)

const hashLength = 4
const baseUrl = "http://localhost:8080?u="

type fullUrl struct {
    Cnt int
    Url string
}

type getHashFromRequestError struct {
	s string
}

func (e *getHashFromRequestError) Error() string {
	return e.s
}

// таблица ссылок
var tableOfLinks = make( map[string]fullUrl )

func generateShortURL (link string) string {

    sum := sha256.Sum256( []byte( link ) )
    slice := sum[:]
    encoded := base64.StdEncoding.EncodeToString(slice)

    return encoded[:hashLength]
}

func getFullUrl( shortUrl string ) (string, int) {

    val, ok := tableOfLinks[shortUrl]

    if ok == true {
        // увеличивает счетик посещений
        val.Cnt += 1
        tableOfLinks[shortUrl] = val
        return val.Url, 303
    }
    return "", 404
}

func getHashFromRequest( r *http.Request ) (string, error) {

	if (len(r.URL.RawQuery) > hashLength) {
		return r.URL.RawQuery[len(r.URL.RawQuery)- hashLength:], nil

	} else {
		err := getHashFromRequestError{s: "Невозможно получить короткую ссылку из запроса \n"}
		return "", &err
	}
}

func redirect( w http.ResponseWriter, r *http.Request ) {

	link, err := getHashFromRequest( r )
	if err == nil {
		shortUrl, _ := url.QueryUnescape(link)
		url, statusCode := getFullUrl( shortUrl )

		if statusCode == 303 {
			http.Redirect(w, r, url, http.StatusSeeOther)

		} else {
			http.Redirect(w, r, shortUrl, http.StatusNotFound)
		}
	} else {
		fmt.Println(err)
	}
}

func getLinkCnt( shortUrl string ) int {

	val, ok := tableOfLinks[shortUrl]
	if ok == true {
		return val.Cnt
	}
	return -1
}

func getLinkStatistiсs( w http.ResponseWriter, r *http.Request ) {

	hash, err := getHashFromRequest( r )

	if err == nil {
		encHash, _ := url.QueryUnescape( hash )
		cnt := getLinkCnt( encHash )

		if cnt >= 0 {
			str := fmt.Sprintf("Адрес %s посещали %d раз \n", baseUrl + encHash, cnt )
			w.Write([]byte(str))

		} else {
			str := fmt.Sprintf("Адрес %s не существует в системе\n",
				baseUrl + encHash )
			w.Write([]byte(str))
		}
	} else {
		fmt.Println(err)
	}
}

func registerNewLink( w http.ResponseWriter, r *http.Request ) {

	query, _ := url.QueryUnescape(r.URL.RawQuery)

	if len(query) > 8 {
		url := query[4:]
		newStruct := fullUrl{ Url: url }
		hash := generateShortURL(url)
		tableOfLinks[hash] = newStruct

		str := fmt.Sprintf("Короткая ссылка для %s - %s \n", url, baseUrl + hash)
		w.Write([]byte(str))

	} else {
		w.Write([]byte("Введенный url неполный! \n"))
	}
}

func main () {

	s := &http.Server{
		Addr:           ":8080",
	}

	// обработчики запросов:
	// - перейти по короткой ссылке
	http.HandleFunc("/", redirect)
	// - получить новую короткую ссылку
	http.HandleFunc("/reg_new_link/", registerNewLink)
	// - получить статистику переходов по короткой ссылке
	http.HandleFunc("/get_link_statistiсs/", getLinkStatistiсs)

	//запускаем сервер
	s.ListenAndServe()
}
