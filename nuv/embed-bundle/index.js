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
 */

const fs = require('fs');

ctypes = {
    "gif": "image/gif",
    "jpg": "image/jpeg",
    "png": "image/png",
    "ico": "image/vnd.microsoft.icon",
    "ttf": "font/ttf",
    "woff": "font/woff",
    "woff2": "font/woff2",
    "svg": "image/svg"
}

function isBinary(file) {
    return file.endsWith(".gif") ||
        file.endsWith(".jpg") ||
        file.endsWith(".png") ||
        file.endsWith(".ico") ||
        file.endsWith(".ttf") ||
        file.endsWith(".woff") ||
        file.endsWith(".woff2") ||
        file.endsWith(".svg")        
}

// replace base in html and css
function replaceBase(path, body) {
    // filter .html
    if(!path.endsWith(".html"))
        return body
    // calculate toReplace    
    let a = process.env['__OW_ACTION_NAME'].split("/")
    if(a.length == 3) a.splice(-1, 0, "default")
    let toReplace = "/api/v1/web"+ a.join("/")+"/";
    // replace all
    const toFind = /(?<=(src|href)=['"])\//g
    return body.replace(toFind, toReplace)
}

function body(path) {
    let file = `${__dirname}${path}`
    if (!fs.existsSync(file)) {
        file = ${__dirname}/index.html`
    }
    let data = fs.readFileSync(file)
    if(isBinary(path)) 
        return {
            body: data.toString("base64"),
            statusCode: 200,
            headers: {
                "Content-Type": ctypes[path.split(".").pop()]
            }
        }
    return {
        body: replaceBase(path, data.toString("utf-8")),
        statusCode: 200,
    }
}

function check(args) {
    let res = ""
    if (!("__ow_path" in args)) {
        res = "<h1>Error: not deployed as <tt>--web=true</tt></h1>"
    }
    return res
}

function main(args) {
    // check parameters
    res = check(args)
    if (res!="") {
        return { "body": res }
    }
    let path = args['__ow_path'];
    // send body
    if (path != "") {
        return body(path) 
    }
    // return redirect if no path
    return { "body": `<script>location.href += "/"</script>` }
}

module.exports.main = main
