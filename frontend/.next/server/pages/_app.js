/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(() => {
var exports = {};
exports.id = "pages/_app";
exports.ids = ["pages/_app"];
exports.modules = {

/***/ "./src/lib/apollo.ts":
/*!***************************!*\
  !*** ./src/lib/apollo.ts ***!
  \***************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   apolloClient: () => (/* binding */ apolloClient)\n/* harmony export */ });\n/* harmony import */ var _apollo_client__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @apollo/client */ \"@apollo/client\");\n/* harmony import */ var _apollo_client__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_apollo_client__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var _apollo_client_link_context__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @apollo/client/link/context */ \"@apollo/client/link/context\");\n/* harmony import */ var _apollo_client_link_context__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(_apollo_client_link_context__WEBPACK_IMPORTED_MODULE_1__);\n/* harmony import */ var _apollo_client_link_error__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @apollo/client/link/error */ \"@apollo/client/link/error\");\n/* harmony import */ var _apollo_client_link_error__WEBPACK_IMPORTED_MODULE_2___default = /*#__PURE__*/__webpack_require__.n(_apollo_client_link_error__WEBPACK_IMPORTED_MODULE_2__);\n\n\n\nconst httpLink = (0,_apollo_client__WEBPACK_IMPORTED_MODULE_0__.createHttpLink)({\n    uri: \"https://127.0.0.1:8080/query\" || 0,\n    // Add fetch options for development with self-signed certificates\n    fetchOptions: {\n        // This helps with CORS in development\n        mode: \"cors\",\n        credentials: \"include\"\n    }\n});\nconst authLink = (0,_apollo_client_link_context__WEBPACK_IMPORTED_MODULE_1__.setContext)((_, { headers })=>{\n    // Get the authentication token from local storage if it exists\n    const token =  false ? 0 : null;\n    return {\n        headers: {\n            ...headers,\n            authorization: token ? `Bearer ${token}` : \"\",\n            // Add additional headers for better CORS handling\n            \"Content-Type\": \"application/json\"\n        }\n    };\n});\n// Add error handling link\nconst errorLink = (0,_apollo_client_link_error__WEBPACK_IMPORTED_MODULE_2__.onError)(({ graphQLErrors, networkError, operation, forward })=>{\n    if (graphQLErrors) {\n        graphQLErrors.forEach(({ message, locations, path })=>console.error(`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`));\n    }\n    if (networkError) {\n        console.error(`[Network error]: ${networkError}`);\n        // Provide helpful error messages for common HTTPS/CORS issues\n        if (networkError.message.includes(\"CORS\") || networkError.message.includes(\"fetch\")) {\n            console.error(\"CORS/Network Error - Make sure:\\n\" + \"1. Backend server is running on https://127.0.0.1:8080\\n\" + \"2. You have accepted the self-signed certificate by visiting https://127.0.0.1:8080 in your browser\\n\" + \"3. Backend CORS is configured for https://127.0.0.1:3000\");\n        }\n    }\n});\nconst apolloClient = new _apollo_client__WEBPACK_IMPORTED_MODULE_0__.ApolloClient({\n    link: (0,_apollo_client__WEBPACK_IMPORTED_MODULE_0__.from)([\n        errorLink,\n        authLink,\n        httpLink\n    ]),\n    cache: new _apollo_client__WEBPACK_IMPORTED_MODULE_0__.InMemoryCache(),\n    defaultOptions: {\n        watchQuery: {\n            errorPolicy: \"all\"\n        },\n        query: {\n            errorPolicy: \"all\"\n        }\n    }\n});\n//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9zcmMvbGliL2Fwb2xsby50cyIsIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7O0FBS3dCO0FBQ2lDO0FBQ0w7QUFFcEQsTUFBTU0sV0FBV0osOERBQWNBLENBQUM7SUFDOUJLLEtBQUtDLDhCQUFtQyxJQUFJLENBQThCO0lBQzFFLGtFQUFrRTtJQUNsRUcsY0FBYztRQUNaLHNDQUFzQztRQUN0Q0MsTUFBTTtRQUNOQyxhQUFhO0lBQ2Y7QUFDRjtBQUVBLE1BQU1DLFdBQVdWLHVFQUFVQSxDQUFDLENBQUNXLEdBQUcsRUFBRUMsT0FBTyxFQUFFO0lBQ3pDLCtEQUErRDtJQUMvRCxNQUFNQyxRQUNKLE1BQTZCLEdBQUdDLENBQWtDLEdBQUc7SUFFdkUsT0FBTztRQUNMRixTQUFTO1lBQ1AsR0FBR0EsT0FBTztZQUNWSSxlQUFlSCxRQUFRLENBQUMsT0FBTyxFQUFFQSxNQUFNLENBQUMsR0FBRztZQUMzQyxrREFBa0Q7WUFDbEQsZ0JBQWdCO1FBQ2xCO0lBQ0Y7QUFDRjtBQUVBLDBCQUEwQjtBQUMxQixNQUFNSSxZQUFZaEIsa0VBQU9BLENBQ3ZCLENBQUMsRUFBRWlCLGFBQWEsRUFBRUMsWUFBWSxFQUFFQyxTQUFTLEVBQUVDLE9BQU8sRUFBRTtJQUNsRCxJQUFJSCxlQUFlO1FBQ2pCQSxjQUFjSSxPQUFPLENBQUMsQ0FBQyxFQUFFQyxPQUFPLEVBQUVDLFNBQVMsRUFBRUMsSUFBSSxFQUFFLEdBQ2pEQyxRQUFRQyxLQUFLLENBQ1gsQ0FBQywwQkFBMEIsRUFBRUosUUFBUSxZQUFZLEVBQUVDLFVBQVUsUUFBUSxFQUFFQyxLQUFLLENBQUM7SUFHbkY7SUFFQSxJQUFJTixjQUFjO1FBQ2hCTyxRQUFRQyxLQUFLLENBQUMsQ0FBQyxpQkFBaUIsRUFBRVIsYUFBYSxDQUFDO1FBRWhELDhEQUE4RDtRQUM5RCxJQUNFQSxhQUFhSSxPQUFPLENBQUNLLFFBQVEsQ0FBQyxXQUM5QlQsYUFBYUksT0FBTyxDQUFDSyxRQUFRLENBQUMsVUFDOUI7WUFDQUYsUUFBUUMsS0FBSyxDQUNYLHNDQUNFLDZEQUNBLDBHQUNBO1FBRU47SUFDRjtBQUNGO0FBR0ssTUFBTUUsZUFBZSxJQUFJakMsd0RBQVlBLENBQUM7SUFDM0NrQyxNQUFNL0Isb0RBQUlBLENBQUM7UUFBQ2tCO1FBQVdQO1FBQVVSO0tBQVM7SUFDMUM2QixPQUFPLElBQUlsQyx5REFBYUE7SUFDeEJtQyxnQkFBZ0I7UUFDZEMsWUFBWTtZQUNWQyxhQUFhO1FBQ2Y7UUFDQUMsT0FBTztZQUNMRCxhQUFhO1FBQ2Y7SUFDRjtBQUNGLEdBQUciLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly9tdXNlLWZyb250ZW5kLy4vc3JjL2xpYi9hcG9sbG8udHM/OWMyNCJdLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQge1xuICBBcG9sbG9DbGllbnQsXG4gIEluTWVtb3J5Q2FjaGUsXG4gIGNyZWF0ZUh0dHBMaW5rLFxuICBmcm9tLFxufSBmcm9tIFwiQGFwb2xsby9jbGllbnRcIjtcbmltcG9ydCB7IHNldENvbnRleHQgfSBmcm9tIFwiQGFwb2xsby9jbGllbnQvbGluay9jb250ZXh0XCI7XG5pbXBvcnQgeyBvbkVycm9yIH0gZnJvbSBcIkBhcG9sbG8vY2xpZW50L2xpbmsvZXJyb3JcIjtcblxuY29uc3QgaHR0cExpbmsgPSBjcmVhdGVIdHRwTGluayh7XG4gIHVyaTogcHJvY2Vzcy5lbnYuTkVYVF9QVUJMSUNfR1JBUEhRTF9VUkwgfHwgXCJodHRwczovLzEyNy4wLjAuMTo4MDgwL3F1ZXJ5XCIsXG4gIC8vIEFkZCBmZXRjaCBvcHRpb25zIGZvciBkZXZlbG9wbWVudCB3aXRoIHNlbGYtc2lnbmVkIGNlcnRpZmljYXRlc1xuICBmZXRjaE9wdGlvbnM6IHtcbiAgICAvLyBUaGlzIGhlbHBzIHdpdGggQ09SUyBpbiBkZXZlbG9wbWVudFxuICAgIG1vZGU6IFwiY29yc1wiLFxuICAgIGNyZWRlbnRpYWxzOiBcImluY2x1ZGVcIixcbiAgfSxcbn0pO1xuXG5jb25zdCBhdXRoTGluayA9IHNldENvbnRleHQoKF8sIHsgaGVhZGVycyB9KSA9PiB7XG4gIC8vIEdldCB0aGUgYXV0aGVudGljYXRpb24gdG9rZW4gZnJvbSBsb2NhbCBzdG9yYWdlIGlmIGl0IGV4aXN0c1xuICBjb25zdCB0b2tlbiA9XG4gICAgdHlwZW9mIHdpbmRvdyAhPT0gXCJ1bmRlZmluZWRcIiA/IGxvY2FsU3RvcmFnZS5nZXRJdGVtKFwiYXV0aC10b2tlblwiKSA6IG51bGw7XG5cbiAgcmV0dXJuIHtcbiAgICBoZWFkZXJzOiB7XG4gICAgICAuLi5oZWFkZXJzLFxuICAgICAgYXV0aG9yaXphdGlvbjogdG9rZW4gPyBgQmVhcmVyICR7dG9rZW59YCA6IFwiXCIsXG4gICAgICAvLyBBZGQgYWRkaXRpb25hbCBoZWFkZXJzIGZvciBiZXR0ZXIgQ09SUyBoYW5kbGluZ1xuICAgICAgXCJDb250ZW50LVR5cGVcIjogXCJhcHBsaWNhdGlvbi9qc29uXCIsXG4gICAgfSxcbiAgfTtcbn0pO1xuXG4vLyBBZGQgZXJyb3IgaGFuZGxpbmcgbGlua1xuY29uc3QgZXJyb3JMaW5rID0gb25FcnJvcihcbiAgKHsgZ3JhcGhRTEVycm9ycywgbmV0d29ya0Vycm9yLCBvcGVyYXRpb24sIGZvcndhcmQgfSkgPT4ge1xuICAgIGlmIChncmFwaFFMRXJyb3JzKSB7XG4gICAgICBncmFwaFFMRXJyb3JzLmZvckVhY2goKHsgbWVzc2FnZSwgbG9jYXRpb25zLCBwYXRoIH0pID0+XG4gICAgICAgIGNvbnNvbGUuZXJyb3IoXG4gICAgICAgICAgYFtHcmFwaFFMIGVycm9yXTogTWVzc2FnZTogJHttZXNzYWdlfSwgTG9jYXRpb246ICR7bG9jYXRpb25zfSwgUGF0aDogJHtwYXRofWBcbiAgICAgICAgKVxuICAgICAgKTtcbiAgICB9XG5cbiAgICBpZiAobmV0d29ya0Vycm9yKSB7XG4gICAgICBjb25zb2xlLmVycm9yKGBbTmV0d29yayBlcnJvcl06ICR7bmV0d29ya0Vycm9yfWApO1xuXG4gICAgICAvLyBQcm92aWRlIGhlbHBmdWwgZXJyb3IgbWVzc2FnZXMgZm9yIGNvbW1vbiBIVFRQUy9DT1JTIGlzc3Vlc1xuICAgICAgaWYgKFxuICAgICAgICBuZXR3b3JrRXJyb3IubWVzc2FnZS5pbmNsdWRlcyhcIkNPUlNcIikgfHxcbiAgICAgICAgbmV0d29ya0Vycm9yLm1lc3NhZ2UuaW5jbHVkZXMoXCJmZXRjaFwiKVxuICAgICAgKSB7XG4gICAgICAgIGNvbnNvbGUuZXJyb3IoXG4gICAgICAgICAgXCJDT1JTL05ldHdvcmsgRXJyb3IgLSBNYWtlIHN1cmU6XFxuXCIgK1xuICAgICAgICAgICAgXCIxLiBCYWNrZW5kIHNlcnZlciBpcyBydW5uaW5nIG9uIGh0dHBzOi8vMTI3LjAuMC4xOjgwODBcXG5cIiArXG4gICAgICAgICAgICBcIjIuIFlvdSBoYXZlIGFjY2VwdGVkIHRoZSBzZWxmLXNpZ25lZCBjZXJ0aWZpY2F0ZSBieSB2aXNpdGluZyBodHRwczovLzEyNy4wLjAuMTo4MDgwIGluIHlvdXIgYnJvd3NlclxcblwiICtcbiAgICAgICAgICAgIFwiMy4gQmFja2VuZCBDT1JTIGlzIGNvbmZpZ3VyZWQgZm9yIGh0dHBzOi8vMTI3LjAuMC4xOjMwMDBcIlxuICAgICAgICApO1xuICAgICAgfVxuICAgIH1cbiAgfVxuKTtcblxuZXhwb3J0IGNvbnN0IGFwb2xsb0NsaWVudCA9IG5ldyBBcG9sbG9DbGllbnQoe1xuICBsaW5rOiBmcm9tKFtlcnJvckxpbmssIGF1dGhMaW5rLCBodHRwTGlua10pLFxuICBjYWNoZTogbmV3IEluTWVtb3J5Q2FjaGUoKSxcbiAgZGVmYXVsdE9wdGlvbnM6IHtcbiAgICB3YXRjaFF1ZXJ5OiB7XG4gICAgICBlcnJvclBvbGljeTogXCJhbGxcIixcbiAgICB9LFxuICAgIHF1ZXJ5OiB7XG4gICAgICBlcnJvclBvbGljeTogXCJhbGxcIixcbiAgICB9LFxuICB9LFxufSk7XG4iXSwibmFtZXMiOlsiQXBvbGxvQ2xpZW50IiwiSW5NZW1vcnlDYWNoZSIsImNyZWF0ZUh0dHBMaW5rIiwiZnJvbSIsInNldENvbnRleHQiLCJvbkVycm9yIiwiaHR0cExpbmsiLCJ1cmkiLCJwcm9jZXNzIiwiZW52IiwiTkVYVF9QVUJMSUNfR1JBUEhRTF9VUkwiLCJmZXRjaE9wdGlvbnMiLCJtb2RlIiwiY3JlZGVudGlhbHMiLCJhdXRoTGluayIsIl8iLCJoZWFkZXJzIiwidG9rZW4iLCJsb2NhbFN0b3JhZ2UiLCJnZXRJdGVtIiwiYXV0aG9yaXphdGlvbiIsImVycm9yTGluayIsImdyYXBoUUxFcnJvcnMiLCJuZXR3b3JrRXJyb3IiLCJvcGVyYXRpb24iLCJmb3J3YXJkIiwiZm9yRWFjaCIsIm1lc3NhZ2UiLCJsb2NhdGlvbnMiLCJwYXRoIiwiY29uc29sZSIsImVycm9yIiwiaW5jbHVkZXMiLCJhcG9sbG9DbGllbnQiLCJsaW5rIiwiY2FjaGUiLCJkZWZhdWx0T3B0aW9ucyIsIndhdGNoUXVlcnkiLCJlcnJvclBvbGljeSIsInF1ZXJ5Il0sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./src/lib/apollo.ts\n");

/***/ }),

/***/ "./src/pages/_app.tsx":
/*!****************************!*\
  !*** ./src/pages/_app.tsx ***!
  \****************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"default\": () => (/* binding */ App)\n/* harmony export */ });\n/* harmony import */ var react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react/jsx-dev-runtime */ \"react/jsx-dev-runtime\");\n/* harmony import */ var react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var _apollo_client__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @apollo/client */ \"@apollo/client\");\n/* harmony import */ var _apollo_client__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(_apollo_client__WEBPACK_IMPORTED_MODULE_1__);\n/* harmony import */ var _lib_apollo__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../lib/apollo */ \"./src/lib/apollo.ts\");\n/* harmony import */ var _index_css__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../index.css */ \"./src/index.css\");\n/* harmony import */ var _index_css__WEBPACK_IMPORTED_MODULE_3___default = /*#__PURE__*/__webpack_require__.n(_index_css__WEBPACK_IMPORTED_MODULE_3__);\n\n\n\n\nfunction App({ Component, pageProps }) {\n    return /*#__PURE__*/ (0,react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxDEV)(_apollo_client__WEBPACK_IMPORTED_MODULE_1__.ApolloProvider, {\n        client: _lib_apollo__WEBPACK_IMPORTED_MODULE_2__.apolloClient,\n        children: /*#__PURE__*/ (0,react_jsx_dev_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxDEV)(Component, {\n            ...pageProps\n        }, void 0, false, {\n            fileName: \"/home/samuel/code/Muse/frontend/src/pages/_app.tsx\",\n            lineNumber: 9,\n            columnNumber: 7\n        }, this)\n    }, void 0, false, {\n        fileName: \"/home/samuel/code/Muse/frontend/src/pages/_app.tsx\",\n        lineNumber: 8,\n        columnNumber: 5\n    }, this);\n}\n//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9zcmMvcGFnZXMvX2FwcC50c3giLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7O0FBQ2dEO0FBQ0g7QUFDdkI7QUFFUCxTQUFTRSxJQUFJLEVBQUVDLFNBQVMsRUFBRUMsU0FBUyxFQUFZO0lBQzVELHFCQUNFLDhEQUFDSiwwREFBY0E7UUFBQ0ssUUFBUUoscURBQVlBO2tCQUNsQyw0RUFBQ0U7WUFBVyxHQUFHQyxTQUFTOzs7Ozs7Ozs7OztBQUc5QiIsInNvdXJjZXMiOlsid2VicGFjazovL211c2UtZnJvbnRlbmQvLi9zcmMvcGFnZXMvX2FwcC50c3g/ZjlkNiJdLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgdHlwZSB7IEFwcFByb3BzIH0gZnJvbSBcIm5leHQvYXBwXCI7XG5pbXBvcnQgeyBBcG9sbG9Qcm92aWRlciB9IGZyb20gXCJAYXBvbGxvL2NsaWVudFwiO1xuaW1wb3J0IHsgYXBvbGxvQ2xpZW50IH0gZnJvbSBcIi4uL2xpYi9hcG9sbG9cIjtcbmltcG9ydCBcIi4uL2luZGV4LmNzc1wiO1xuXG5leHBvcnQgZGVmYXVsdCBmdW5jdGlvbiBBcHAoeyBDb21wb25lbnQsIHBhZ2VQcm9wcyB9OiBBcHBQcm9wcykge1xuICByZXR1cm4gKFxuICAgIDxBcG9sbG9Qcm92aWRlciBjbGllbnQ9e2Fwb2xsb0NsaWVudH0+XG4gICAgICA8Q29tcG9uZW50IHsuLi5wYWdlUHJvcHN9IC8+XG4gICAgPC9BcG9sbG9Qcm92aWRlcj5cbiAgKTtcbn1cbiJdLCJuYW1lcyI6WyJBcG9sbG9Qcm92aWRlciIsImFwb2xsb0NsaWVudCIsIkFwcCIsIkNvbXBvbmVudCIsInBhZ2VQcm9wcyIsImNsaWVudCJdLCJzb3VyY2VSb290IjoiIn0=\n//# sourceURL=webpack-internal:///./src/pages/_app.tsx\n");

/***/ }),

/***/ "./src/index.css":
/*!***********************!*\
  !*** ./src/index.css ***!
  \***********************/
/***/ (() => {



/***/ }),

/***/ "@apollo/client":
/*!*********************************!*\
  !*** external "@apollo/client" ***!
  \*********************************/
/***/ ((module) => {

"use strict";
module.exports = require("@apollo/client");

/***/ }),

/***/ "@apollo/client/link/context":
/*!**********************************************!*\
  !*** external "@apollo/client/link/context" ***!
  \**********************************************/
/***/ ((module) => {

"use strict";
module.exports = require("@apollo/client/link/context");

/***/ }),

/***/ "@apollo/client/link/error":
/*!********************************************!*\
  !*** external "@apollo/client/link/error" ***!
  \********************************************/
/***/ ((module) => {

"use strict";
module.exports = require("@apollo/client/link/error");

/***/ }),

/***/ "react/jsx-dev-runtime":
/*!****************************************!*\
  !*** external "react/jsx-dev-runtime" ***!
  \****************************************/
/***/ ((module) => {

"use strict";
module.exports = require("react/jsx-dev-runtime");

/***/ })

};
;

// load runtime
var __webpack_require__ = require("../webpack-runtime.js");
__webpack_require__.C(exports);
var __webpack_exec__ = (moduleId) => (__webpack_require__(__webpack_require__.s = moduleId))
var __webpack_exports__ = (__webpack_exec__("./src/pages/_app.tsx"));
module.exports = __webpack_exports__;

})();