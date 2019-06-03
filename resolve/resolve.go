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

package resolve

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

func Resolve(myurl string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := http.Client{
		Transport: tr,
		Timeout:   10 * time.Second}

	req, err := http.NewRequest("GET", myurl, nil)
	if err != nil {
		err = fmt.Errorf("could not make http request: %v\n", err)
		return "", err
	}
	req.Header.Set("User-Agent", "redirect-resolver")
	r, reqerr := c.Do(req)
	if reqerr != nil {
		// Get https://github.com/cornellius-gp/gpytorch: remote error: protocol version not supported
		match, err := regexp.MatchString("Get .*: remote error: protocol version not supported", reqerr.Error())
		if err != nil {
			err = fmt.Errorf("regexp compilation error: %v\n", err)
			return "", err
		}
		if match {
			re1 := regexp.MustCompile("^Get ")
			re2 := regexp.MustCompile(": remote error: protocol version not supported$")
			intermediate := re1.ReplaceAllString(reqerr.Error(), "")
			bestguess := re2.ReplaceAllString(intermediate, "")
			return bestguess, nil
		}
		err = fmt.Errorf("could not handle request: %v\n", reqerr)
		return "", err
	}
	return r.Request.URL.String(), nil
}
