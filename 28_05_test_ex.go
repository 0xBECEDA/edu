package main

import
(
    "crypto/sha256"
    "fmt"
    "net/http"
    "encoding/base64"
)


type fullUrl struct {
    Cnt int
    Url string
}

var tableOfLinks = make( map[string]fullUrl )

func generateShortURL (link string) string {

    sum := sha256.Sum256( []byte( link ) )
    slice := sum[:]
    encoded := base64.StdEncoding.EncodeToString(slice)

    fmt.Printf("encoded %s \n", encoded[:8])
    return encoded[:8]
}

func getFullUrl( shortUrl string ) (string, int) {

    val, ok := tableOfLinks[shortUrl]

    if ok == true {
        val.Cnt += 1
        tableOfLinks[shortUrl] = val
        return val.Url, 303
    }
    return "", 404
}

func getShortLinkFromRequest( r *http.Request ) string {

    return r.URL.Path[len(r.URL.Path)-8:]
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
        str := fmt.Sprintf("Этот адрес %s посещали %d раз \n", shortUrl, cnt )
        w.Write([]byte(str))

    } else {
        str := fmt.Sprintf("Адрес %s не существует в системе\n", shortUrl )
        w.Write([]byte(str))
    }
}

func main () {

    s := &http.Server{
          Addr:           ":8080",
    }

    link := "https://www.google.com"
    hash := generateShortURL(link)
    fmt.Printf("hash %s \n", hash )

    test := fullUrl{ Url: link }
    tableOfLinks[hash] = test

    http.HandleFunc("/get_full_url/", redirect)
    http.HandleFunc("/get_link_statistiсs/", getLinkStatistiсs)

    //запускаем сервер
    s.ListenAndServe()
}
