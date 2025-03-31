# Setup

Before continuing with the setup process, please make sure to include the One Step GPS API key in the `.env` file of the `assessmentAPI` directory, as well as the Google Maps Javascript API key in the `.env` file of the `assessmentUI` directory.

## Docker
From the root of the project directory, run the command `docker compose up`. This will build a compose stack called `one-step-gps` that contains three containers:

1. `mongodb` => MongoDB database
2. `assessment-api` => Golang server
3. `assessment-ui` => Nginx hosting Vue.js application

## Alternative
`assessmentAPI` directory:
1. In the `.env` file, change the value of the `MONGO_URI` variable to `mongodb://localhost:27017`.
2. In the terminal, run the command `go mod download` to get all required packages.
3. Finally, run the command `make build && make run` to have the server listening on port `8080`.

`assessmentUI` directory:
1. In the terminal, run the command `npm ci`.
2. Then, run the command `npm run build`.
3. Finally, run the command `npm run preview`

Regardless of which method you choose, you can navigate to the user interface in your browser at `http://localhost:4173`.

# Design

## Golang API Server

The API I built is responsible for the essential tasks of storing user preferences, as well as fetching data from the One Step GPS (OSGPS) API and restructuring it into a digestible format for the client application. The user preferences include sorting the devices list based on various conditions, showing/hiding devices in the list, and setting the polling frequency in which the data is sent from the server to the client.

In order to store user preferences, I decided to use anonymous sessions rather than the conventional user sign-in. Essentially, a unique session ID is generated and sent to the client as an `HTTPOnly` cookie. This ID is also used as an identifier within MongoDB to retrieve the user's preference data.

My server uses the url and API key the team provided to communicate with the OSGPS API and retrieve the data for all devices that the given account follows. Simply, my server receives **all** data about every device, parses and rebuilds the data into a simpler structure, and sends the new structure to the client. To establish live updates from the server to the client, I implemented Server-side events (SSE), in which the client uses the `EventSource` API to subscribe to the server for frequent updates. The frequency of the updates can be updated on the user's end.

While implementing SSE, I discovered that the client cannot send cookies along with its request to subscribe to the server. This was especially an issue because it meant that I wouldn't be able to identify the client and its preferences (particularly the polling frequency) and ultimately not provide a personalized experience for the user. I came up with a workaround that involved including a Short-lived token (SLT) in the query parameters of the SSE subscription request. Here is a step-by-step of my workflow to get an authenticated and personalized SSE subscription:

1. The client sends a prefacing `GET` request (before the SSE request) that has the authentication cookie attached. In this case, it is to `/user-preferences`.
2. Upon receiving this request, the server generates a SLT key, which points to an SLT that contains values such as the user's session ID and expiration time. The SLT key is sent back to the client in the response's `Authorization` header.
3. The client takes the SLT key and attaches it as a query parameter to the SSE subscription request.
4. The server receives the SSE request and parses the query parameter for the SLT key. This temporary identifier allows the server to identify the user making the SSE request through the user session ID stored in the token, and apply their user preferences to shape the way the server communicates with the client. In this case, it affects the frequency at which the client receives updates.

## Vue.js Application

The dashboard utilizes the Google Maps JavaScript API to display a map and render markers representing the geolocation of each device. For fetching specific location information, such as addresses, I originally implemented a function in my API that used the Google Geocode API, however, I decided to move away from that implementation as it would have been too expensive. Instead, I used the free, open-source `nominatim` API for reverse geocoding. The implementation details of the API calls and data processing of `nominatim` can be found in the `Geocode.ts` file in the `classes` directory.

For styling, I decided to stick with CSS. While I find component libraries such as ChakraUI to be extremely helpful (especially in writing DRY code), I thought I would demonstrate my skills with the basic tool.

Finally, I have implemented my own compositions to provide context of various components of the application. This was necessary, as the `provide` functionality that Vue offered wasn't sufficient for my needs, being state encapsulation and persistence, and better reusability.

If I had more time to spend on this assignment, I would focus on refactoring the code into even smaller components to make it more manageable as the codebase grew.
