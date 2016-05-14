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

export function getPackage() {
  return {
    type: GET_PACKAGE,
  };
}

export function getPackageRequest() {
  return {
    type: GET_PACKAGE_REQUEST
  };
}

export function getPackageSuccess() {
  return {
    type: GET_PACKAGE_SUCCESS
  };
}

export function getPackageFailure() {
  return {
    type: GET_PACKAGE_FAILURE
  };
}

export function getPackages() {
  return {
    type: GET_PACKAGES,
  };
}

export function getPackagesRequest() {
  return {
    type: GET_PACKAGES_REQUEST
  };
}

export function getPackagesSuccess() {
  return {
    type: GET_PACKAGES_SUCCESS
  };
}

export function getPackagesFailure() {
  return {
    type: GET_PACKAGES_FAILURE
  };
}

export function register() {
  return {
    type: REGISTER
  };
}

export function registerRequest() {
  return {
    type: REGISTER_REQUEST
  };
}

export function registerSuccess() {
  return {
    type: REGISTER_SUCCESS
  };
}

export function registerFailure() {
  return {
    type: REGISTER_FAILURE
  };
}

export function validateRegisterForm() {
  return {
    type: VALIDATE_REGISTER_FORM
  };
}

export function login() {
  return {
    type: LOGIN
  };
}

export function logout() {
  return {
    type: LOGOUT
  };
}

export function authenticationRequest() {
  return {
    type: AUTHENTICATION_REQUEST
  };
}

export function authenticationSuccess() {
  return {
    type: AUTHENTICATION_SUCCESS
  };
}

export function authenticationFailure() {
  return {
    type: AUTHENTICATED_FAILURE
  };
}

export function validateLoginForm() {
  return {
    type: VALIDATE_LOGIN_FORM
  };
}

export function setPassword() {
  return {
    type: SET_PASSWORD
  };
}

export function setPasswordRequest() {
  return {
    type: SET_PASSWORD_REQUEST
  };
}

export function setPasswordSuccess() {
  return {
    type: SET_PASSWORD_SUCCESS
  };
}

export function setPasswordFailure() {
  return {
    type: SET_PASSWORD_FAILURE
  };
}

export function validatePasswordForm() {
  return {
    type: VALIDATE_PASSWORD_FORM
  };
}

export function setEmail() {
  return {
    type: SET_EMAIL
  };
}

export function setEmailRequest() {
  return {
    type: SET_EMAIL_REQUEST
  };
}

export function setEmailSuccess() {
  return {
    type: SET_EMAIL_SUCCESS
  };
}

export function setEmailFailure() {
  return {
    type: SET_EMAIL_FAILURE
  };
}

export function validateEmailForm() {
  return {
    type: VALIDATE_EMAIL_FORM
  };
}

export function updateProfile() {
  return {
    type: UPDATE_PROFILE
  };
}

export function updateProfileRequest() {
  return {
    type: UPDATE_PROFILE_REQUEST
  };
}

export function updateProfileSuccess() {
  return {
    type: UPDATE_PROFILE_SUCCESS
  };
}

export function updateProfileFailure() {
  return {
    type: UPDATE_PROFILE_FAILURE
  };
}

export function validateProfileForm() {
  return {
    type: VALIDATE_PROFILE_FORM
  };
}

export function getProfile() {
  return {
    type: GET_PROFILE
  };
}

export function getProfileRequest() {
  return {
    type: GET_PROFILE_REQUEST
  };
}

export function getProfileSuccess() {
  return {
    type: GET_PROFILE_SUCCESS
  };
}

export function getProfileFailure() {
  return {
    type: GET_PROFILE_FAILURE
  };
}

export function getSubscriptions() {
  return {
    type: GET_SUBCRIPTIONS
  };
}

export function getSubscriptionsRequest() {
  return {
    type: GET_SUBCRIPTIONS_REQUEST
  };
}

export function getSubscriptionsSuccess() {
  return {
    type: GET_SUBCRIPTIONS_SUCCESS
  };
}

export function getSubscriptionsFailure() {
  return {
    type: GET_SUBCRIPTIONS_FAILURE
  };
}

export function getTokens() {
  return {
    type: GET_TOKENS
  };
}

export function getTokensRequest() {
  return {
    type: GET_TOKENS_REQUEST
  };
}

export function getTokensSuccess() {
  return {
    type: GET_TOKENS_SUCCESS
  };
}

export function getTokensFailure() {
  return {
    type: GET_TOKENS_FAILURE
  };
}
