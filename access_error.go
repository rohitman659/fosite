/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 *
 */

package fosite

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (f *Fosite) WriteAccessError(rw http.ResponseWriter, req AccessRequester, err error) {
	f.writeJsonError(rw, req, err)
}

func (f *Fosite) writeJsonError(rw http.ResponseWriter, requester AccessRequester, err error) {
	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")

	rfcerr := ErrorToRFC6749Error(err).WithLegacyFormat(f.UseLegacyErrorFormat).WithExposeDebug(f.SendDebugMessagesToClients)

	if requester != nil {
		rfcerr = rfcerr.WithLocalizer(f.MessageCatalog, getLangFromRequester(requester))
	}

	js, err := json.Marshal(rfcerr)
	if err != nil {
		if f.SendDebugMessagesToClients {
			errorMessage := EscapeJSONString(err.Error())
			http.Error(rw, fmt.Sprintf(`{"error":"server_error","error_description":"%s"}`, errorMessage), http.StatusInternalServerError)
		} else {
			http.Error(rw, `{"error":"server_error"}`, http.StatusInternalServerError)
		}
		log.Println("here in the marshal error")
		return
	}

	log.Println("the code is", rfcerr.CodeField)
	rw.WriteHeader(rfcerr.CodeField)
	// ignoring the error because the connection is broken when it happens
	_, _ = rw.Write(js)
}
