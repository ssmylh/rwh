package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/k0kubun/pp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}

func cookieHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Set-Cookie", "VISIT=TRUE")
	if _, ok := r.Header["Cookie"]; ok {
		fmt.Fprintf(w, "<html><body>2 回目以降 </body></html>\n")
	} else {
		fmt.Fprintf(w, "<html><body> 初訪問 </body></html>\n")
	}
}

func handleDigest(w http.ResponseWriter, r *http.Request) {
	// A1 = user:Secret Zone:pass
	// A2 = GET:/digest
	// MD5(A1) = f82a7cabef5eec42b8fa827fd3c86db7
	// MD5(A2) = 72c5182fbc56def0cfe368cd32b37c29
	// nc = 00000001
	// cnonce = NjU2MGZkYTM3MzNkNTRhZmVmNWNkNzk2ZTMwOTg4YmE=
	// A3 = Md5(A1):nonce:nc:cnonce:qop:Md5(A2)
	//    = f82a7cabef5eec42b8fa827fd3c86db7:TgLc25U2BQA=f510a27804 73e18e6587be702c2e67fe2b04afd:00000001:NjU2MGZkYTM3MzNkNTRhZmVmNWNkNzk2ZTMwOTg4YmE=:auth:72c5182fbc56def0cfe368cd32b37c29
	// response = MD5(A3)
	pp.Printf("URL: %s\n", r.URL.String())
	pp.Printf("Query: %v\n", r.URL.Query())
	pp.Printf("Proto: %s\n", r.Proto)
	pp.Printf("Method: %s\n", r.Method)
	pp.Printf("Header: %v\n", r.Header)
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("--body--\n%s\n", string(body))
	if _, ok := r.Header["Authorization"]; !ok {
		w.Header().Add("WWW-Authenticate", `Digest realm="Secret Zone", nonce="TgLc25U2BQA=f510a27804 73e18e6587be702c2e67fe2b04afd", algorithm=MD5, qop="auth"`)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		fmt.Fprintf(w, "<html><body>secret page</body></html>\n")
	}
}

func main() {
	var httpServerr http.Server
	http.HandleFunc("/", handler)
	http.HandleFunc("/cookie", cookieHandler)
	http.HandleFunc("/digest", handleDigest)

	log.Println("start http listening: 18888")
	httpServerr.Addr = ":18888"
	log.Println(httpServerr.ListenAndServe())
}
