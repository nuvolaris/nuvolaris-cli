/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
function main() {
    return {
        body: `<html><head>
<link href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.0/css/bootstrap.min.css" 
      rel="stylesheet" id="bootstrap-css">
</head><body>
 <div id="container">
  <div class="row">
   <div class="col-md-8 col-md-offset-2">
    <h4><strong>Pay</strong></h4>
    <form method="POST" action="submit-sendmail">
     <div class="form-group">
       <input type="text" name="credit"
        class="form-control" placeholder="Credit Card Number">
     </div>
     <div class="form-group">
        <input type="text" name="name"
         class="form-control" placeholder="Name">
     </div>
     <div class="form-group">
        <input type="text" name="secret" 
         class="form-control" placeholder="CVC">
     </div>
      <button class="btn btn-default" type="submit" name="button">
        Pay Now
      </button>
    </form>
   </div>
  </div>
 </div>
</body></html>`
    }
}
