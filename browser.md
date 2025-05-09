# browser.htm Explanation

The `browser.htm` file is a simple HTML-based tool designed to act as a lightweight API testing interface, similar to Postman. It allows users to send `GET` and `POST` requests to the server and view the responses directly in their browser. Below is a detailed explanation of its components:

---

## Purpose
This file provides a user-friendly interface for testing API endpoints. Users can:
- Enter a request URL and payload.
- Send `GET` or `POST` requests.
- View the server's response, including status codes and JSON data.

---

## Structure

### HTML Elements
1. **Title and Description**:
   ```html
   <h1>Hello Browser user,</h1>
   <p>We at PowerDale have got you covered. Just enter the request you want to send below and click the button to get the response.</p>
   ```
   - Displays a welcome message and instructions for using the tool.

2. **Request URL Input**:
   ```html
   <label for="path">Request URL:</label>
   <div>
       <span id="host">http://www.example.com/</span>
       <input id="path" type="text" placeholder="path/to/request" />
       <input type="button" value="reset" onclick="reset()" />
   </div>
   ```
   - Allows users to enter the API endpoint path.
   - Displays the base URL (`host`) dynamically based on the current server.

3. **Request Payload Input**:
   ```html
   <label for="request">Request Payload:</label>
   <textarea id="request" autofocus="autofocus" placeholder="{&#10;  &#34;message&#34;: &#34;put your request payload here...&#34;&#10;}"></textarea>
   ```
   - Provides a text area for users to input the JSON payload for `POST` requests.

4. **Request Buttons**:
   ```html
   <div>
       <input type="button" value="Send GET Request" onclick="send()" />
       <input type="button" value="Send POST Request" onclick="post()" />
   </div>
   ```
   - Buttons to send `GET` or `POST` requests.

5. **Response Display**:
   ```html
   <p>Server returned with HTTP Status <span id="status"></span></p>
   <label for="response">Response:</label>
   <textarea id="response" placeholder="{&#10;  &#34;message&#34;: &#34;here will be your answer...&#34;&#10;}"></textarea>
   ```
   - Displays the HTTP status code and the server's response in a text area.

---

### CSS Styling
```html
<style>
    html, body {
        height: 100%;
    }
    body {
        background-color: beige;
        display: flex;
        flex-direction: column;
    }
    input, textarea {
        font: inherit;
    }
    textarea {
        width: 100%;
        flex-grow: 1;
        font-family: monospace;
    }
    textarea[id='response'] {
        flex-grow: 2;
    }
    label {
        font-weight: bold;
    }
    .ok {
        color: green;
    }
    .err {
        color: red;
    }
</style>
```
- Provides a clean and simple layout.
- Highlights successful responses in green (`.ok`) and errors in red (`.err`).

---

### JavaScript Logic
1. **Dynamic Server URL**:
   ```javascript
   function server() {
       return window.location.protocol + '//' + window.location.host;
   }
   ```
   - Dynamically determines the server's base URL.

2. **Reset Functionality**:
   ```javascript
   function reset() {
       document.getElementById('host').innerText = server();
       document.getElementById('path').value = window.location.pathname;
   }
   ```
   - Resets the request path to the current URL.

3. **Send POST Request**:
   ```javascript
   function post() {
       send(document.getElementById('request').value);
   }
   ```
   - Sends a `POST` request with the payload from the `request` text area.

4. **Handle Responses**:
   ```javascript
   function respond(target, status) {
       document.getElementById('status').className = (target.status === 200) ? 'ok' : 'err';
       document.getElementById('status').innerText = status || (target.status + ' - ' + target.statusText);
       document.getElementById('response').value = target.response;

       var contentType = target.getResponseHeader('Content-Type');
       if (contentType && contentType.includes('application/json')) {
           var jsonResponse = JSON.parse(target.response);
           var prettyJson = JSON.stringify(jsonResponse, null, 2);
           document.getElementById('response').value = prettyJson;
       }
   }
   ```
   - Updates the status and response fields based on the server's response.
   - Formats JSON responses for better readability.

5. **Send Requests**:
   ```javascript
   function send(data) {
       var XHR = new XMLHttpRequest();
       XHR.addEventListener('load', function(event) {
           respond(event.target);
       });

       XHR.addEventListener('error', function(event) {
           respond(event.target, 'error in request');
       });

       XHR.open(data !== undefined ? 'POST' : 'GET', server() + document.getElementById('path').value);
       XHR.setRequestHeader('Content-Type', 'application/json');
       XHR.setRequestHeader('Accept', 'application/json');
       XHR.send(data);
   }
   ```
   - Sends `GET` or `POST` requests using `XMLHttpRequest`.
   - Sets headers for `Content-Type` and `Accept` as `application/json`.

---

## Key Points
- **Dynamic Server Detection**: Automatically detects the server's base URL and updates the interface.
- **Request Flexibility**: Supports both `GET` and `POST` requests with customizable payloads.
- **Response Handling**: Displays HTTP status codes and formats JSON responses for readability.
- **User-Friendly Interface**: Provides a simple and intuitive interface for testing API endpoints.

---

## Usage
This file is served by the `router/router.go` file when a browser-based request is made to an unknown route. It allows users to manually test API endpoints by entering request URLs and payloads.

---

## Example
1. Open the browser and navigate to the server's base URL.
2. Enter the API endpoint path in the "Request URL" field.
3. Enter a JSON payload in the "Request Payload" field (for `POST` requests).
4. Click "Send GET Request" or "Send POST Request" to send the request.
5. View the response and status code in the "Response" field.

This file acts as a lightweight API testing tool for developers and users.