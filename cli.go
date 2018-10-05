/*
 * Copyright (C) 2018  CERN for the benefit of the LHCb collaboration
 * Author: Paul Seyfert <pseyfert@cern.ch>
 *
 * This software is distributed under the terms of the GNU General Public
 * Licence version 3 (GPL Version 3), copied verbatim in the file "LICENSE".
 *
 * In applying this licence, CERN does not waive the privileges and immunities
 * granted to it by virtue of its status as an Intergovernmental Organization
 * or submit itself to any jurisdiction.
 */

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "no url for resolution provided\n")
		flag.Usage()
		os.Exit(1)
	}

	myurl := flag.Args()[0]

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := http.Client{
		Transport: tr,
		Timeout:   10 * time.Second}

	req, err := http.NewRequest("GET", myurl, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not make http request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("User-Agent", "redirect-resolver")
	r, reqerr := c.Do(req)
	if reqerr != nil {
		// Get https://github.com/cornellius-gp/gpytorch: remote error: protocol version not supported
		match, err := regexp.MatchString("Get .*: remote error: protocol version not supported", reqerr.Error())
		if err != nil {
			fmt.Fprintf(os.Stderr, "regexp compilation error: %s\n", err)
			os.Exit(2)
		}
		if match {
			re1 := regexp.MustCompile("^Get ")
			re2 := regexp.MustCompile(": remote error: protocol version not supported$")
			intermediate := re1.ReplaceAllString(reqerr.Error(), "")
			bestguess := re2.ReplaceAllString(intermediate, "")
			fmt.Println(bestguess)
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "could not handle request: %s\n", reqerr)
		os.Exit(1)
	}
	// l, err := r.Location()
	// if err != nil {
	//  fmt.Printf("no location: %s", err)
	// } else {
	//  fmt.Printf("returned: %s", l)
	// }
	lastUrl := r.Request.URL.String()
	fmt.Println(lastUrl)
}
