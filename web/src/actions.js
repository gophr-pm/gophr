export const GET_PACKAGE = 'GET_PACKAGE';
export const GET_PACKAGE_REQUEST = 'GET_PACKAGE_REQUEST';
export const GET_PACKAGE_SUCCESS = 'GET_PACKAGE_SUCCESS';
export const GET_PACKAGE_FAILURE = 'GET_PACKAGE_FAILURE';

export const GET_PACKAGES = 'GET_PACKAGES';
export const GET_PACKAGES_REQUEST = 'GET_PACKAGES_REQUEST';
export const GET_PACKAGES_SUCCESS = 'GET_PACKAGES_SUCCESS';
export const GET_PACKAGES_FAILURE = 'GET_PACKAGES_FAILURE';

export const REGISTER = 'REGISTER';
export const AUTHENTICATION_REQUEST = 'REGISTER_REQUEST';
export const AUTHENTICATION_SUCCESS = 'REGISTER_SUCCESS';
export const AUTHENTICATION_FAILURE = 'REGISTER_FAILURE';
export const VALIDATE_REGISTER_FORM = 'VALIDATE_REGISTER_FORM';

export const LOGIN = 'LOGIN';
export const LOGOUT = 'LOGOUT';
export const AUTHENTICATION_REQUEST = 'AUTHENTICATION_REQUEST';
export const AUTHENTICATION_SUCCESS = 'AUTHENTICATION_SUCCESS';
export const AUTHENTICATION_FAILURE = 'AUTHENTICATION_FAILURE';
export const VALIDATE_LOGIN_FORM = 'VALIDATE_LOGIN_FORM';

export const SET_PASSWORD = 'SET_PASSWORD';
export const SET_PASSWORD_REQUEST = 'SET_PASSWORD_REQUEST';
export const SET_PASSWORD_SUCCESS = 'SET_PASSWORD_SUCCESS';
export const SET_PASSWORD_FAILURE = 'SET_PASSWORD_FAILURE';
export const VALIDATE_PASSWORD_FORM = 'VALIDATE_PASSWORD_FORM';

export const SET_EMAIL = 'SET_EMAIL';
export const SET_EMAIL_REQUEST = 'SET_EMAIL_REQUEST';
export const SET_EMAIL_SUCCESS = 'SET_EMAIL_SUCCESS';
export const SET_EMAIL_FAILURE = 'SET_EMAIL_FAILURE';
export const VALIDATE_EMAIL_FORM = 'VALIDATE_EMAIL_FORM';

export const UPDATE_PROFILE = 'UPDATE_PROFILE';
export const UPDATE_PROFILE_REQUEST = 'UPDATE_PROFILE_REQUEST';
export const UPDATE_PROFILE_SUCCESS = 'UPDATE_PROFILE_SUCCESS';
export const UPDATE_PROFILE_FAILURE = 'UPDATE_PROFILE_FAILURE';
export const PROFILE_PROFILE_FORM = 'VALIDATE_PROFILE_FORM';

export const GET_PROFILE = 'GET_PROFILE';
export const GET_PROFILE_REQUEST = 'GET_PROFILE_REQUEST';
export const GET_PROFILE_SUCCESS = 'GET_PROFILE_SUCCESS';
export const GET_PROFILE_FAILURE = 'GET_PROFILE_FAILURE';

export const GET_SUBCRIPTIONS = 'GET_SUBCRIPTIONS'
export const GET_SUBCRIPTIONS_REQUEST = 'GET_SUBCRIPTIONS_REQUEST';
export const GET_SUBCRIPTIONS_SUCCESS = 'GET_SUBCRIPTIONS_SUCCESS';
export const GET_SUBCRIPTIONS_FAILURE = 'GET_SUBCRIPTIONS_FAILURE';

export const GET_TOKENS = 'GET_TOKENS';
export const GET_TOKENS_REQUEST = 'GET_TOKENS_REQUEST';
export const GET_TOKENS_SUCCESS = 'GET_TOKENS_SUCCESS';
export const GET_TOKENS_FAILURE = 'GET_TOKENS_FAILURE';

const api = '/api/v0/';

export function getPackage(packageName) {
  return {
    // Types of actions to emit before and after
    types: [GET_PACKAGE_REQUEST, GET_PACKAGE_SUCCESS, GET_PACKAGE_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'get',
    resource: 'package',
    payload: { name: packageName},
  };
}

export function getPackages(query) {
  return {
    types: [GET_PACKAGES_REQUEST, GET_PACKAGES_SUCCESS, GET_PACKAGES_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'get',
    resource: 'packages',
    payload: query,
  };
}

export function register(account)  {
  return {
    types: [REGISTER_REQUEST, REGISTER_SUCCESS, REGISTER_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'post',
    resource: 'register',
    payload: account,
  };

export function validateRegisterForm() {
  return {
    type: VALIDATE_REGISTER_FORM
  };
}

export function login(credentials) {
  return {
    types: [LOGIN_REQUEST, LOGIN_SUCCESS, LOGIN_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'post',
    resource: 'login',
    payload: credentials,
  };
}

export function logout() {
  return {
    type: LOGOUT
  };
}

export function validateLoginForm() {
  return {
    type: VALIDATE_LOGIN_FORM
  };
}

export function setPassword(password) {
  return {
    types: [SET_PASSWORD_REQUEST, SET_PASSWORD_SUCCESS, SET_PASSWORD_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'post',
    resource: 'password',
    payload: password,
  };
}

export function validatePasswordForm() {
  return {
    type: VALIDATE_PASSWORD_FORM
  };
}

export function setEmail(email) {
  return {
    types: [SET_EMAIL_REQUEST, SET_EMAIL_SUCCESS, SET_EMAIL_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'post',
    resource: 'email',
    payload: email,
  };
}

export function validateEmailForm() {
  return {
    type: VALIDATE_EMAIL_FORM
  };
}

export function updateProfile(profile) {
  return {
    types: [UPDATE_PROFILE_REQUEST, UPDATE_PROFILE_SUCCESS, UPDATE_PROFILE_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'post',
    resource: 'profile',
    payload: profile,
  };
}

export function validateProfileForm() {
  return {
    type: VALIDATE_PROFILE_FORM
  };
}

export function getProfile(profileName) {
  return {
    types: [GET_PROFILE_REQUEST, GET_PROFILE_SUCCESS, GET_PROFILE_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'get',
    resource: 'profile',
    payload: profileName,
  };
}

export function getSubscriptions(query) {
  return {
    types: [GET_SUBCRIPTIONS_REQUEST, GET_SUBCRIPTIONS_SUCCESS, GET_SUBCRIPTIONS_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'get',
    resource: 'subscription',
    payload: query,
  };
}

export function getTokens() {
  return {
    types: [GET_TOKENS_REQUEST, GET_TOKENS_SUCCESS, GET_TOKENS_FAILURE],

    // Check the cache (optional):
    // shouldCallAPI: (state) => shouldFetchPost(state),

    method: 'get',
    resource: 'tokens',
  };
}
