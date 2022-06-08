// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

function main(args) {
    let url = "https://welcome.nuvolaris.io/nuv/"
    if("__OW_API_HOST" in process.env)
      url += process.env["__OW_API_HOST"].replace(/:/g, "-")
    return new Promise(function(resolve) {
        let req = require("https").get(url, 
          (res) => {
           res.on('data', (data) => resolve( {"body":data.toString()}))
          })
        req.end()
        req.on("error", (err) => resolve({"body":"Welcome to Nuvolaris!!"}))
    })
}
