// internal utils
export function jsonGet(url: string) {
  return _fetch("GET", url);
}

export function jsonPost(url: string, body: object) {
  return _fetch("POST", url, body);
}

export function jsonDelete(url: string) {
  return _fetch("DELETE", url);
}

export function jsonPut(url: string, body: object) {
  return _fetch("PUT", url, body);
}

export function _fetch(method: string, url: string, body?: object) {
  return fetch(`/api${url}`, {
    method: method,
    body: JSON.stringify(body),
    headers: {
      // CSRF prevention
      // https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)_Prevention_Cheat_Sheet#Use_of_Custom_Request_Headers
      "X-Requested-With": "XMLHttpRequest",
    },
  })
    .then((response) => response.json())
    .then((result) => {
      if (result.error) {
        return Promise.reject(result.error);
      }
      return Promise.resolve(result);
    })
    .catch((error) => {
      return Promise.reject(error.toString());
    });
}
