package jans

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/moabu/kubelogin/pkg/kubeconfig"
	v1 "k8s.io/api/authentication/v1"
)

type proxy struct {
	targetUrl string
	verbose   bool
}

func StartHandler(listenAddress, targetURL string, verbose bool) error {

	log.Printf("Starting Jans proxy for URL %s\n", targetURL)
	p := &proxy{
		targetUrl: targetURL,
		verbose:   verbose,
	}

	log.Printf("Listening on %v for incoming requests...\n", listenAddress)
	http.HandleFunc("/", p.handler)
	return http.ListenAndServe(listenAddress, nil)
}

func (p *proxy) handler(w http.ResponseWriter, r *http.Request) {

	if p.verbose {
		res, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Printf("could not dump request: %v\n", err)
		} else {
			log.Printf("Received request: \n%v\n", string(res))
		}
	}

	// Read body of POST request
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Could not read request body: %v\n", err)
		writeError(w, err)
		return
	}

	// Unmarshal JSON from POST request to TokenReview object
	// TokenReview: https://github.com/kubernetes/api/blob/master/authentication/v1/types.go
	var tr v1.TokenReview
	err = json.Unmarshal(b, &tr)
	if err != nil {
		log.Printf("could not unmarshal JWT payload: %v\n", err)
		writeError(w, err)
		return
	}

	token, err := kubeconfig.DecodeToken(tr.Spec.Token)
	if err != nil {
		log.Printf("could not decode TokenReview: %v\n", err)
		writeError(w, err)
		return
	}

	if token.Url != p.targetUrl {
		err = fmt.Errorf("target URL does not match")
		log.Println(err)
		writeError(w, err)
		return
	}

	// don't use URL provided in TokenReview object, as this
	// could be tempered with. Instead, use the URL from config.
	cl, err := NewClient(p.targetUrl, token.ClientID, token.ClientPassword)
	if err != nil {
		log.Printf("could not initiate Jans client: %v\n", err)
		writeError(w, err)
		return
	}

	jansUserInfo, err := cl.GetUserInfo(r.Context(), token.AccessToken)
	if err != nil {
		log.Printf("Could not get user info: %v\n", err)
		writeError(w, err)
		return
	}

	if p.verbose {
		log.Printf("Got user info: %+v\n", jansUserInfo)
	}

	respToken := v1.TokenReview{
		Status: v1.TokenReviewStatus{},
	}

	// Set status of TokenReview object
	if jansUserInfo == nil {
		respToken.Status.Authenticated = false
	} else {
		respToken.Status.Authenticated = true
		respToken.Status.User = v1.UserInfo{
			Username: jansUserInfo.UserName,
			UID:      jansUserInfo.UID,
		}
	}

	// Marshal the TokenReview to JSON and send it back
	b, err = json.Marshal(respToken)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Write(b)
	log.Printf("Returning: %s\n", string(b))
}

func writeError(w http.ResponseWriter, err error) {

	respToken := v1.TokenReview{
		Status: v1.TokenReviewStatus{
			Authenticated: false,
			Error:         err.Error(),
		},
	}

	// Marshal the TokenReview to JSON and send it back
	b, err := json.Marshal(respToken)
	if err != nil {
		log.Printf("could not marshal TokenReview: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		fmt.Fprintln(w, err)
		return
	}

	w.Write(b)

}
