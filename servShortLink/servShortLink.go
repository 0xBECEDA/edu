package servShortLink

import
(
    "crypto/sha256"
    "fmt"
    "net/http"
    "encoding/base64"
)

const shortLinkLength = 4

type fullUrl struct {
    Cnt int
    Url string
}

// таблица ссылок
var tableOfLinks = make( map[string]fullUrl )

func generateShortURL (link string) string {

    sum := sha256.Sum256( []byte( link ) )
    slice := sum[:]
    encoded := base64.StdEncoding.EncodeToString(slice)

    return encoded[:shortLinkLength]
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

func getShortLinkFromRequest( r *http.Request ) string {

    return r.URL.Path[len(r.URL.Path)-shortLinkLength:]
}

func redirect( w http.ResponseWriter, r *http.Request ) {

    shortUrl := getShortLinkFromRequest( r )
    url, statusCode := getFullUrl( shortUrl )

    if statusCode == 303 {
        http.Redirect(w, r, url, http.StatusSeeOther)

    } else {
        http.Redirect(w, r, shortUrl, http.StatusNotFound)
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

    shortUrl := getShortLinkFromRequest( r )
    cnt := getLinkCnt( shortUrl)

    if cnt >= 0 {
        str := fmt.Sprintf("Адрес %s посещали %d раз \n", shortUrl, cnt )
        w.Write([]byte(str))

    } else {
        str := fmt.Sprintf("Адрес %s не существует в системе\n", shortUrl )
        w.Write([]byte(str))
    }
}

func registerNewLink( w http.ResponseWriter, r *http.Request ) {

    if len(r.URL.Path) > 14 {

        url := r.URL.Path[14:]
        newStruct := fullUrl{ Url: url }
        hash := generateShortURL(url)
        tableOfLinks[hash] = newStruct

        str := fmt.Sprintf("Короткая ссылка для %s - %s \n", url, hash)
        w.Write([]byte(str))

    } else {
    w.Write([]byte("Введена пустая ссылка! Попробуйте еще раз!\n"))

    }
}

func main () {

    s := &http.Server{
          Addr:           ":8080",
    }

    // обработчики запросов:
    // - перейти по короткой ссылке
    http.HandleFunc("/get_full_url/", redirect)
    // - получить новую короткую ссылку
    http.HandleFunc("/reg_new_link/", registerNewLink)
    // - получить статистику переходов по короткой ссылке
    http.HandleFunc("/get_link_statistiсs/", getLinkStatistiсs)

    //запускаем сервер
    s.ListenAndServe()
}
