// import axios from 'axios';
// import {
//     AUTH_LOGIN,
//     AUTH_LOGOUT,
//     AUTH_ERROR,
//     AUTH_CHECK,
//     AUTH_GET_PERMISSIONS
// } from 'react-admin';
// import { removeLocalUser } from 'Helpers/localUser';
// import { getUserProperty } from 'Helpers/localUser';

// export default async (type, params) => {
//     if (type === AUTH_LOGIN) {
//         const { isAuth } = params;

//         return isAuth ? Promise.resolve() : Promise.reject();
//     }
//     if (type === AUTH_LOGOUT) {
//         removeLocalUser();
//         return Promise.resolve();
//     }
//     if (type === AUTH_ERROR) {
//         // make this reject in production (probably) and uncomment removeLocalUser
//         // removeLocalUser();
//         return Promise.resolve();
//     }
//     if (type === AUTH_CHECK) {
//         const token = getUserProperty('token');
//         const localUserId = getUserProperty('id');
//         const localUserType = getUserProperty('userType');

//         var config = {
//             headers: { Authorization: token }
//         };
//         const {
//             data: { id, userType }
//         } = await axios.get(`/user/verify`, config);

//         return localUserId === id && localUserType === userType
//             ? Promise.resolve()
//             : Promise.reject();
//     }
//     if (type === AUTH_GET_PERMISSIONS) {
//         const permissions = getUserProperty('permissions');
//         return Promise.resolve(permissions);
//     }
//     return Promise.reject('Unknown method');
// };

const tokenKey = 'token';

const getTokenClaims = token => {
  try {
    return JSON.parse(atob(token.split('.')[1]));
  } catch {
    throw new Error('Invalid key format');
  }
};

export const authProvider = {
  login: async ({ username, password }) => {
    const request = new Request('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
      headers: new Headers({ 'Content-Type': 'application/json' }),
    });
    return fetch(request)
      .then(response => {
        if (response.status < 200 || response.status >= 300) {
          throw new Error(response.statusText);
        }
        return response.json();
      })
      .then(({ token }) => {
        localStorage.setItem('token', token);
      });
  },
  logout: () => {
    localStorage.removeItem(tokenKey);
    return Promise.resolve();
  },
  checkAuth: async () => {
    const token = localStorage.getItem(tokenKey);
    if (!token) {
      throw new Error('No key');
    }

    let isValid = false;
    try {
      const payload = getTokenClaims(token);
      const exp = new Date(payload.exp * 1000);
      if (exp > new Date()) {
        isValid = true;
      }
    } catch {
      throw new Error('Invalid key format');
    }

    if (!isValid) {
      throw new Error('Token is expired');
    }

    return;
  },
  checkError: error => {
    const status = error.status;
    if (status === 401) {
      localStorage.removeItem(tokenKey);
      return Promise.reject();
    }
    return Promise.resolve();
  },
  getPermissions: async () => {
    const token = localStorage.getItem(tokenKey);
    if (!token) {
      return 'guest';
    }

    const payload = getTokenClaims(token);

    return payload.admin ? 'admin' : 'user';
  },
};
